# Forum

ViableForum is a server-rendered web forum built with Go, SQLite, HTML, CSS, and vanilla JavaScript. Users can register, sign in, publish categorized posts, discuss posts through comments, react to posts and comments, and filter the forum feed.

## Features

- User registration with email, username, and password
- Bcrypt password hashing before storage
- Cookie-based login sessions with expiration
- One active session per user; a new login invalidates previous sessions
- Public post and comment viewing
- Authenticated post creation
- Authenticated comment creation
- Multiple categories per post
- Category/subforum filtering
- Created-post filtering for the logged-in user
- Liked-post filtering for the logged-in user
- Like/dislike reactions for posts and comments
- Publicly visible like and dislike counts
- SQLite persistence with foreign-key enforcement and indexes
- Graceful HTTP server shutdown
- Docker and Docker Compose support
- Unit and integration-style tests across database, repository, service, and HTTP layers

## Technology

- Go 1.24
- SQLite through `github.com/mattn/go-sqlite3`
- Bcrypt through `golang.org/x/crypto/bcrypt`
- UUID generation through `github.com/google/uuid`
- Server-rendered Go `html/template` templates
- Vanilla JavaScript and CSS
- Docker with a multi-stage Alpine build

No frontend framework is required. The application does not use React, Vue, Angular, or another frontend framework.

## Requirements

### Recommended: Docker

- Docker Engine
- Docker Compose v2, available as either `docker compose` or `docker-compose`

### Local Go execution

- Go 1.24 or newer
- A C compiler, because `go-sqlite3` uses CGO
- SQLite development headers/libraries suitable for the host system

On Debian/Ubuntu systems, the local CGO prerequisites can usually be installed with:

```bash
sudo apt-get update
sudo apt-get install -y gcc libc6-dev libsqlite3-dev
```

Docker is recommended because the provided Dockerfile installs the Alpine CGO build dependencies automatically.

## Run With Docker Compose

From the repository root:

```bash
docker-compose up --build
```

Or, with the modern Compose command:

```bash
docker compose up --build
```

The server listens on:

```text
http://localhost:8080
```

To start in the background:

```bash
docker-compose up --build -d
```

To view logs:

```bash
docker-compose logs -f web-forum
```

To stop the application while keeping the database volume:

```bash
docker-compose down
```

To stop the application and delete its persisted SQLite data as well:

```bash
docker-compose down -v
```

> `docker-compose down -v` permanently deletes the named `forum-data` volume and all database records stored in it.

## Run Locally Without Docker

Install dependencies and run the application from the repository root:

```bash
go mod download
CGO_ENABLED=1 go run ./cmd/web
```

The default local paths are:

```text
Database: ./forum.db
Schema:   ./internal/database/schema.sql
```

You can override them with environment variables:

```bash
DB_PATH=./data/forum.db \
SCHEMA_PATH=./internal/database/schema.sql \
CGO_ENABLED=1 \
go run ./cmd/web
```

Create the database directory first when using a nested database path:

```bash
mkdir -p data
go run ./cmd/web
```

## Configuration

The application reads these environment variables:

| Variable | Default | Purpose |
| --- | --- | --- |
| `DB_PATH` | `./forum.db` | SQLite database file location |
| `SCHEMA_PATH` | `./internal/database/schema.sql` | SQL schema file location |

Docker Compose overrides them with:

```text
DB_PATH=/app/data/forum.db
SCHEMA_PATH=/app/internal/database/schema.sql
```

The SQLite database is stored in the named Docker volume `forum-data`, mounted at `/app/data`. Application binaries, templates, and static files remain in the image and are not hidden by the database volume.

## Initial Seed Data

On startup, the application:

1. Creates the database schema if it does not exist.
2. Seeds the default categories.
3. Seeds a demo user if it does not already exist.

The demo account is:

```text
Email:    alpha@zone01.edu
Password: password123
Username: architect_alpha
```

This account exists for local development and testing only. Change or remove seeded credentials before deploying the application outside a development environment.

The default seeded categories are:

- Entertainment
- Sports
- Politics
- Technology
- Others

The homepage currently exposes selected subforums in its sidebar. The database can contain additional categories, and posts can be associated with any category that exists in the database.

## Main Routes

| Method | Route | Access | Description |
| --- | --- | --- | --- |
| `GET` | `/` | Public | Forum homepage and post feed |
| `GET` | `/post/view?id=<id>` | Public | View a post and its comments |
| `GET` | `/login` | Public | Registration and login forms |
| `POST` | `/register` | Public | Create a user account |
| `POST` | `/login` | Public | Authenticate and create a session cookie |
| `GET` | `/logout` | Public/session-aware | Delete the current session and clear the cookie |
| `POST` | `/post/create` | Authenticated | Create a post with one or more categories |
| `POST` | `/comment/create` | Authenticated | Add a comment to a post |
| `POST` | `/interaction` | Authenticated | Like or dislike a post or comment |
| `GET` | `/static/...` | Public | Serve CSS and JavaScript assets |

