# Gator - Blog Aggregator CLI

A command-line interface for aggregating and browsing RSS feeds from multiple sources. Gator allows you to register users, follow feeds, and browse posts from your favorite blogs all in one place.

## Features

- **User Management**: Register and login to manage your personal feed collection
- **Feed Management**: Add and browse available RSS feeds
- **Feed Following**: Follow/unfollow feeds to personalize your content stream
- **Post Browsing**: Browse aggregated posts from your followed feeds
- **Automatic Scraping**: Automatically fetch and store new posts from feeds

## Prerequisites

Before running Gator, ensure you have the following installed:

- **Go** (version 1.25.1 or later): [Download and install Go](https://golang.org/dl/)
- **PostgreSQL** (version 12 or later): [Install PostgreSQL](https://www.postgresql.org/download/)

## Installation

### Install via Go Install

To install the `gator` CLI tool globally on your system:

```bash
go install github.com/SyntaxSamurai/Bootdev/BlogAggregator@latest
```

This will compile the project and place the `gator` binary in your `$GOPATH/bin` directory (typically `~/go/bin/` or `/usr/local/go/bin/`).

**Note**: Make sure your `$GOPATH/bin` is in your `$PATH`:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Build from Source

Alternatively, you can build the binary manually:

```bash
git clone https://github.com/SyntaxSamurai/Bootdev/BlogAggregator.git
cd BlogAggregator
go build -o gator
```

## Database Setup

Gator uses PostgreSQL to store users, feeds, and posts. You'll need to set up a database and run migrations.

### 1. Create the Database

```bash
# Connect to PostgreSQL and create the database
createdb gator
```

### 2. Run Database Migrations

Gator uses `goose` for database migrations. Install goose if you haven't already:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Then run the migrations:

```bash
# Set your database connection string
export DB_URL="postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

# Run migrations to create all necessary tables
goose postgres "$DB_URL" up

Example: 
goose postgres "postgres://postgres:postgres@localhost:5433/gator" up

```

The migration files will create the following tables:
- `users` - User accounts
- `feeds` - RSS feed information
- `feed_follows` - User feed subscriptions
- `posts` - Aggregated blog posts

### 3. Verify Database Setup

Connect to your database to verify the tables were created:

```bash
psql "$DB_URL"
```

Then run:
```sql
\dt
```

You should see the tables: `users`, `feeds`, `feed_follows`, `posts`.

## Configuration

Gator uses a configuration file named `.gatorconfig.json` in your home directory to store database connection information and the current user session.

### Create the Configuration File

Create `.gatorconfig.json` in your home directory:

```bash
cat > ~/.gatorconfig.json << EOF
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
EOF
```

Replace the connection string with your actual database credentials if different.

**Important**: The `current_user_name` field will be automatically updated when you login to the application.

## Usage

### User Commands

#### Register a New User
```bash
gator register <username>
```

#### Login
```bash
gator login <username>
```

This sets the current user in the config file and enables authenticated commands.

#### List All Users
```bash
gator users
```

#### Reset Database (Development Only)
```bash
gator reset
```
⚠️ **Warning**: This deletes all data in the database. Use only for development/testing!

### Feed Commands

#### Add a New Feed
```bash
gator addfeed <feed_url>
```

**Example**:
```bash
gator addfeed https://blog.boot.dev/index.xml
```

#### List All Feeds
```bash
gator feeds
```

This shows all feeds in the database with their creation date and fetch count.

#### Follow a Feed
```bash
gator follow <feed_url>
```

Subscribes the current user to the specified feed.

#### Unfollow a Feed
```bash
gator unfollow <feed_url>
```

Unsubscribes the current user from the specified feed.

#### List Followed Feeds
```bash
gator following
```

Shows all feeds that the current user is following.

### Aggregation Commands

#### Aggregate Posts from Feeds
```bash
gator agg [duration]
```

Fetches new posts from all feeds in the database. The optional duration specifies how far back to fetch posts.

**Examples**:
```bash
# Fetch posts from the last 30 minutes
gator agg 30m

# Fetch posts from the last 2 hours
gator agg 2h

# Fetch posts from the last 1 day
gator agg 24h
```

### Browse Commands

#### Browse Posts
```bash
gator browse [limit]
```

Shows posts from feeds that the current user follows. The optional limit parameter specifies the maximum number of posts to display.

**Examples**:
```bash
# Show the default number of posts
gator browse

# Show the last 10 posts
gator browse 10

# Show the last 50 posts
gator browse 50
```

## Development

### Running in Development Mode

For development, you can run the project directly without installing:

```bash
go run . <command> [args]
```

**Examples**:
```bash
go run . register alice
go run . login alice
go run . addfeed https://blog.boot.dev/index.xml
```

### Database Schema

The database schema is managed through SQL migration files in the `sql/schema/` directory:

- `001_users.sql` - Creates users table
- `002_feeds.sql` - Creates feeds table
- `003_feed_follows.sql` - Creates feed_follows table
- `004_add_last_fetched_at_to_feeds.sql` - Adds last_fetched_at column to feeds
- `005_posts.sql` - Creates posts table

### Adding New Queries

To add new database queries:

1. Create a `.sql` file in `sql/queries/`
2. Add your SQL query with the `-- name: QueryName :one/:many/:exec` comment
3. Run `sqlc generate` to generate the Go code
4. Use the generated functions in your handlers

### Project Structure

```
BlogAggregator/
├── internal/
│   ├── config/      # Configuration and command handlers
│   └── database/    # Generated database code
├── sql/
│   ├── schema/      # Database migration files
│   └── queries/     # SQL query files for sqlc
├── .env             # Environment variables (for development)
├── go.mod           # Go module definition
├── go.sum           # Go dependencies
├── main.go          # Application entry point
└── sqlc.yaml        # sqlc configuration
```

## Troubleshooting

### Common Issues

**Error: "connection refused" when running commands**
- Ensure PostgreSQL is running: `pg_ctl status` or `systemctl status postgresql`
- Verify your database URL in `.gatorconfig.json`
- Check that the database exists: `psql -l | grep gator`

**Error: "user already exists"**
- This username is already registered. Choose a different username.

**Error: "no such table"**
- Run the migrations: `goose postgres "$DB_URL" up sql/schema`
- Verify the migrations completed successfully.

**Error: "duplicate key value violates unique constraint"**
- This usually occurs when trying to add a feed URL that already exists
- This is normal behavior for posts with duplicate URLs (they're automatically skipped)

### Development Tips

- Use `go run .` for quick testing during development
- Use `gator` (after `go install`) for production-like testing
- Check logs for detailed error messages
- Use PostgreSQL's `psql` to inspect database state directly

## Building for Production

Gator compiles to a static binary, making it easy to distribute without requiring the Go toolchain on the target machine:

```bash
# Build for the current platform
go build -o gator

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o gator-linux

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o gator-mac

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o gator.exe
```

## Contributing

Feel free to submit issues and enhancement requests!

## License

This project is part of a Boot.dev guided course and is available for educational purposes.

---

**Note**: This CLI tool is designed for educational purposes to demonstrate building a complete Go application with database integration, CLI commands, and RSS feed processing.