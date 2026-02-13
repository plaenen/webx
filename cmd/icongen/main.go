// borrowed from https://github.com/templui/templui/blob/main/cmd/icongen/main.go
// MIT License

// Copyright (c) 2023 Axel Adrian

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	outputDir     = "./ui/icon/"
	lucideVersion = "0.563.1" // Current Lucide version - update as needed
)

type GitHubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	DownloadUrl string `json:"download_url"`
}

func main() {
	tag := lucideVersion

	// Fetch commit SHA for the version tag
	commitSHA, err := fetchCommitSHA(tag)
	if err != nil {
		panic(fmt.Errorf("failed to fetch commit SHA for %s: %w", tag, err))
	}
	fmt.Printf("Lucide %s → commit %s\n", tag, commitSHA)

	// Fetch the list of files from GitHub (pinned to tag)
	files, err := fetchGitHubContents(tag)
	if err != nil {
		panic(fmt.Errorf("failed to fetch GitHub contents: %w", err))
	}

	// Initialize slice for icon definitions
	var iconDefs []string
	iconDefs = append(iconDefs, "package icon\n")
	iconDefs = append(iconDefs, "// This file is auto generated\n")
	iconDefs = append(iconDefs, fmt.Sprintf("// Using Lucide icons version %s (commit %s)\n", lucideVersion, commitSHA))

	// Create the output directory if it doesn't exist
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("failed to create output directory %s: %w", outputDir, err))
	}

	// Create a wait group to synchronize the goroutines
	var wg sync.WaitGroup
	// Create a channel to limit the number of concurrent downloads
	semaphore := make(chan struct{}, 10) // Allow up to 10 concurrent downloads
	// Create a mutex to protect shared data structures
	var mu sync.Mutex
	iconDataEntries := make(map[string]string) // To store "iconName": "innerSvgContent"

	// Filter SVG files
	var svgFiles []GitHubContent
	for _, file := range files {
		if filepath.Ext(file.Name) == ".svg" {
			svgFiles = append(svgFiles, file)
		}
	}

	// Setup progress tracking
	total := len(svgFiles)
	fmt.Printf("Found %d SVG icons to process\n", total)

	var processed int32
	statusChan := make(chan bool)
	doneChan := make(chan bool)

	// Start progress reporter goroutine
	go func() {
		for {
			select {
			case <-statusChan:
				current := atomic.AddInt32(&processed, 1)
				percent := float64(current) / float64(total) * 100
				fmt.Printf("\rProcessing icons... [%d/%d] %.1f%% complete", current, total, percent)
			case <-doneChan:
				fmt.Println("\rProcessing icons... [100%] Complete!                ")
				return
			}
		}
	}()

	// Process each SVG file
	for _, file := range svgFiles {
		wg.Add(1)
		go func(file GitHubContent) {
			defer wg.Done()
			// Acquire a token from the semaphore
			semaphore <- struct{}{}
			defer func() {
				// Release the token back to the semaphore when done
				<-semaphore
			}()

			name := strings.TrimSuffix(file.Name, ".svg")
			funcName := toPascalCase(name)

			// Add icon definition (thread-safe)
			mu.Lock()
			iconDefs = append(iconDefs, fmt.Sprintf("var %s = Icon(%q)\n", funcName, name))
			mu.Unlock()

			// Download icon content
			contentBytes, err := downloadFile(file.DownloadUrl)
			if err != nil {
				fmt.Printf("\nError downloading %s: %v\n", file.Name, err)
				statusChan <- true
				return
			}

			// Extract inner SVG content
			innerContent := extractSVGContent(string(contentBytes))

			// Store inner content for icon_data.go (thread-safe)
			mu.Lock()
			iconDataEntries[name] = innerContent
			mu.Unlock()

			// Update status
			statusChan <- true
		}(file)
	}

	// Wait for all downloads to complete
	wg.Wait()

	// Signal the progress reporter to stop
	doneChan <- true

	// Write icon_defs.go
	outputFileDefs := filepath.Join(outputDir, "icon_defs.go")
	err = os.WriteFile(outputFileDefs, []byte(strings.Join(iconDefs, "")), 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write icon_defs.go: %w", err))
	}
	fmt.Printf("Generated %s successfully!\n", outputFileDefs)

	// Download and write LICENSE file
	licenseBytes, err := downloadFile(fmt.Sprintf("https://raw.githubusercontent.com/lucide-icons/lucide/%s/LICENSE", tag))
	if err != nil {
		panic(fmt.Errorf("failed to download Lucide LICENSE: %w", err))
	}
	licenseFile := filepath.Join(outputDir, "LICENSE")
	err = os.WriteFile(licenseFile, licenseBytes, 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write LICENSE: %w", err))
	}
	fmt.Printf("Downloaded %s\n", licenseFile)

	// Write icon_data.go
	outputFileData := filepath.Join(outputDir, "icon_data.go")
	var iconDataContent strings.Builder
	iconDataContent.WriteString("package icon\n\n")
	iconDataContent.WriteString("// This file is auto generated\n")
	iconDataContent.WriteString(fmt.Sprintf("// Using Lucide icons version %s (commit %s)\n\n", lucideVersion, commitSHA))
	iconDataContent.WriteString(fmt.Sprintf("const LucideVersion = %q\n", lucideVersion))
	iconDataContent.WriteString(fmt.Sprintf("const LucideCommitSHA = %q\n\n", commitSHA))
	iconDataContent.WriteString("var internalSvgData = map[string]string{\n")
	for name, data := range iconDataEntries {
		// Escape backticks in SVG data for multi-line string literals
		escapedData := strings.ReplaceAll(data, "`", "`+\"`\"+`")
		iconDataContent.WriteString(fmt.Sprintf("\t%q: `%s`,\n", name, escapedData))
	}
	iconDataContent.WriteString("}\n")

	err = os.WriteFile(outputFileData, []byte(iconDataContent.String()), 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write icon_data.go: %w", err))
	}
	fmt.Printf("Generated %s successfully!\n", outputFileData)

	// Write icon.go
	outputFileIconGo := filepath.Join(outputDir, "icon.go")
	iconGoContent := `package icon

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/a-h/templ"
)

// iconContents caches the fully generated SVG strings for icons that have been used,
// keyed by a composite key of name and props to handle different stylings.
var (
	iconContents = make(map[string]string)
	iconMutex    sync.RWMutex
)

// Props defines the properties that can be set for an icon.
type Props struct {
	Size        int
	Color       string
	Fill        string
	Stroke      string
	StrokeWidth string // Stroke Width of Icon, Usage: "2.5"
	Class       string
}

// IconType is a typed icon reference for use in component Props.
//
// Usage:
//
//	type Props struct {
//	    Icon icon.IconType
//	}
//
// Pass an icon:    Props{Icon: icon.Bell}
// Render in templ: @props.Icon() or @props.Icon(icon.Props{Size: 16})
//
// To also accept arbitrary templ.Component icons (e.g. custom SVGs), add a
// second field to your Props and render whichever is set:
//
//	type Props struct {
//	    Icon     icon.IconType   // typed Lucide icon
//	    IconComp templ.Component // arbitrary component fallback
//	}
//
//	// in templ:
//	if props.Icon != nil {
//	    @props.Icon()
//	} else if props.IconComp != nil {
//	    @props.IconComp
//	}
type IconType func(...Props) templ.Component

// Icon returns an IconType for the specified icon name.
func Icon(name string) IconType {
	return func(props ...Props) templ.Component {
		var p Props
		if len(props) > 0 {
			p = props[0]
		}

		// Create a unique key for the cache based on icon name and all relevant props.
		// This ensures different stylings of the same icon are cached separately.
		cacheKey := fmt.Sprintf("%s|s:%d|c:%s|f:%s|sk:%s|sw:%s|cl:%s",
			name, p.Size, p.Color, p.Fill, p.Stroke, p.StrokeWidth, p.Class)

		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			iconMutex.RLock()
			svg, cached := iconContents[cacheKey]
			iconMutex.RUnlock()

			if cached {
				_, err = w.Write([]byte(svg))
				return err
			}

			// Not cached, generate it
			// The actual generation now happens once and is cached.
			generatedSvg, err := generateSVG(name, p) // p (Props) is passed to generateSVG
			if err != nil {
				// Provide more context in the error message
				return fmt.Errorf("failed to generate svg for icon '%s' with props %+v: %w", name, p, err)
			}

			iconMutex.Lock()
			iconContents[cacheKey] = generatedSvg
			iconMutex.Unlock()

			_, err = w.Write([]byte(generatedSvg))
			return err
		})
	}
}

// generateSVG creates an SVG string for the specified icon with the given properties.
// This function is called when an icon-prop combination is not yet in the cache.
func generateSVG(name string, props Props) (string, error) {
	// Get the raw, inner SVG content for the icon name from our internal data map.
	content, err := getIconContent(name) // This now reads from internalSvgData
	if err != nil {
		return "", err // Error from getIconContent already includes icon name
	}

	size := props.Size
	if size <= 0 {
		size = 24 // Default size
	}

	fill := props.Fill
	if fill == "" {
		fill = "none" // Default fill
	}

	stroke := props.Stroke
	if stroke == "" {
		stroke = props.Color // Fallback to Color if Stroke is not set
	}
	if stroke == "" {
		stroke = "currentColor" // Default stroke color
	}

	strokeWidth := props.StrokeWidth
	if strokeWidth == "" {
		strokeWidth = "2" // Default stroke width
	}

	// Construct the final SVG string.
	// The data-lucide attribute helps identify these as Lucide icons if needed.
	return fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%d\" height=\"%d\" viewBox=\"0 0 24 24\" fill=\"%s\" stroke=\"%s\" stroke-width=\"%s\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"%s\" data-lucide=\"icon\">%s</svg>",
		size, size, fill, stroke, strokeWidth, props.Class, content), nil
}

// getIconContent retrieves the raw inner SVG content for a given icon name.
// It reads from the pre-generated internalSvgData map from icon_data.go.
func getIconContent(name string) (string, error) {
	content, exists := internalSvgData[name]
	if !exists {
		return "", fmt.Errorf("icon '%s' not found in internalSvgData map", name)
	}
	return content, nil
}
`
	err = os.WriteFile(outputFileIconGo, []byte(iconGoContent), 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write icon.go: %w", err))
	}
	fmt.Printf("Generated %s successfully!\n", outputFileIconGo)

	fmt.Println("Icon component files generated successfully in " + outputDir)
}

