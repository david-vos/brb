package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	appName      = "brb"
	appDir       = appName + ".app"
	contentsDir  = appDir + "/Contents"
	macOSDir     = contentsDir + "/MacOS"
	resourcesDir = contentsDir + "/Resources"
)

func main() {
	fmt.Printf("Building %s...\n", appName)

	// Clean previous build
	if err := os.RemoveAll(appDir); err != nil && !os.IsNotExist(err) {
		fatal("Failed to remove old app bundle: %v", err)
	}

	// Create app bundle structure
	if err := os.MkdirAll(macOSDir, 0755); err != nil {
		fatal("Failed to create MacOS directory: %v", err)
	}
	if err := os.MkdirAll(resourcesDir, 0755); err != nil {
		fatal("Failed to create Resources directory: %v", err)
	}

	// Build the Go binary
	fmt.Println("Compiling Go binary...")
	binaryPath := filepath.Join(macOSDir, appName)
	cmd := exec.Command("go", "build", "-o", binaryPath, "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatal("Failed to build Go binary: %v", err)
	}

	// Copy Info.plist
	if err := copyFile("Info.plist", filepath.Join(contentsDir, "Info.plist")); err != nil {
		fatal("Failed to copy Info.plist: %v", err)
	}

	// Copy app icon
	if err := copyFile("appIcon.icns", filepath.Join(resourcesDir, "appIcon.icns")); err != nil {
		fmt.Printf("Warning: Failed to copy app icon: %v\n", err)
		fmt.Println("The app will use the default icon.")
	}

	// Make binary executable
	if err := os.Chmod(binaryPath, 0755); err != nil {
		fatal("Failed to make binary executable: %v", err)
	}

	// Success message
	fmt.Printf("\nApp bundle created: %s\n\n", appDir)

	// Install to ~/Applications
	installApp()

	// Register with Launch Services
	registerApp()

	printFinalInstructions()
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return destFile.Sync()
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(1)
}

func installApp() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fatal("Failed to get home directory: %v", err)
	}

	destPath := filepath.Join(homeDir, "Applications", appDir)

	// Remove existing app if present
	if err := os.RemoveAll(destPath); err != nil && !os.IsNotExist(err) {
		fatal("Failed to remove existing app: %v", err)
	}

	fmt.Printf("Installing to ~/Applications/...\n")

	// Copy the app bundle
	cmd := exec.Command("cp", "-r", appDir, destPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatal("Failed to copy app to ~/Applications: %v", err)
	}

	fmt.Printf("✓ Installed to %s\n\n", destPath)
}

func registerApp() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fatal("Failed to get home directory: %v", err)
	}

	appPath := filepath.Join(homeDir, "Applications", appDir)
	lsregisterPath := "/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister"

	fmt.Println("Registering with Launch Services...")

	cmd := exec.Command(lsregisterPath, "-f", appPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatal("Failed to register with Launch Services: %v", err)
	}

	fmt.Println("✓ Registered with Launch Services\n")
}

func printFinalInstructions() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "~"
	}
	appPath := filepath.Join(homeDir, "Applications", appDir)

	fmt.Println("To set as default browser:")
	fmt.Printf("  1. Open the app: open '%s'\n", appPath)
	fmt.Println("  2. Click 'Set as Default Browser' in the menu (opens System Settings)")
	fmt.Println("  3. In System Settings > Desktop & Dock > Default web browser, select 'Browser Redirect Bar'")
	fmt.Println()

	fmt.Println("To verify registration, check if app appears in default browser list:")
	fmt.Println("  open 'x-apple.systempreferences:com.apple.Desktop-Settings'")
}
