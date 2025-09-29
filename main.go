package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sidis405/gator/internal/database"
	"github.com/sidis405/gator/internal/rss"
)
import (
	"errors"
	"fmt"
	"os"

	"github.com/sidis405/gator/internal/config"
)

type state struct {
	db *database.Queries
	c  *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	list map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	command, ok := c.list[cmd.name]
	if !ok {
		return errors.New("command does not exist")
	}
	return command(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.list[name] = f
}

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Printf("cannot read config file: %q", err)
		return
	}

	db, err := sql.Open("postgres", c.DbUrl)
	if err != nil {
		fmt.Printf("error connecting to db: %q", err)
		return
	}
	s := state{
		c:  &c,
		db: database.New(db),
	}
	cmds := commands{list: map[string]func(*state, command) error{
		"login":     handlerLogin,
		"register":  handlerRegister,
		"reset":     handlerReset,
		"users":     handlerUsers,
		"agg":       handlerAgg,
		"addfeed":   middlewareLoggedIn(handlerAddFeed),
		"feeds":     handlerFeeds,
		"follow":    middlewareLoggedIn(handlerFollow),
		"following": middlewareLoggedIn(handlerFollowing),
		"unfollow":  middlewareLoggedIn(handlerUnfollow),
	}}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}
	cmdName := args[1]
	otherParams := args[2:]
	err = cmds.run(&s, command{
		name:      cmdName,
		arguments: otherParams,
	})

	if err != nil {
		fmt.Printf("%q", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func handlerLogin(s *state, cmd command) error {
	ctx := context.Background()
	if len(cmd.arguments) == 0 {
		return errors.New("username is required")
	}
	userName := cmd.arguments[0]

	_, err := s.db.GetUser(ctx, userName)
	if err != nil {
		return errors.New("user does not exist")
	}

	err = s.c.SetUser(userName)
	if err != nil {
		return err
	}
	fmt.Println("the user has been set to", userName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("username is required")
	}
	userName := cmd.arguments[0]

	ctx := context.Background()
	existingUser, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		Name:      userName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	s.c.SetUser(existingUser.Name)
	fmt.Printf("the user %s was registered\n", existingUser.Name)
	fmt.Println(existingUser)
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.DeleteUsers(ctx)
	if err != nil {
		return err
	}

	fmt.Println("all users were deleted")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.c.CurrentUserName {
			fmt.Printf(" * %s (current)\n", user.Name)
		} else {
			fmt.Printf(" * %s\n", user.Name)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("please specify interval: 1s, 1m, 1h")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, currentUser database.User) error {
	if len(cmd.arguments) < 2 {
		return errors.New("usage is: addfeed name url")
	}

	ctx := context.Background()

	feedName := cmd.arguments[0]
	feedUrl := cmd.arguments[1]

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		Name:      feedName,
		Url:       feedUrl,
		UserID:    currentUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, _ command) error {
	ctx := context.Background()
	feeds, err := s.db.FeedsWithUsers(ctx)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf(" - %s (%v) - [%s]\n", feed.Name, feed.Username.String, feed.Url)
	}
	return nil
}

func handlerFollow(s *state, cmd command, currentUser database.User) error {
	if len(cmd.arguments) == 0 {
		return errors.New("url is required")
	}
	ctx := context.Background()
	feed, err := s.db.GetFeed(ctx, cmd.arguments[0])
	if err != nil {
		return err
	}

	newFeedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	fmt.Printf("now you (%s) are following %s\n", newFeedFollow.UserName, newFeedFollow.FeedName)

	return nil
}

func handlerUnfollow(s *state, cmd command, currentUser database.User) error {
	if len(cmd.arguments) == 0 {
		return errors.New("url is required")
	}
	ctx := context.Background()
	feed, err := s.db.GetFeed(ctx, cmd.arguments[0])
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: currentUser.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}

	return nil
}

func handlerFollowing(s *state, _ command, currentUser database.User) error {
	ctx := context.Background()
	following, err := s.db.GetFeedFollowsForUser(ctx, currentUser.ID)
	if err != nil {
		return err
	}
	for _, feed := range following {
		fmt.Printf(" - %s - by %s\n", feed.FeedName, feed.UserName)
	}

	return nil
}

func getCurrentUser(s *state, ctx context.Context) (database.User, error) {
	currentUserName := s.c.CurrentUserName
	if currentUserName == "" {
		return database.User{}, errors.New("no user logged in; run login <username>")
	}
	currentUser, err := s.db.GetUser(ctx, currentUserName)
	if err != nil {
		return database.User{}, err
	}
	return currentUser, nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		currentUser, err := getCurrentUser(s, ctx)
		if err != nil {
			return err
		}
		return handler(s, cmd, currentUser)
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	err = s.db.MarkFeedAsFetched(ctx, database.MarkFeedAsFetchedParams{
		ID:        feed.ID,
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	feedContent, err := rss.FetchFeed(ctx, feed.Url)

	if err != nil {
		return err
	}

	fmt.Sprintf("-- %v --\n", feedContent.Channel.Title)
	for _, item := range feedContent.Channel.Item {
		fmt.Printf(" - %v\n", item.Title)
	}
	return nil

	return nil
}
