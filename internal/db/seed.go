package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/eecopilot/go-course-social/internal/store"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

var usernames = []string{
	"alice", "bob", "charlie", "dave", "eve",
	"frank", "grace", "heidi", "ivan", "judy",
	"karl", "laura", "mallory", "nina", "olivia",
	"peter", "quinn", "rachel", "sam", "tina",
	"ursula", "victor", "wendy", "xander", "yara",
	"zach", "aaron", "bella", "carl", "diana",
	"edward", "fiona", "george", "hannah", "ian",
	"jasmine", "kevin", "lily", "mike", "nora",
	"oscar", "paula", "quincy", "rose", "steve",
	"troy", "uma", "vicky", "will", "xena",
	"yasmine", "zane",
}

var titles = []string{
	"The Art of Minimalism",
	"10 Tips for Healthy Living",
	"Traveling on a Budget",
	"Mastering Time Management",
	"The Benefits of Daily Meditation",
	"Easy Recipes for Busy People",
	"How to Boost Your Creativity",
	"Understanding Mental Health",
	"The Power of Positive Thinking",
	"Exploring Local Hidden Gems",
	"DIY Home Decor Ideas",
	"Getting Started with Gardening",
	"The Future of Technology",
	"Sustainable Living Made Easy",
	"Essential Skills for Remote Work",
	"Fitness Tips for Beginners",
	"The Impact of Climate Change",
	"A Beginner's Guide to Investing",
	"Unplugging from Social Media",
	"The Joy of Reading Books",
}

var contents = []string{
	"Explore the principles of minimalism and how decluttering your space can lead to a more fulfilling life. Discover practical steps to start your journey today.",
	"Learn about the importance of nutrition and exercise. This post will provide you with ten easy-to-implement tips that can make a significant difference in your health.",
	"Traveling doesn't have to break the bank. This guide shares insider tips on how to find affordable accommodations, budget-friendly activities, and money-saving hacks.",
	"Time management is crucial for success. Discover effective techniques like the Pomodoro Technique and time blocking to help you maximize your productivity.",
	"Meditation can transform your mental well-being. Dive into the benefits of daily practice, including stress reduction, improved focus, and emotional health.",
	"Cooking doesn't have to be time-consuming. This post features quick and easy recipes that require minimal ingredients but deliver maximum flavor.",
	"Unlock your creative potential with exercises designed to stimulate your imagination and boost original thinking. Find out how to overcome creative blocks!",
	"Mental health is just as important as physical health. This post addresses common issues and offers strategies for maintaining a healthy mindset.",
	"Positive thinking can change your life. Learn how to cultivate a positive mindset and the benefits it can bring to your personal and professional life.",
	"Discover the hidden gems in your local area that are often overlooked. This post will guide you to unique spots worth exploring, from parks to cafes.",
	"Transform your living space with DIY home decor projects that are easy to execute and budget-friendly. Get inspired to create a space that reflects your style.",
	"Gardening can be therapeutic and rewarding. This guide covers the basics of starting your own garden, including choosing plants and essential care tips.",
	"Stay informed about the latest advancements in technology and how they affect our daily lives. Explore trends like AI, VR, and sustainable tech innovations.",
	"Learn simple steps to incorporate sustainable practices into your daily routine. This post covers eco-friendly habits that can make a difference.",
	"Remote work is here to stay. Discover essential skills and tools that can help you thrive in a virtual work environment and maintain productivity.",
	"Fitness doesn't have to be intimidating. This post offers beginner-friendly workouts and tips to help you get started on your fitness journey.",
	"Climate change is a pressing issue. Understand its impact and learn about actions you can take to contribute to a healthier planet.",
	"Investing can seem daunting. This beginner's guide breaks down the basics of investing, including stocks, bonds, and mutual funds.",
	"Unplugging from social media can be refreshing. Explore the benefits of taking a break and tips for reconnecting with the real world.",
	"Reading can enrich your life. Discover the joys of reading and get recommendations for must-read books across various genres.",
}

var tags = []string{
	"lifestyle",
	"health",
	"travel",
	"food",
	"wellness",
	"technology",
	"fitness",
	"sustainability",
	"personal development",
	"DIY",
	"mental health",
	"productivity",
	"finance",
	"parenting",
	"self-care",
	"books",
	"entrepreneurship",
	"fashion",
	"home improvement",
	"relationships",
}

var comments = []string{
	"Great read! Thanks for sharing!",
	"I love this idea! Can't wait to try it out.",
	"Very insightful. I learned something new today.",
	"This is exactly what I needed. Thank you!",
	"Fantastic tips! I appreciate your advice.",
	"I completely agree with your point!",
	"This post is so relatable. Well done!",
	"Thanks for the inspiration!",
	"I never thought of it that way. Interesting perspective!",
	"Keep up the great work!",
}

func Seed(store store.Storage) error {
	ctx := context.Background()

	// 这里有问题，generate函数都是从0开始。只适应第一次运行，再次运行。user_id为从0开始会有问题，因为上一次已经执行过了。

	users := generateUsers(20)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Printf("failed to create user: %v", err)
		}
	}

	posts := generatePosts(50, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Printf("failed to create post: %v", err)
		}
	}

	comments := generateComments(100, posts, users)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Printf("failed to create comment: %v", err)
		}
	}
	log.Printf("seeding completed")
	return nil
}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)
	for i := 0; i < n; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
		}
	}
	return users
}

func generatePosts(n int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, n)
	for i := 0; i < n; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(n int, posts []*store.Post, users []*store.User) []*store.Comment {
	cms := make([]*store.Comment, n)
	// Create a slice of post IDs
	postIDs := make([]int64, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}
	// Randomly select a post ID from the slice
	// Create a slice of user IDs
	userIDs := make([]int64, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}
	for i := 0; i < n; i++ {
		// Randomly select a user ID from the slice
		postID := postIDs[rand.Intn(len(postIDs))]
		userID := userIDs[rand.Intn(len(userIDs))]
		comment := comments[rand.Intn(len(comments))]
		cms[i] = &store.Comment{
			PostID:  postID,
			UserID:  userID,
			Content: comment,
		}
	}
	return cms
}
