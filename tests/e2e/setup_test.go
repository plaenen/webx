package e2e_test

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"

	pw "github.com/playwright-community/playwright-go"
)

var (
	baseURL    string
	playwright *pw.Playwright
	browser    pw.Browser
	serverCmd  *exec.Cmd
)

func TestMain(m *testing.M) {
	// Build the showcase binary.
	build := exec.Command("go", "build", "-o", "showcase-test", "./cmd/showcase")
	build.Dir = "../.." // project root relative to tests/e2e/
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		log.Fatalf("build showcase: %v", err)
	}
	defer os.Remove("../../showcase-test")

	// Start the server (binary is in project root, Dir is project root).
	serverCmd = exec.Command("./showcase-test", "serve", "--port", "0")
	serverCmd.Dir = "../.."
	stderr, err := serverCmd.StderrPipe()
	if err != nil {
		log.Fatalf("stderr pipe: %v", err)
	}
	serverCmd.Stdout = os.Stdout

	if err := serverCmd.Start(); err != nil {
		log.Fatalf("start server: %v", err)
	}

	// Parse port from server log output: "server started" address=http://localhost:PORT
	portRe := regexp.MustCompile(`localhost:(\d+)`)
	scanner := bufio.NewScanner(stderr)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if matches := portRe.FindStringSubmatch(line); len(matches) > 1 {
			baseURL = fmt.Sprintf("http://127.0.0.1:%s", matches[1])
			found = true
			break
		}
	}
	if !found {
		serverCmd.Process.Kill()
		log.Fatal("could not determine server port from log output")
	}

	// Wait for server to be reachable.
	for i := 0; i < 50; i++ {
		resp, err := http.Get(baseURL)
		if err == nil {
			resp.Body.Close()
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Install Playwright browsers if needed and launch.
	if err := pw.Install(&pw.RunOptions{Browsers: []string{"chromium"}}); err != nil {
		log.Fatalf("install playwright browsers: %v", err)
	}

	playwright, err = pw.Run()
	if err != nil {
		log.Fatalf("start playwright: %v", err)
	}

	browser, err = playwright.Chromium.Launch()
	if err != nil {
		log.Fatalf("launch chromium: %v", err)
	}

	code := m.Run()

	browser.Close()
	playwright.Stop()
	serverCmd.Process.Kill()
	serverCmd.Wait()
	os.Exit(code)
}

// newPage creates a new browser page for a test.
func newPage(t *testing.T) pw.Page {
	t.Helper()
	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("new page: %v", err)
	}
	t.Cleanup(func() { page.Close() })
	return page
}
