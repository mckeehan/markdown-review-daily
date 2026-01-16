# Understanding the Two Modes

The program now supports two modes of operation:

## 1. Daily Mode (Default)

Shows all files scheduled for review **today**, regardless of what time it currently is.

```bash
./markdown_review_daily ./docs
```

### Example

If you have a file with `review_schedule: "0 9 * * *"` (daily at 9 AM):
- Running at 8:00 AM → **File IS shown** (it's scheduled for today)
- Running at 9:00 AM → **File IS shown** (it's scheduled for today)
- Running at 5:00 PM → **File IS shown** (it's scheduled for today)

This is useful for getting a morning reminder of all reviews due that day.

## 2. Exact Time Mode

Shows only files scheduled for the **current hour and minute**.

```bash
./markdown_review_daily --exact ./docs
```

### Example

If you have a file with `review_schedule: "0 9 * * *"` (daily at 9 AM):
- Running at 8:00 AM → File is NOT shown (wrong time)
- Running at 9:00 AM → **File IS shown** (exact match)
- Running at 9:15 AM → File is NOT shown (wrong minute)
- Running at 5:00 PM → File is NOT shown (wrong time)

This is useful for:
- Running as a cron job that triggers at specific times
- Getting notifications only when reviews are due
- Automation that should run only at scheduled times

## Practical Examples

### Morning Reminder (Daily Mode)

Run once in the morning to see everything due today:

```bash
#!/bin/bash
# Run at 8 AM daily
0 8 * * * /usr/local/bin/markdown_review_daily ~/documents | mail -s "Today's Reviews" you@example.com
```

Output at 8 AM would show:
```
✓ standup.md (scheduled for 9 AM)
✓ weekly_meeting.md (scheduled for 10 AM on Mondays)
✓ monthly_report.md (scheduled for 8 AM on 1st of month)
```

### Just-In-Time Notifications (Exact Mode)

Run every 15 minutes to notify only when reviews are due:

```bash
#!/bin/bash
# Run every 15 minutes
*/15 * * * * /usr/local/bin/markdown_review_daily --exact ~/documents && notify-send "Review Due!"
```

At 9:00 AM would show only:
```
✓ standup.md (scheduled for 9:00 AM)
```

At 9:15 AM would show:
```
(no files - nothing scheduled for 9:15)
```

## Which Mode Should You Use?

**Use Daily Mode (default) if:**
- You want a once-per-day summary
- You review your tasks in the morning
- You want to see all reviews scheduled for today at once
- You're planning your day

**Use Exact Mode (--exact) if:**
- You want notifications at specific times
- You're automating actions based on schedules
- You run the tool multiple times per day
- You want reminders just when something is due

## Testing Both Modes

Create a test file `test.md`:
```markdown
---
title: Morning Meeting
review_schedule: "0 9 * * *"
---
# Morning Meeting Notes
```

Test daily mode (should show if today matches):
```bash
./markdown_review_daily
```

Test exact mode (should show only at 9:00 AM):
```bash
./markdown_review_daily --exact
```
