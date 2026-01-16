# Markdown Review Daily

A Go program that scans markdown files for YAML frontmatter containing `review_schedule` fields in cron format and identifies which files are scheduled for review at the current time.

## How It Works

The program checks if **today** is a day that matches the cron schedule, regardless of the current time. For example:
- A file with schedule `"0 9 * * *"` (daily at 9 AM) will be shown **all day long** on every day
- A file with schedule `"0 10 * * 1"` (Mondays at 10 AM) will be shown **all day on Mondays**

This behavior makes the tool useful as a daily review reminder - you can run it once in the morning to see all files scheduled for review that day.

## Installation

### Prerequisites
- Go 1.21 or higher

### Building

```bash
go build -o markdown_review_daily markdown_review_daily.go
```

Or build for different platforms:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o markdown_review_daily markdown_review_daily.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o markdown_review_daily markdown_review_daily.go

# Windows
GOOS=windows GOARCH=amd64 go build -o markdown_review_daily.exe markdown_review_daily.go
```

## Usage

```bash
./markdown_review_daily [options] [directory]
```

### Options

- `--exact`: Match exact time (hour and minute must match). By default, the tool shows all files scheduled for today regardless of the current time.

### Examples

```bash
# Scan current directory (show all files scheduled for today)
./markdown_review_daily

# Scan specific directory
./markdown_review_daily ./documents

# Show only files scheduled for RIGHT NOW (exact hour and minute)
./markdown_review_daily --exact

# Scan specific directory with exact time matching
./markdown_review_daily --exact /home/user/notes
```

## Output Format

The program outputs results as a markdown list with wiki-style links, perfect for Obsidian:

```
Scanning directory: /Users/you/Documents/notes
Current time: 2026-01-16 15:07
Mode: Daily matching (any file scheduled for today)
----------------------------------------
- [[movies|Movies]]
- [[weekly/review|Weekly Review]]
- [[projects/planning|Project Planning]]
----------------------------------------
Found 3 file(s) scheduled for review today
```

The output uses Obsidian's wikilink format with the vertical bar (`|`) to set custom display text. The links can be directly copied into an Obsidian note!

## Cron Format

The `review_schedule` field should use standard cron format:

```
* * * * *
│ │ │ │ │
│ │ │ │ └── Day of week (0-7, where both 0 and 7 are Sunday)
│ │ │ └──── Month (1-12)
│ │ └────── Day of month (1-31)
│ └──────── Hour (0-23)
└────────── Minute (0-59)
```

## Markdown File Examples

### Daily at 9 AM
```markdown
---
title: Daily Review
review_schedule: "0 9 * * *"
---

# Daily Review Document

This file is scheduled for review every day at 9:00 AM.
```

### Every Monday at 10 AM
```markdown
---
title: Weekly Monday Review
review_schedule: "0 10 * * 1"
---

# Weekly Review Document

This file is scheduled for review every Monday at 10:00 AM.
```

### First day of every month at 8 AM
```markdown
---
review_schedule: "0 8 1 * *"
---

# Monthly Review
```

### Every 15 minutes
```markdown
---
review_schedule: "*/15 * * * *"
---

# Frequent Check
```

### Weekdays at 2 PM
```markdown
---
review_schedule: "0 14 * * 1-5"
---

# Weekday Afternoon Review
```

### Multiple times (8 AM, 12 PM, 5 PM)
```markdown
---
review_schedule: "0 8,12,17 * * *"
---

# Multiple Daily Checks
```

## Features

- **Full cron syntax support**:
  - Wildcards (`*`)
  - Specific values (`5`)
  - Ranges (`1-5`)
  - Steps (`*/15`, `10-20/2`)
  - Lists (`1,15,30`)
- **YAML frontmatter parsing**: Extracts `review_schedule` from markdown files
- **Recursive directory scanning**: Processes all `.md` files in subdirectories
- **Clear output**: Shows matching files with their schedules
- **Error handling**: Warns about invalid cron expressions without stopping

## Output Example

```
Scanning directory: /Users/you/Documents
Current time: 2026-01-15 18:55
Mode: Daily matching (any file scheduled for today)
----------------------------------------
- [[daily_standup|Daily Standup]]
- [[weekly_review|Weekly Review]]
----------------------------------------
Found 2 file(s) scheduled for review today
```

## Cron Expression Examples

| Expression | Description |
|------------|-------------|
| `0 9 * * *` | Daily at 9:00 AM |
| `0 10 * * 1` | Every Monday at 10:00 AM |
| `0 8 1 * *` | First day of month at 8:00 AM |
| `*/15 * * * *` | Every 15 minutes |
| `0 14 * * 1-5` | Weekdays at 2:00 PM |
| `0 8,12,17 * * *` | Daily at 8 AM, 12 PM, and 5 PM |
| `0 9 * * 0` | Every Sunday at 9:00 AM |
| `30 8 * * 1-5` | Weekdays at 8:30 AM |
| `0 */2 * * *` | Every 2 hours |
| `0 9-17 * * 1-5` | Every hour from 9 AM to 5 PM on weekdays |

## Code Structure

The program is organized into several key functions:

- **`ParseCronSchedule`**: Parses a cron expression string into a structured format
- **`MatchesToday`**: Checks if today is a scheduled day (ignores hour/minute)
- **`MatchesTime`**: Checks if a cron schedule matches a given time exactly (includes hour/minute)
- **`matchesCronField`**: Handles matching individual cron fields (supports wildcards, ranges, steps, and lists)
- **`ExtractFrontmatterData`**: Reads YAML frontmatter from markdown files to extract title and review_schedule
- **`ScanDirectory`**: Recursively scans directories for matching files
- **`main`**: Orchestrates the scanning and displays results as markdown links

## Error Handling

The program handles several error conditions gracefully:

- Invalid cron expressions (displays warning and continues)
- Unreadable files (displays warning and continues)
- Missing frontmatter (silently skips file)
- Invalid directory paths (exits with error)

## Performance

The program efficiently processes large directory trees by:
- Using `filepath.Walk` for optimized directory traversal
- Processing files in a single pass
- Minimal memory footprint (processes files one at a time)

## Testing

You can test the program with example markdown files:

1. Create test files with different schedules
2. Run the program at different times
3. Verify that only files matching the current time are reported

## License

This program is provided as-is for use in managing markdown file review schedules.
