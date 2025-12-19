/*
Copyright © 2024 tienvu461@gmail.com
*/
package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tienvu461/gosak/utils"
)

// autoupdateCmd represents the autoupdate command
var autoupdateCmd = &cobra.Command{
	Use:   "autoupdate",
	Short: "Update gosak to the latest version",
	Long:  `Check for the latest release, download, and provide command to update gosak.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Checking for updates...")
		latest, err := getLatestRelease()
		if err != nil {
			return fmt.Errorf("failed to fetch latest release: %w", err)
		}

		// Remove 'v' prefix if present for comparison logic if needed,
		// but simple string comparison might be enough if following ver.
		// Usually tags are vX.Y.Z
		currentVer := utils.Version
		if currentVer == "" {
			currentVer = "dev"
		}

		if latest.TagName == currentVer {
			fmt.Printf("You are already using the latest version: %s\n", currentVer)
			return nil
		}

		fmt.Printf("New version available: %s (current: %s)\n", latest.TagName, currentVer)

		assetURL, assetName, err := findAssetURL(latest.Assets)
		if err != nil {
			return fmt.Errorf("failed to find compatible asset: %w", err)
		}

		fmt.Printf("Downloading %s...\n", assetName)
		downloadPath := filepath.Join(os.TempDir(), assetName)
		if err := downloadFile(assetURL, downloadPath); err != nil {
			return fmt.Errorf("failed to download file: %w", err)
		}

		// Extract
		fmt.Println("Extracting...")
		extractedPath, err := extractBinary(downloadPath, assetName, latest.TagName)
		if err != nil {
			return fmt.Errorf("failed to extract binary: %w", err)
		}

		exePath := strings.ReplaceAll("$GOPATH/bin/gosak", "$GOPATH", os.Getenv("GOPATH"))
		if exePath == "/bin/gosak" { // fallback if GOPATH is empty
			exePath = filepath.Join(os.Getenv("HOME"), "go", "bin", "gosak")
		}

		fmt.Println("\nUpdate ready! Run the following command to complete the update:")
		fmt.Printf("sudo mv %s %s\n", extractedPath, exePath)
		fmt.Println("\nOr if you don't have sudo access and installed locally:")
		fmt.Printf("mv %s %s\n", extractedPath, exePath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(autoupdateCmd)
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func getLatestRelease() (*Release, error) {
	resp, err := http.Get("https://api.github.com/repos/tienvu461/gosak/releases/latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}

func findAssetURL(assets []Asset) (string, string, error) {
	osName := runtime.GOOS // linux, darwin, windows
	arch := runtime.GOARCH // amd64, arm64, 386

	// Map to goreleaser naming convention
	// {{ .ProjectName }}_{{ title .Os }}_{{ .Arch }}.tar.gz
	// Linux, Darwin, Windows
	// x86_64, i386, arm64, ...

	var expectedOS string
	switch osName {
	case "linux":
		expectedOS = "Linux"
	case "darwin":
		expectedOS = "Darwin"
	case "windows":
		expectedOS = "Windows"
	default:
		return "", "", fmt.Errorf("unsupported OS: %s", osName)
	}

	var expectedArch string
	switch arch {
	case "amd64":
		expectedArch = "x86_64"
	case "386":
		expectedArch = "i386"
	case "arm64":
		expectedArch = "arm64"
	default:
		expectedArch = arch
	}

	// Suffix
	expectedSuffix := ".tar.gz"
	if osName == "windows" {
		expectedSuffix = ".zip"
	}

	// loose matching
	// Look for string containing OS and Arch
	for _, asset := range assets {
		if strings.Contains(asset.Name, expectedOS) &&
			strings.Contains(asset.Name, expectedArch) &&
			strings.HasSuffix(asset.Name, expectedSuffix) {
			return asset.BrowserDownloadURL, asset.Name, nil
		}
	}

	return "", "", fmt.Errorf("no asset found for %s/%s", expectedOS, expectedArch)
}

func downloadFile(url, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractBinary(archivePath, assetName, version string) (string, error) {
	tempDir := os.TempDir()

	if strings.HasSuffix(assetName, ".zip") {
		return extractZip(archivePath, tempDir, version)
	}
	return extractTarGz(archivePath, tempDir, version)
}

func extractTarGz(archivePath, destDir, version string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Find the binary named 'gosak' or 'gosak.exe', possibly inside a folder
		if strings.HasSuffix(header.Name, "gosak") || strings.HasSuffix(header.Name, "gosak.exe") {
			target := filepath.Join(destDir, fmt.Sprintf("gosak_%s", version))
			outFile, err := os.Create(target)
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return "", err
			}
			outFile.Close()
			os.Chmod(target, 0755)
			return target, nil
		}
	}
	return "", fmt.Errorf("gosak binary not found in archive")
}

func extractZip(archivePath, destDir, version string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "gosak.exe") || strings.HasSuffix(f.Name, "gosak") {
			target := filepath.Join(destDir, fmt.Sprintf("gosak_%s", version))
			if runtime.GOOS == "windows" {
				target += ".exe"
			}

			rc, err := f.Open()
			if err != nil {
				return "", err
			}

			outFile, err := os.Create(target)
			if err != nil {
				rc.Close()
				return "", err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()

			if err != nil {
				return "", err
			}

			os.Chmod(target, 0755)
			return target, nil
		}
	}
	return "", fmt.Errorf("gosak binary not found in archive")
}
