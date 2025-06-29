# Go Auto Backup Tool

A cross-platform (macOS and Windows) Go application to periodically back up chosen directories to an external HDD. Backups are saved by date, and configuration is prompted on first run and saved for future use.

## Features
- Choose source directories
- Choose target directory (e.g., external HDD)
- Choose sync time (daily at HH:MM or interval like 1m, 10m, 1h, etc.)
- Immediate backup on app start
- Auto-refreshes configuration before each backup (edit the config file and changes are picked up automatically)
- Saves configuration for future runs
- Backups stored in date-based folders
- Cross-platform: macOS and Windows
- Native GUI for configuration (Fyne)
- **Version info is printed in logs and on startup**

## Usage

1. Build the app:
   ```sh
   go build -o go-auto-backup ./main.go
   ```
2. Run the app:
   ```sh
   ./go-auto-backup
   ```
   - On first run, a GUI window will prompt you to configure sources, target, and schedule.
   - On subsequent runs, backup will run automatically per schedule, and the config can be edited at any time.
   - The app prints its version on startup and in the logs for every backup run.

## Configuration
- The configuration is saved in your home directory as `.go-dir-backup-config.json`.
- **Sync Time**:
  - Use `HH:MM` (e.g., `02:00`) for once-a-day backups at a specific time.
  - Use Go duration strings (e.g., `1m`, `10m`, `1h`) for interval backups.
- The app performs an immediate backup on start, then continues per schedule.
- The config is reloaded before every backup, so you can change settings without restarting the app.

## Notes
- Ensure the target directory (external HDD) is connected before scheduled backups.

## CI/CD
- GitHub Actions workflow lints and builds the app for macOS and Windows.
- Built binaries are uploaded as artifacts for each run.
- Version info is injected into the binary from the branch or tag name.
