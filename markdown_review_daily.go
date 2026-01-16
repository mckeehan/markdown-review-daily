package main

import (
	"bufio"
	"flag"
	"fmt"
    "log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// CronSchedule represents a parsed cron expression
type CronSchedule struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
}

// ParseCronSchedule parses a cron expression string
func ParseCronSchedule(cronExpr string) (*CronSchedule, error) {
	fields := strings.Fields(cronExpr)
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(fields))
	}

	return &CronSchedule{
		Minute:     fields[0],
		Hour:       fields[1],
		DayOfMonth: fields[2],
		Month:      fields[3],
		DayOfWeek:  fields[4],
	}, nil
}

// MatchesToday checks if the cron schedule is triggered on the given date
// (ignoring the specific hour/minute - just checking if today is a scheduled day)
func (cs *CronSchedule) MatchesToday(t time.Time) bool {
	// Check day of month
	if !matchesCronField(cs.DayOfMonth, t.Day(), 31) {
		return false
	}

	// Check month
	if !matchesCronField(cs.Month, int(t.Month()), 12) {
		return false
	}

	// Check day of week (0-7, where both 0 and 7 are Sunday)
	dow := int(t.Weekday())
	if !matchesCronField(cs.DayOfWeek, dow, 7) {
		// Special case: Sunday can be represented as both 0 and 7
		if dow == 0 {
			if !matchesCronField(cs.DayOfWeek, 7, 7) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

// MatchesTime checks if the cron schedule matches the given time exactly
// (including hour and minute)
func (cs *CronSchedule) MatchesTime(t time.Time) bool {
	// Check minute
	if !matchesCronField(cs.Minute, t.Minute(), 59) {
		return false
	}

	// Check hour
	if !matchesCronField(cs.Hour, t.Hour(), 23) {
		return false
	}

	// Check if today matches
	return cs.MatchesToday(t)
}

// matchesCronField checks if a cron field matches the current value
func matchesCronField(field string, current, max int) bool {
	// Wildcard matches everything
	if field == "*" {
		return true
	}

	// Handle comma-separated values (e.g., "1,15,30")
	if strings.Contains(field, ",") {
		values := strings.Split(field, ",")
		for _, val := range values {
			v, err := strconv.Atoi(strings.TrimSpace(val))
			if err == nil && v == current {
				return true
			}
		}
		return false
	}

	// Handle step values (e.g., "*/5" or "10-20/2")
	if strings.Contains(field, "/") {
		parts := strings.Split(field, "/")
		if len(parts) != 2 {
			return false
		}

		step, err := strconv.Atoi(parts[1])
		if err != nil {
			return false
		}

		rangeExpr := parts[0]

		// Handle */step
		if rangeExpr == "*" {
			return current%step == 0
		}

		// Handle start-end/step
		if strings.Contains(rangeExpr, "-") {
			rangeParts := strings.Split(rangeExpr, "-")
			if len(rangeParts) != 2 {
				return false
			}
			start, err1 := strconv.Atoi(rangeParts[0])
			end, err2 := strconv.Atoi(rangeParts[1])
			if err1 != nil || err2 != nil {
				return false
			}
			if current >= start && current <= end {
				return (current-start)%step == 0
			}
		}
		return false
	}

	// Handle ranges (e.g., "1-5")
	if strings.Contains(field, "-") {
		parts := strings.Split(field, "-")
		if len(parts) != 2 {
			return false
		}
		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return false
		}
		return current >= start && current <= end
	}

	// Direct match
	val, err := strconv.Atoi(field)
	if err != nil {
		return false
	}
	return val == current
}

// ExtractFrontmatterData extracts the review_schedule and title from a markdown file's YAML frontmatter
func ExtractFrontmatterData(filename string) (schedule string, title string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inFrontmatter := false
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Check for YAML frontmatter delimiter
		if line == "---" {
			if !inFrontmatter && lineNum == 1 {
				inFrontmatter = true
				continue
			} else if inFrontmatter {
				// End of frontmatter
				break
			}
		}

		// Extract fields if we're in frontmatter
		if inFrontmatter {
			// Trim leading whitespace to handle indented YAML
			trimmedLine := strings.TrimSpace(line)
			
			if strings.HasPrefix(trimmedLine, "review_schedule:") {
				parts := strings.SplitN(trimmedLine, ":", 2)
				if len(parts) == 2 {
					schedule = strings.TrimSpace(parts[1])
					// Remove quotes if present
					schedule = strings.Trim(schedule, `"'`)
				}
			}
			
			if strings.HasPrefix(trimmedLine, "title:") {
				parts := strings.SplitN(trimmedLine, ":", 2)
				if len(parts) == 2 {
					title = strings.TrimSpace(parts[1])
					// Remove quotes if present
					title = strings.Trim(title, `"'`)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	return schedule, title, nil
}

// FileInfo holds information about a matching file
type FileInfo struct {
	Path     string
	Schedule string
	Title    string
}

// ScanDirectory scans a directory for markdown files with matching review schedules
func ScanDirectory(dir string, currentTime time.Time, exactTime bool) ([]FileInfo, error) {
	var matches []FileInfo

	// Get absolute path of the base directory for computing relative paths
	absDir, err := filepath.Abs(dir)
	if err != nil {
		absDir = dir
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process markdown files
		if !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		// Extract review schedule and title
		schedule, title, err := ExtractFrontmatterData(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: error reading %s: %v\n", path, err)
			return nil
		}

		// Skip files without a review schedule
		if schedule == "" {
			return nil
		}

		// Parse the cron schedule
		cronSchedule, err := ParseCronSchedule(schedule)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: invalid cron expression in %s: %v\n", path, err)
			return nil
		}

		// Check if the schedule matches based on mode
		var isMatch bool
		if exactTime {
			isMatch = cronSchedule.MatchesTime(currentTime)
		} else {
			isMatch = cronSchedule.MatchesToday(currentTime)
		}

		if isMatch {
			// Compute relative path from base directory
			absPath, _ := filepath.Abs(path)
			relPath, err := filepath.Rel(absDir, absPath)
			if err != nil {
				relPath = path
			}

			// If no title found, use filename without extension
			if title == "" {
				title = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			}

			matches = append(matches, FileInfo{
				Path:     relPath,
				Schedule: schedule,
				Title:    title,
			})
		}

		return nil
	})

	return matches, err
}

func main() {
	// Define command-line flags
	exactTime := flag.Bool("exact", false, "Match exact time (hour and minute) instead of just today's date")
	flag.Parse()

	// Get directory from command line args or use current directory
	dir := "."
	args := flag.Args()
	if len(args) > 0 {
		dir = args[0]
	}

	// Resolve symbolic links
	resolvedDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		// If we can't resolve, try to use the original path
		fmt.Fprintf(os.Stderr, "Warning: couldn't resolve symlink for %s: %v\n", dir, err)
		resolvedDir = dir
	}

	// Get current time
	currentTime := time.Now()

	log.Printf("Scanning directory: %s\n", resolvedDir)
	if resolvedDir != dir {
		log.Printf("(resolved from: %s)\n", dir)
	}
	log.Printf("Current time: %s\n", currentTime.Format("2006-01-02 15:04"))
	if *exactTime {
		log.Println("Mode: Exact time matching (hour and minute must match)")
	} else {
		log.Println("Mode: Daily matching (any file scheduled for today)")
	}
	log.Println("----------------------------------------")

	// Scan directory for matching files
	matches, err := ScanDirectory(resolvedDir, currentTime, *exactTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	// Display results as markdown links
	for _, match := range matches {
		// Remove .md extension from path for wiki-link
		pathWithoutExt := strings.TrimSuffix(match.Path, ".md")
		log.Printf("- [%s]([[%s]])\n", match.Title, pathWithoutExt)
		fmt.Printf("- [%s]([[%s]])\n", match.Title, pathWithoutExt)
	}

	log.Println("----------------------------------------")
	if *exactTime {
		log.Printf("Found %d file(s) scheduled for review right now\n", len(matches))
	} else {
		log.Printf("Found %d file(s) scheduled for review today\n", len(matches))
	}
}
