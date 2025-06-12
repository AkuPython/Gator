package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AkuPython/Gator/internal/database"
	"github.com/google/uuid"
)

type commands struct {
	Command map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.Command[cmd.Name]
	if !ok {
		return fmt.Errorf("Command Does not Exist")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.Command[name] = f
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Could not get ID of user: %v - %v", s.cfg.CurrentUserName, err)
		}
		return handler(s, cmd, user)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Must provide (only) Username")
	}
	username := cmd.Args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		log.Fatalf("Non-Existant Username '%v' - %v", username, err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error setting username '%v' - %v", username, err)
	}
	fmt.Printf("Username set to: %v\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Must provide (only) Username")
	}
	username := cmd.Args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: username,
	})
	if err != nil {
		log.Fatalf("Error creating user with username '%v' - %v", username, err)
	}
	fmt.Printf("Username: '%v' created \n", username)
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error setting username '%v' - %v", username, err)
	}
	fmt.Printf("Username set to: %v\n", username)
	log.Print(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	return s.db.DeleteUsers(context.Background())
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("Could not get users from DB - %v", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Println("*", user.Name, "(current)")
		} else {
			fmt.Println("*", user.Name)
		}
	}
	return nil
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("Could not get next feed to fetch", err)
		return
	}
	fmt.Println("fetching feed:", feed.Name)
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Printf("Couldn't mark feed: %v as fetched: %v", feed.ID, err)
		return
	}
		feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		fmt.Println("-", item.Title)
	}
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Must provide interval (1m, 1h, etc)")
	}
	scrapeInterval, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %v", err)
	}
	if scrapeInterval < time.Duration(30 * time.Second) {
		return fmt.Errorf("Interval too short, must be at least 30s")
	}
	log.Printf("Collecting feeds every %s...", scrapeInterval)

	ticker := time.NewTicker(scrapeInterval)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	// return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Must provide name & url")
	}
	
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: cmd.Args[0],
		Url: cmd.Args[1],
		UserID: user.ID,
	})

	if err != nil {
		return fmt.Errorf("Could not create Feed for User:\n\t%v\n\t%v", user.ID, err)
	}
	// fmt.Println(feed)
	
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("Could not create Feed Follow for User: %v, Feed: %v", user.ID, feed.ID)
	}
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Could not get feeds from DB: %v", err)
	}
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Could not get user ID: %v from DB: %v", feed.UserID, err)
		}
		fmt.Printf("Name: %v - URL: %v - Added by: %v\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func handlerCreateFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Must provide (only) url")
	}
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Could not get feed using URL: %v from DB: %v", cmd.Args[0], err)
	}
	
	feed_create, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Could not create feed follow in DB: %v", err)
	}

	fmt.Printf("Feed Name: %v - Current User: %v", feed_create.FeedName, feed_create.UserName)
	return nil
}

func handlerFeedFollowsForUser(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Only runs on current user")
	}
	
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Could not get feeds for user: %v from DB: %v", user.Name, err)
	}
	for _, follow := range feedFollows {
		user, err := s.db.GetUserByID(context.Background(), follow.UserID)
		if err != nil {
			return fmt.Errorf("Could not get user ID: %v from DB: %v", follow.UserID, err)
		}
		fmt.Printf("Name: %v - User: %v\n", follow.Name, user.Name)
	}
	return nil
}

func handlerUnfollowURL(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Must provide URL")
	}
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Could not get feeds for URL: %v from DB: %v", cmd.Args[0], err)
	}
	
	err = s.db.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
		UserID: user.ID,
		FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("Could not remove follow for Feed: %v - User: %v - URL: %v\n%v", feed.Name, user.Name, cmd.Args[0], err)
	}
	return nil
}
