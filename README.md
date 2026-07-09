# github-achievements-tracker

A Go CLI tool to track a user's progress towards GitHub achievements using the GitHub API.

## Common GitHub Achievements

| Badge | Name | Description |
| :---: | :---: | :--- |
| <img src="https://raw.githubusercontent.com/kavicastelo/github-achievements-guide/main/Media/Badges/Pull-Shark/PNG/PullShark.png" width="40"/> | **Pull Shark** | Opened pull requests that have been merged. |
| <img src="https://raw.githubusercontent.com/kavicastelo/github-achievements-guide/main/Media/Badges/Pair-Extraordinaire/PNG/PairExtraordinaire.png" width="40"/> | **Pair Extraordinaire** | Coauthored commits on merged pull requests. |
| <img src="https://raw.githubusercontent.com/kavicastelo/github-achievements-guide/main/Media/Badges/YOLO/PNG/YOLO_Badge.png" width="40"/> | **YOLO** | Merged a pull request without code review. |
| <img src="https://raw.githubusercontent.com/kavicastelo/github-achievements-guide/main/Media/Badges/Quick-Draw/PNG/Skin-Tones/QuickDraw_SkinTone1.png" width="40"/> | **Quickdraw** | Closed an issue or pull request within 5 minutes of opening. |
| <img src="https://raw.githubusercontent.com/kavicastelo/github-achievements-guide/main/Media/Badges/Star-Struck/PNG/Skin-Tones/StarStruck_SkinTone1.png" width="40"/> | **Starstruck** | Created a repository that has received 16+ stars. |
| <img src="https://raw.githubusercontent.com/kavicastelo/github-achievements-guide/main/Media/Badges/Galaxy-Brain/PNG/GalaxyBrain.png" width="40"/> | **Galaxy Brain** | Had answers accepted in GitHub Discussions. |

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
