# boot.dev Blog Aggregator project

## USAGE:
create "~/.gatorconfig.json" with contents:
{"db_url":"postgres://<postgres_user>:<postgres_passwd>@localhost:5432/gator?sslmode=disable","current_user_name":"<username>"}

gator login <username>
 - Login as user
gator register <username>
 - Add user to DB
gator reset
 - Wipe users from DB (deletes will cascade wiping all DBs)
gator users
 - Get a list of current registered users
gator agg
 - Gather and print all feeds to screen
gator addfeed <name> <url>
 - Add a feed and register to current logged in user
gator feeds
 - List al feeds
gator follow <url>
 - Follow <url> for current logged in user
gator following
 - Get all followed URLs for the current user
gator unfollow <url>
 - Unfollow <url> for current logged in user


## Some ideas to come back to:
- Add sorting and filtering options to the browse command
- Add pagination to the browse command
- Add concurrency to the agg command so that it can fetch more frequently
- Add a search command that allows for fuzzy searching of posts
- Add bookmarking or liking posts
- Add a TUI that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)
- Add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely
- Write a service manager that keeps the agg command running in the background and restarts it if it crashes

