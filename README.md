# github-achievements-tracker

A Go CLI tool to track a user's progress towards GitHub achievements using the GitHub API.

## Common GitHub Achievements

| Badge | Name | Description |
| :---: | :---: | :--- |
| <img src="images/PullShark.png" width="40"/> | **Pull Shark** | Opened pull requests that have been merged. |
| <img src="images/PairExtraordinaire.png" width="40"/> | **Pair Extraordinaire** | Coauthored commits on merged pull requests. |
| <img src="images/YOLO_Badge.png" width="40"/> | **YOLO** | Merged a pull request without code review. |
| <img src="images/QuickDraw_SkinTone1.png" width="40"/> | **Quickdraw** | Closed an issue or pull request within 5 minutes of opening. |
| <img src="images/StarStruck_SkinTone1.png" width="40"/> | **Starstruck** | Created a repository that has received 16+ stars. |
| <img src="images/GalaxyBrain.png" width="40"/> | **Galaxy Brain** | Had answers accepted in GitHub Discussions. |

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
