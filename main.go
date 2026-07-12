package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/LogicSamurai/BlogAggregator/internal/config"
	"github.com/LogicSamurai/BlogAggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Errorf("Error: %v", err)
		os.Exit(1)
	}

	//load the database
	db, err := sql.Open("postgres", cfg.DbURL)
	dbQueries := database.New(db)

	state := config.State{
		Config: &cfg,
		DB:     dbQueries,
	}

	commands := config.Commands{
		Handlers: make(map[string]func(*config.State, config.Command) error),
	}
	commands.Register("login", config.HandlerLogin)
	commands.Register("register", config.HandlerRegister)
	commands.Register("reset", config.HandlerReset)
	commands.Register("users", config.HandlerUsers)
	commands.Register("agg", config.HandlerAgg)
	commands.Register("addfeed", config.MiddlewareLoggedIn(config.HandlerAddFeed))
	commands.Register("feeds", config.HandlerFeeds)
	commands.Register("follow", config.MiddlewareLoggedIn(config.HandlerFollow))
	commands.Register("following", config.MiddlewareLoggedIn(config.HandlerFollowing))
	commands.Register("unfollow", config.MiddlewareLoggedIn(config.HandlerUnFollow))
	commands.Register("browse", config.MiddlewareLoggedIn(config.HandlerBrowse))

	if len(os.Args) < 2 {
		fmt.Println("Usage: myprogram <command> [args]")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := config.Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = commands.Run(&state, cmd)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
