package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func getLaunchDir() string {
	usr, err := user.Current()
	if err != nil {
		// Fallback to environment variable or default
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			homeDir = os.Getenv("USERPROFILE") // Windows
		}
		return filepath.Join(homeDir, ".warp", "launch_configurations")
	}
	return filepath.Join(usr.HomeDir, ".warp", "launch_configurations")
}

type ConfigFile struct {
	Name     string
	Path     string
	ModTime  time.Time
	IsTemp   bool
}

func main() {
	fmt.Println("🚀 Warp Launch Config Updater")
	fmt.Println("==============================")
	fmt.Println()

	launchDir := getLaunchDir()
	
	// Find temp files
	tempFiles := findTempFiles(launchDir)
	if len(tempFiles) == 0 {
		fmt.Println("❌ No temp files found! Please save your current session as a temp config first.")
		fmt.Println("   Use Cmd+P → 'Save New Launch Configuration' → name it 'temp-something'")
		return
	}

	// Get the most recent temp file
	sort.Slice(tempFiles, func(i, j int) bool {
		return tempFiles[i].ModTime.After(tempFiles[j].ModTime)
	})
	latestTemp := tempFiles[0]
	fmt.Printf("✓ Found temp config: %s\n", latestTemp.Name)
	fmt.Println()

	// Get all non-temp configs
	allConfigs := findAllConfigs(launchDir)
	nonTempConfigs := make([]ConfigFile, 0)
	for _, config := range allConfigs {
		if !config.IsTemp {
			nonTempConfigs = append(nonTempConfigs, config)
		}
	}

	if len(nonTempConfigs) == 0 {
		fmt.Println("❌ No existing launch configurations found!")
		return
	}

	// Show menu
	fmt.Println("Available launch configurations:")
	for i, config := range nonTempConfigs {
		fmt.Printf("  %d) %s\n", i+1, config.Name)
	}
	fmt.Println()

	// Get user selection
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the number of the config to UPDATE: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	selection, err := strconv.Atoi(input)
	if err != nil || selection < 1 || selection > len(nonTempConfigs) {
		fmt.Println("❌ Invalid selection!")
		return
	}

	targetConfig := nonTempConfigs[selection-1]
	fmt.Printf("Updating: %s\n", targetConfig.Name)

	// Backup existing config
	backupPath := fmt.Sprintf("%s.backup.%s", targetConfig.Path, time.Now().Format("20060102_150405"))
	if err := copyFile(targetConfig.Path, backupPath); err != nil {
		fmt.Printf("❌ Failed to backup: %v\n", err)
		return
	}
	fmt.Printf("✓ Backed up to: %s\n", filepath.Base(backupPath))

	// Update the config
	if err := updateConfig(latestTemp.Path, targetConfig.Path, targetConfig.Name); err != nil {
		fmt.Printf("❌ Failed to update config: %v\n", err)
		return
	}
	fmt.Printf("✅ Successfully updated '%s'!\n", targetConfig.Name)

	// Clean up temp files
	fmt.Println()
	fmt.Println("🧹 Cleaning up temp files...")
	cleanupCount := 0
	for _, tempFile := range tempFiles {
		if err := os.Remove(tempFile.Path); err == nil {
			cleanupCount++
			fmt.Printf("  ✓ Removed %s\n", tempFile.Name)
		}
	}

	if cleanupCount == 0 {
		fmt.Println("  No temp files to clean up")
	} else {
		fmt.Printf("  🎉 Cleaned up %d temp file(s)\n", cleanupCount)
	}

	// Clean up old backup files (older than 7 days)
	fmt.Println()
	fmt.Println("🧹 Cleaning up old backup files...")
	backupCleanupCount := 0
	files, _ := ioutil.ReadDir(launchDir)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	
	for _, file := range files {
		if strings.Contains(file.Name(), ".backup.") && file.ModTime().Before(sevenDaysAgo) {
			backupPath := filepath.Join(launchDir, file.Name())
			if err := os.Remove(backupPath); err == nil {
				backupCleanupCount++
				fmt.Printf("  ✓ Removed old backup: %s\n", file.Name())
			}
		}
	}

	if backupCleanupCount == 0 {
		fmt.Println("  No old backup files to clean up")
	} else {
		fmt.Printf("  🎉 Cleaned up %d old backup file(s)\n", backupCleanupCount)
	}
}

func findTempFiles(launchDir string) []ConfigFile {
	var tempFiles []ConfigFile
	files, _ := ioutil.ReadDir(launchDir)
	
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yaml") && strings.HasPrefix(file.Name(), "temp") && !strings.HasPrefix(file.Name(), ".") {
			tempFiles = append(tempFiles, ConfigFile{
				Name:    strings.TrimSuffix(file.Name(), ".yaml"),
				Path:    filepath.Join(launchDir, file.Name()),
				ModTime: file.ModTime(),
				IsTemp:  true,
			})
		}
	}
	return tempFiles
}

func findAllConfigs(launchDir string) []ConfigFile {
	var configs []ConfigFile
	files, _ := ioutil.ReadDir(launchDir)
	
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yaml") && !strings.HasPrefix(file.Name(), ".") {
			configs = append(configs, ConfigFile{
				Name:    strings.TrimSuffix(file.Name(), ".yaml"),
				Path:    filepath.Join(launchDir, file.Name()),
				ModTime: file.ModTime(),
				IsTemp:  strings.HasPrefix(file.Name(), "temp"),
			})
		}
	}
	return configs
}

func copyFile(src, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, data, 0644)
}

func updateConfig(tempPath, targetPath, targetName string) error {
	data, err := ioutil.ReadFile(tempPath)
	if err != nil {
		return err
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	
	// Update the name field
	for i, line := range lines {
		if strings.HasPrefix(line, "name:") {
			lines[i] = fmt.Sprintf("name: %s", targetName)
			break
		}
	}

	newContent := strings.Join(lines, "\n")
	return ioutil.WriteFile(targetPath, []byte(newContent), 0644)
}
