package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jefferykk77/github-achievements-tracker/tracker"
)

const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
)

func getActiveUser(token string) (string, error) {
	cmd := exec.Command("gh", "api", "user", "--jq", ".login")
	if token != "" {
		cmd.Env = append(cmd.Environ(), fmt.Sprintf("GITHUB_TOKEN=%s", token))
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get active user: %s: %w", stderr.String(), err)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func getActiveToken() string {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token
	}

	cmd := exec.Command("gh", "auth", "token")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err == nil {
		if t := strings.TrimSpace(stdout.String()); t != "" {
			return t
		}
	}

	return ""
}

func printProgressBar(percent float64, width int) {
	filled := int(percent / 100.0 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled

	fmt.Print("[")
	fmt.Print(ColorGreen)
	for i := 0; i < filled; i++ {
		fmt.Print("█")
	}
	fmt.Print(ColorGray)
	for i := 0; i < empty; i++ {
		fmt.Print("░")
	}
	fmt.Print(ColorReset)
	fmt.Printf("] %.1f%%\n", percent)
}

func getLevelLabel(level int) string {
	switch level {
	case 0:
		return "None"
	case 1:
		return ColorBlue + "Base / Default" + ColorReset
	case 2:
		return ColorYellow + "Bronze (x2)" + ColorReset
	case 3:
		return ColorCyan + "Silver (x3)" + ColorReset
	case 4:
		return ColorPurple + "Gold (x4)" + ColorReset
	default:
		return "Unknown"
	}
}

func printProgressDashboard(progress tracker.Progress, usePrivate bool) {
	level, current, target, percent := progress.GetLevel(usePrivate)

	fmt.Println()
	fmt.Printf("%s%s%s\n", ColorBold+ColorCyan, progress.Achievement.Name, ColorReset)
	fmt.Printf("%s%s%s\n", ColorGray, progress.Achievement.Description, ColorReset)
	fmt.Println(strings.Repeat("-", 50))
	
	fmt.Printf("Current Level : %s\n", getLevelLabel(level))
	
	if usePrivate {
		fmt.Printf("Progress      : %s%d%s / %d %s (Public: %d, Private: %d)\n", 
			ColorBold, current, ColorReset, target, progress.Achievement.Units, 
			progress.PublicCount, progress.PrivateCount)
	} else {
		fmt.Printf("Progress      : %s%d%s / %d %s (excluding private)\n", 
			ColorBold, current, ColorReset, target, progress.Achievement.Units)
	}
	
	fmt.Print("Bar           : ")
	printProgressBar(percent, 25)

	if level < len(progress.Achievement.Levels) {
		nextThreshold := progress.Achievement.Levels[level]
		needed := nextThreshold - current
		if needed > 0 {
			fmt.Printf("Next Level    : Need %s%d%s more %s to reach %s\n", 
				ColorGreen, needed, ColorReset, progress.Achievement.Units, getLevelLabel(level+1))
		}
	} else {
		fmt.Println("Next Level    : 🎉 Fully maxed out!")
	}
}

func main() {
	userFlag := flag.String("user", "", "GitHub username")
	privateFlag := flag.Bool("private", true, "Include private contributions")
	flag.Parse()

	token := getActiveToken()

	username := *userFlag
	if username == "" {
		activeUser, err := getActiveUser(token)
		if err != nil {
			log.Fatalf("%sError: cannot determine active user, specify with -user%s", ColorRed, ColorReset)
		}
		username = activeUser
	}

	fmt.Printf("%sFetching merged pull requests for %s...\n", ColorGray, username)
	
	prs, err := tracker.FetchPullRequests(username, token)
	if err != nil {
		log.Fatalf("%sError fetching data: %v%s", ColorRed, err, ColorReset)
	}

	pullShark, pairExtra := tracker.CalculateAchievements(prs)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("   %s%sGitHub Achievements Progress Tracker%s\n", ColorBold+ColorYellow, strings.Repeat(" ", 8), ColorReset)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Target User: %s%s%s\n", ColorBold, username, ColorReset)
	fmt.Printf("Total Merged PRs Found: %s%d%s\n", ColorBold, len(prs), ColorReset)

	printProgressDashboard(pullShark, *privateFlag)
	printProgressDashboard(pairExtra, *privateFlag)
	
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("%sNote: Achievements take 24-48 hours to sync on GitHub.%s\n", ColorGray, ColorReset)
}
