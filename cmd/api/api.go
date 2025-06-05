package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eecopilot/go-course-social/docs" // This is required to generate swagger docs
	"github.com/eecopilot/go-course-social/internal/auth"
	"github.com/eecopilot/go-course-social/internal/mailer"
	"github.com/eecopilot/go-course-social/internal/store"
	"github.com/eecopilot/go-course-social/internal/store/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	cacheStorage  cache.Storage
}

type config struct {
	addr        string
	env         string
	version     string
	db          dbConfig
	apiUrl      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redisCfg    redisConfig
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type authConfig struct {
	basicAuth basicConfig
	token     tokenConfig
}

type tokenConfig struct {
	secretKey string
	aud       string
	iss       string
	exp       time.Duration
}

type basicConfig struct {
	username string
	password string
}

type mailConfig struct {
	exp       time.Duration
	sendGrid  sendGridConfig
	fromEmail string
}

type sendGridConfig struct {
	apiKey string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
	// Group
	r.Route("/v1", func(r chi.Router) {
		// r.With(app.BasicAuthMiddleware()).Get("/health", app.healthCheckHandler)

		r.Get("/health", app.healthCheckHandler)

		// docs 文档
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)
			r.Route("/{postId}", func(r chi.Router) {
				// 添加中间件
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
			})
		})

		// users
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)
			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.getUserHandler)

				// follower
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			// feed
			// Group 可以给这个group设置中间件
			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r
}

func (app *application) run(mux *chi.Mux) error {
	// consts for server settings
	const WTime = 20 * time.Second
	const RTime = 10 * time.Second
	const ITime = time.Minute

	// Docs
	docs.SwaggerInfo.Version = app.config.version
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Host = app.config.apiUrl
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: WTime, // max time to write response to the client
		ReadTimeout:  RTime, // max time to read request from the client
		IdleTimeout:  ITime, // max time for connections using TCP Keep-Alive
	}

	// 创建一个用于接收shutdown错误的channel
	// 这个channel用于在goroutine和主线程之间传递shutdown的结果
	shutdown := make(chan error)

	// 启动一个goroutine来监听系统信号，实现优雅关闭
	go func() {
		// 创建一个接收系统信号的channel，缓冲区大小为1
		// 缓冲区大小为1是为了确保信号不会被阻塞
		quit := make(chan os.Signal, 1)

		// 注册要监听的信号：SIGINT (Ctrl+C) 和 SIGTERM (终止信号)
		// 当收到这些信号时，会发送到quit channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// 阻塞等待接收信号，一旦收到信号就继续执行
		s := <-quit

		// 创建一个带超时的context，给服务器5秒时间来完成正在处理的请求
		// 这是优雅关闭的关键：不会立即强制关闭，而是等待现有连接处理完成
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() // 确保context资源被释放

		// 记录服务器正在关闭的日志
		app.logger.Infow("Server is shutting down", "signal", s.String())

		// 调用服务器的Shutdown方法进行优雅关闭，并将结果发送到shutdown channel
		// Shutdown会停止接受新的连接，并等待现有连接完成处理（在超时时间内）
		shutdown <- srv.Shutdown(ctx)
	}()

	// 记录服务器启动日志
	app.logger.Infow("Server has started", "addr", app.config.addr, "env", app.config.env)

	// 启动HTTP服务器，这个调用会阻塞直到服务器关闭
	err := srv.ListenAndServe()

	// 检查服务器关闭的原因
	// 如果错误不是nil且不是正常的服务器关闭错误，则返回错误
	// http.ErrServerClosed 是调用Shutdown()方法时的正常返回值
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	// 等待从shutdown channel接收优雅关闭的结果
	// 这里会阻塞直到上面的goroutine完成srv.Shutdown()调用
	err = <-shutdown
	if err != nil {
		return err
	}

	// 记录服务器已停止的日志
	app.logger.Infow("Server has stopped", "addr", app.config.addr, "env", app.config.env)
	return nil
}
