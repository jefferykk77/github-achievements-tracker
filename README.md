# github-achievements-tracker

A Go CLI tool to track a user's progress towards GitHub achievements (such as Pull Shark and Pair Extraordinaire) using the GitHub API.

## Requirements

- Go 1.18+
- [GitHub CLI (gh)](https://cli.github.com/) installed and authenticated.

## Installation & Usage

Clone the repository and run:

```bash
go run main.go
```

### Options

- `-user <username>`: Specify a target user to check (defaults to the currently authenticated user).
- `-private=false`: Exclude private repository contributions from the counts (defaults to `true`).
