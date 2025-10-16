# Warp Launch Config Updater

A command-line tool to easily update existing Warp terminal launch configurations with your current session state.

## Problem

Warp's built-in "Save New Launch Configuration" doesn't allow overwriting existing configurations. This tool solves that by:

1. Auto-detecting your most recent temp launch config
2. Providing an interactive menu to select which config to update
3. Safely backing up your old configuration
4. Updating the target config with your current session
5. Automatically cleaning up temp files and old backups

## Installation

### Option 1: Go Install (Recommended)
```bash
go install github.com/yourusername/warp-launch-config-updater@latest
```

### Option 2: Build from Source
```bash
git clone https://github.com/yourusername/warp-launch-config-updater.git
cd warp-launch-config-updater
go build -o warp-config-updater main.go
sudo mv warp-config-updater /usr/local/bin/
```

### Option 3: Download Binary
Download the latest release for your platform from the [Releases](https://github.com/yourusername/warp-launch-config-updater/releases) page.

## Usage

1. **Save your current Warp session:**
   - Arrange your windows, tabs, and panes how you want them
   - Press `Cmd+P` (or `Ctrl+P` on Linux/Windows)
   - Type "Save New Launch Configuration"
   - Name it anything starting with "temp" (e.g., `temp`, `temp-update`, `temp-vibework`)

2. **Run the updater:**
   ```bash
   warp-config-updater
   ```

3. **Select which config to update:**
   - The tool will show you a numbered list of your existing launch configurations
   - Type the number and press Enter

4. **Done!** ✨
   - Your old config is backed up with a timestamp
   - Your current session is now saved as the selected launch config
   - All temp files are cleaned up automatically

## Example

```bash
$ warp-config-updater

🚀 Warp Launch Config Updater
==============================

✓ Found temp config: temp-vibework

Available launch configurations:
  1) Demand planner
  2) Fiddle 1
  3) Fiddle
  4) Jest
  5) fire-print
  6) vibework

Enter the number of the config to UPDATE: 6
Updating: vibework
✓ Backed up to: vibework.yaml.backup.20250116_143022
✅ Successfully updated 'vibework'!

🧹 Cleaning up temp files...
  ✓ Removed temp-vibework
  🎉 Cleaned up 1 temp file(s)

🧹 Cleaning up old backup files...
  No old backup files to clean up
```

## Features

- ✅ **Auto-detection** of most recent temp launch config
- ✅ **Interactive menu** for selecting target configuration
- ✅ **Automatic backup** of existing configurations
- ✅ **Safe updates** with error handling
- ✅ **Automatic cleanup** of temp files and old backups
- ✅ **Cross-platform** support (macOS, Linux, Windows)

## Requirements

- [Warp Terminal](https://www.warp.dev/)
- Go 1.16+ (for building from source)

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/yourusername/warp-launch-config-updater/issues).