## Filtering

### Category filter

Category filtering is available to all visitors:

```text
http://localhost:8080/?category=Technology
```

The category name must match an existing category exactly.

### Created posts

Created-post filtering requires an active session:

```text
http://localhost:8080/?filter=created
```

Only posts belonging to the logged-in user are returned.

### Liked posts

Liked-post filtering also requires an active session:

```text
http://localhost:8080/?filter=liked
```

Only posts that the logged-in user has currently liked are returned. Switching a reaction from like to dislike removes the post from this filter.

## Reactions

Reactions are submitted through the protected `POST /interaction` endpoint.

Form fields:

```text
target_id=<post-or-comment-id>
target_type=post|comment
value=1|-1
```

Examples:

```bash
curl -b cookies.txt -c cookies.txt \
  -X POST http://localhost:8080/interaction \
  -d 'target_id=<post-id>' \
  -d 'target_type=post' \
  -d 'value=1'
```

The database enforces one reaction per user, target, and target type. Setting a new reaction replaces the previous reaction for that target. Valid values are:

- `1`: like
- `-1`: dislike

Guests can see reaction totals but cannot create reactions.

## Authentication Behavior

When login succeeds, the server creates a UUID session record in SQLite and sends an `HttpOnly` cookie named `forum_session_token`.

The cookie:

- Is scoped to `/`
- Has a 24-hour expiration
- Uses `SameSite=Lax`
- Cannot be read by client-side JavaScript
- Is validated against the sessions table on protected requests

The application invalidates existing sessions for the user when a new login is created. This keeps one active session per user account.

## Database Structure

The schema is defined in [internal/database/schema.sql](internal/database/schema.sql). The main tables are:

- `users`: registered user accounts
- `sessions`: active login sessions and expiration dates
- `posts`: forum posts
- `categories`: available category names
- `post_categories`: many-to-many post/category relationship
- `comments`: comments associated with posts
- `interactions`: likes and dislikes for posts and comments

Foreign keys are enabled and indexes are created for common lookup paths, including posts by author, comments by post, and interactions by target.

## Project Structure

```text
cmd/web/main.go                    Application entrypoint and dependency wiring
internal/database/                 SQLite initialization, schema, and seed data
internal/delivery/http/             HTTP handlers, routes, and auth middleware
internal/models/                    Domain data structures
internal/repository/                SQLite persistence operations
internal/service/                   Business rules and validation
ui/templates/                      Server-rendered HTML templates
ui/static/css/                     Application styling
ui/static/js/                      Vanilla JavaScript behavior
Dockerfile                          Multi-stage production image
docker-compose.yml                 Local Docker Compose service definition
go.mod                              Go module and dependencies
```

## Testing

Run the complete test suite:

```bash
go test ./...
```

Run tests with race detection where supported by the environment:

```bash
go test -race ./...
```

Run a specific package:

```bash
go test ./internal/service
go test ./internal/delivery/http
```

The test suite covers:

- SQLite initialization
- Repository CRUD behavior
- User and session lifecycle
- Authentication validation
- Post and comment services
- Category and author filters
- Reaction persistence and validation
- HTTP redirects, cookies, status codes, and method guards

## Common Troubleshooting

### SQLite cannot open the database file

Make sure the parent directory exists:

```bash
mkdir -p data
```

With Docker Compose, use the provided volume mount and current Compose configuration. Do not mount the entire `/app` directory over the image, because that hides the compiled binary and application templates. The correct persistent mount is `/app/data`.

### Go version error during Docker build

The module requires Go 1.24. Confirm that the Dockerfile uses:

```dockerfile
FROM golang:1.24-alpine AS builder
```

Rebuild without stale layers if necessary:

```bash
docker-compose build --no-cache
docker-compose up
```

### Login returns to the homepage

A successful login redirects to `/`. The authenticated homepage should show authenticated navigation such as `New post`, `My Blueprints`, `Liked posts`, and `Logout`.

If login returns `401`, verify the exact email and password used during registration. Check the application logs:

```bash
docker-compose logs --tail=100 web-forum
```

### Port 8080 is already in use

Stop the process using the port or change the Compose mapping, for example:

```yaml
ports:
  - "8081:8080"
```

Then open `http://localhost:8081`.

### Template or static files are missing

Run the application from the repository root. The server expects paths such as:

```text
./ui/templates/base.html
./ui/static/
./internal/database/schema.sql
```

## Development Notes

- Keep user-generated content parameterized through SQL queries; do not concatenate user input into SQL statements.
- Keep session tokens in `HttpOnly` cookies and validate them server-side.
- Use the existing service/repository boundaries when adding features.
- Add or update tests with every behavior change.
- Avoid committing local database files, credentials, or cookies.
- The Docker image statically links the CGO-enabled SQLite application and uses Alpine as its runtime image.

## License

No license file is currently included in this repository.
