package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	githubOwner = "shoxious"
	githubRepo  = "simracing_telemetry"
)

// version is injected at build time via -ldflags "-X main.version=1.2.3"
var version = "dev"

type ghRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// checkForUpdate fetches the latest GitHub release and offers an auto-update.
// Runs synchronously so the user sees the result before the dashboard starts.
func checkForUpdate() {
	if version == "dev" {
		fmt.Println("► Version:  dev build  (update check skipped)")
		return
	}

	fmt.Printf("► Version:  %s  –  checking GitHub for updates...\n", version)

	apiURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/releases/latest",
		githubOwner, githubRepo,
	)
	client := &http.Client{Timeout: 7 * time.Second}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("User-Agent", "simracing-dashboard/"+version)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("  Update check: offline or unreachable – skipping")
		return
	}
	defer resp.Body.Close()

	// 404 = no releases published yet
	if resp.StatusCode == 404 {
		fmt.Println("  No releases found on GitHub yet")
		return
	}
	if resp.StatusCode != 200 {
		fmt.Printf("  Update check: GitHub returned HTTP %d\n", resp.StatusCode)
		return
	}

	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil || rel.TagName == "" {
		return
	}

	latest := strings.TrimPrefix(rel.TagName, "v")
	current := strings.TrimPrefix(version, "v")

	if !isNewerVersion(latest, current) {
		fmt.Println("  Already up to date ✓")
		return
	}

	// ── Update available ──────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════════╗")
	fmt.Printf( "  ║  Update available:  v%s  →  %s%s║\n",
		current, rel.TagName,
		strings.Repeat(" ", max(1, 28-len(current)-len(rel.TagName))),
	)
	fmt.Println("  ╚══════════════════════════════════════════╝")

	// Find the Windows EXE asset
	var downloadURL, assetName string
	for _, a := range rel.Assets {
		if strings.HasSuffix(strings.ToLower(a.Name), ".exe") {
			downloadURL = a.BrowserDownloadURL
			assetName = a.Name
			break
		}
	}

	if downloadURL == "" {
		fmt.Printf("  No EXE asset found. Download manually:\n  %s\n\n", rel.HTMLURL)
		return
	}

	fmt.Printf("  Asset: %s\n", assetName)
	fmt.Print("  Download and install now? [Y/n] ")

	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))

	if input != "" && input != "y" {
		fmt.Println("  Update skipped. Continuing with current version.\n")
		return
	}

	if err := downloadAndApply(downloadURL, rel.TagName); err != nil {
		fmt.Printf("  Update failed: %v\n  Continuing with current version.\n\n", err)
	}
	// downloadAndApply calls os.Exit(0) on success – execution never reaches here.
}

// downloadAndApply downloads the new EXE, writes a self-replacing batch script,
// launches it detached, then exits so the batch can swap the files.
func downloadAndApply(downloadURL, tagName string) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot determine own path: %w", err)
	}
	exePath, _ = filepath.EvalSymlinks(exePath)
	dir := filepath.Dir(exePath)

	newExe := filepath.Join(dir, "simracing-dashboard-new.exe")

	// ── Download with progress bar ────────────────────────────────────────────
	fmt.Printf("\n  Downloading %s ...\n", tagName)

	dlResp, err := http.Get(downloadURL) //nolint:noctx
	if err != nil {
		return err
	}
	defer dlResp.Body.Close()

	f, err := os.Create(newExe)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", newExe, err)
	}

	total := dlResp.ContentLength
	var done int64
	buf := make([]byte, 64*1024)
	for {
		n, readErr := dlResp.Body.Read(buf)
		if n > 0 {
			if _, wErr := f.Write(buf[:n]); wErr != nil {
				f.Close()
				os.Remove(newExe)
				return wErr
			}
			done += int64(n)
			if total > 0 {
				filled := int(done * 40 / total)
				fmt.Printf("\r  [%s%s] %d%%",
					strings.Repeat("█", filled),
					strings.Repeat("░", 40-filled),
					done*100/total,
				)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			f.Close()
			os.Remove(newExe)
			return readErr
		}
	}
	f.Close()
	fmt.Println()

	// ── Write self-replacing batch script ─────────────────────────────────────
	batPath := filepath.Join(dir, "simracing-update.bat")
	backup := exePath + ".old"

	bat := fmt.Sprintf(`@echo off
:: SimRacing Dashboard – auto-updater
timeout /t 2 /nobreak >nul
move /y "%s" "%s"
move /y "%s" "%s"
del "%s" >nul 2>&1
start "" "%s"
del "%%~f0"
`, exePath, backup, newExe, exePath, backup, exePath)

	if err := os.WriteFile(batPath, []byte(bat), 0o755); err != nil {
		os.Remove(newExe)
		return fmt.Errorf("cannot write update script: %w", err)
	}

	fmt.Println("  ✓ Download complete!")
	fmt.Println("  Applying update – dashboard will restart automatically...")
	fmt.Println()

	// Launch the batch file detached (survives parent exit on Windows)
	cmd := exec.Command("cmd", "/c", "start", "/min", "", batPath)
	cmd.Start() //nolint:errcheck

	time.Sleep(300 * time.Millisecond)
	os.Exit(0)
	return nil
}

// isNewerVersion returns true when latest > current (simple semver x.y.z).
func isNewerVersion(latest, current string) bool {
	lv := parseSemver(latest)
	cv := parseSemver(current)
	for i := range lv {
		if lv[i] > cv[i] {
			return true
		}
		if lv[i] < cv[i] {
			return false
		}
	}
	return false
}

func parseSemver(v string) [3]int {
	parts := strings.SplitN(v, ".", 3)
	var out [3]int
	for i, p := range parts {
		if i >= 3 {
			break
		}
		out[i], _ = strconv.Atoi(strings.TrimSpace(p))
	}
	return out
}
