# Example Usage

## Setting Up Test Files

Create a test directory with some example markdown files:

```bash
mkdir -p test_docs
```

Create `test_docs/daily_standup.md`:
```markdown
---
title: Daily Standup Notes
review_schedule: "0 9 * * *"
tags: [daily, meeting]
---

# Daily Standup

Review this every day at 9:00 AM.
```

Create `test_docs/weekly_planning.md`:
```markdown
---
title: Weekly Planning
review_schedule: "0 10 * * 1"
---

# Weekly Planning Session

Review this every Monday at 10:00 AM.
```

Create `test_docs/monthly_review.md`:
```markdown
---
title: Monthly Business Review
review_schedule: "0 8 1 * *"
author: John Doe
---

# Monthly Business Review

Review this on the first day of each month at 8:00 AM.
```

Create `test_docs/no_schedule.md`:
```markdown
---
title: Random Notes
author: Jane Smith
---

# Random Notes

This file has no review schedule.
```

## Running the Program

```bash
# Build the program
go build -o markdown_review_daily markdown_review_daily.go

# Run on test directory
./markdown_review_daily test_docs
```

## Expected Output

If you run this on a Monday at 10:00 AM:

```
Scanning directory: test_docs
Current time: 2026-01-20 10:00
----------------------------------------
✓ test_docs/weekly_planning.md
  Schedule: 0 10 * * 1
----------------------------------------
Found 1 file(s) scheduled for review today
```

If you run this on January 1st at 8:00 AM:

```
Scanning directory: test_docs
Current time: 2026-01-01 08:00
----------------------------------------
✓ test_docs/monthly_review.md
  Schedule: 0 8 1 * *
----------------------------------------
Found 1 file(s) scheduled for review today
```

## Integration with Other Tools

### Cron Job

You can set up a cron job to run this daily and email results:

```bash
# Add to crontab (crontab -e)
0 9 * * * /path/to/markdown_review_daily /path/to/documents | mail -s "Files Due for Review" you@example.com
```

### Shell Script

Create a wrapper script `daily_review_check.sh`:

```bash
#!/bin/bash

DOCS_DIR="/home/user/documents"
REVIEW_TOOL="/usr/local/bin/markdown_review_daily"

# Run the tool
output=$($REVIEW_TOOL $DOCS_DIR)

# Parse the count
count=$(echo "$output" | grep "Found" | awk '{print $2}')

# Only send notification if files are found
if [ "$count" -gt 0 ]; then
    echo "$output"
    # You could also send a notification here
    # notify-send "Review Reminder" "$count file(s) need review"
fi
```

### Git Hook

Use as a pre-commit hook to remind about overdue reviews:

```bash
#!/bin/bash
# .git/hooks/pre-commit

REVIEW_TOOL="./markdown_review_daily"

if [ -x "$REVIEW_TOOL" ]; then
    output=$($REVIEW_TOOL .)
    count=$(echo "$output" | grep "Found" | awk '{print $2}')
    
    if [ "$count" -gt 0 ]; then
        echo "⚠️  Reminder: $count file(s) are scheduled for review today"
        echo "$output"
    fi
fi
```

## Advanced Examples

### Complex Schedules

**Quarterly Review (first Monday of Jan, Apr, Jul, Oct):**
```yaml
---
review_schedule: "0 9 1 1,4,7,10 1"
---
```
Note: This checks for January, April, July, October on day 1 if it's a Monday.

**Business Hours Check (every hour 9 AM - 5 PM, weekdays):**
```yaml
---
review_schedule: "0 9-17 * * 1-5"
---
```

**Semi-annual Review (Jan 1 and July 1):**
```yaml
---
review_schedule: "0 9 1 1,7 *"
---
```

### Filtering Output

**Get only file paths:**
```bash
./markdown_review_daily docs | grep "^✓" | sed 's/^✓ //'
```

**Count files:**
```bash
./markdown_review_daily docs | grep "^✓" | wc -l
```

**Open files in editor:**
```bash
./markdown_review_daily docs | grep "^✓" | sed 's/^✓ //' | xargs code
```

## Testing

Run the unit tests:

```bash
go test -v
```

Expected output:
```
=== RUN   TestMatchesCronField
--- PASS: TestMatchesCronField (0.00s)
=== RUN   TestParseCronSchedule
--- PASS: TestParseCronSchedule (0.00s)
=== RUN   TestMatchesTime
--- PASS: TestMatchesTime (0.00s)
=== RUN   TestSundayMatching
--- PASS: TestSundayMatching (0.00s)
PASS
ok      markdown-review-daily   0.001s
```