// toPascalCase converts a kebab-case string to PascalCase
func toPascalCase(s string) string {
	words := strings.Split(s, "-")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}

// fetchCommitSHA returns the commit SHA for a given tag.
func fetchCommitSHA(tag string) (string, error) {
	refURL := fmt.Sprintf("https://api.github.com/repos/lucide-icons/lucide/git/ref/tags/%s", tag)
	req, err := http.NewRequest("GET", refURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching ref: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d for tag %s", resp.StatusCode, tag)
	}

	var result struct {
		Object struct {
			SHA string `json:"sha"`
		} `json:"object"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}
	return result.Object.SHA, nil
}

// fetchGitHubContents fetches the list of icon files for a given git ref (tag or branch).
func fetchGitHubContents(ref string) ([]GitHubContent, error) {
	fmt.Printf("Fetching icon tree for %s...\n", ref)

	treeURL := fmt.Sprintf("https://api.github.com/repos/lucide-icons/lucide/git/trees/%s?recursive=1", ref)

	req, err := http.NewRequest("GET", treeURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type TreeItem struct {
		Path string `json:"path"`
		Type string `json:"type"`
	}
	type TreeResponse struct {
		Tree []TreeItem `json:"tree"`
	}

	var treeResp TreeResponse
	if err := json.Unmarshal(body, &treeResp); err != nil {
		return nil, err
	}

	var contents []GitHubContent
	for _, item := range treeResp.Tree {
		if strings.HasPrefix(item.Path, "icons/") && strings.HasSuffix(item.Path, ".svg") {
			filename := filepath.Base(item.Path)
			contents = append(contents, GitHubContent{
				Name:        filename,
				Path:        item.Path,
				DownloadUrl: fmt.Sprintf("https://raw.githubusercontent.com/lucide-icons/lucide/%s/%s", ref, item.Path),
			})
		}
	}

	if len(contents) == 0 {
		return nil, fmt.Errorf("no SVG files found in the icons directory")
	}

	fmt.Printf("Found %d icons in the repository\n", len(contents))
	return contents, nil
}

// downloadFile downloads a file from a URL
func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download returned status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// extractSVGContent removes the outer SVG tags from the icon content.
func extractSVGContent(svgContent string) string {
	start := strings.Index(svgContent, ">") + 1
	end := strings.LastIndex(svgContent, "</svg>")
	if start == -1 || end == -1 || start >= end {
		// Return an empty string or the original content if tags are not found or invalid
		// This prevents panics with malformed SVGs.
		// Log a warning if this happens frequently.
		// fmt.Printf("Warning: Could not extract content from SVG: %s\n", svgContent)
		return ""
	}
	return strings.TrimSpace(svgContent[start:end])
}
