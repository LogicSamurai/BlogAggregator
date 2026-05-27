package config

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/SyntaxSamurai/Bootdev/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	Config *Config
	DB     *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (c *Commands) Run(state *State, cmd Command) error {
	handler, exists := c.Handlers[cmd.Name]
	if !exists {
		return errors.New("command doesn't exist")
	}
	return handler(state, cmd)
}

func (c *Commands) Register(name string, handler func(*State, Command) error) {
	c.Handlers[name] = handler
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 || len(cmd.Args) > 1 {
		return errors.New("The login handler expects a single argument, the username")
	}

	// check if the given username exist in DB or not
	_, err := s.DB.GetUser(context.Background(), cmd.Args[0])

	if err != nil {
		return fmt.Errorf("No user found with name : %v\n", cmd.Args[0])
	}

	err = Setuser(*s.Config, cmd.Args[0])
	// cfg := *s.Config
	// cfg.CurrentUserName = cmd.Args[0]

	// err := write(cfg)
	if err != nil {
		return err
	}

	fmt.Println("the user has been set.")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 || len(cmd.Args) > 1 {
		return errors.New("Invalid Arguments")
	}

	// 2. Check if user already exists
	existingUser, err := s.DB.GetUser(context.Background(), cmd.Args[0])

	// 3. Handle the "user exists" case
	if err == nil && existingUser.Name == cmd.Args[0] {
		// User already exists - return error
		return fmt.Errorf("user with name '%s' already exists", cmd.Args[0])
	}

	uid := uuid.New()
	userParams := database.CreateUserParams{
		ID:        uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	user, err := s.DB.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	err = Setuser(*s.Config, cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User was created having Data\n NAME : %v \n Created At : %v\n", user.Name, user.CreatedAt)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.DB.ResetUsers(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())

	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Println("* " + user.Name + " (current)")
		} else {
			fmt.Println("* " + user.Name)
		}
	}

	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return errors.New("agg command takes no arguments")
	}

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil
}

func write(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func Setuser(cfg Config, userName string) error {
	cfg.CurrentUserName = userName

	return write(cfg)
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/" + configFileName, nil
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, nil
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	// fmt.Println("File Content: ", string(fileContent))

	cfg := Config{}
	err = json.Unmarshal(fileContent, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedUrl, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()


	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	
	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return &RSSFeed{}, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}
