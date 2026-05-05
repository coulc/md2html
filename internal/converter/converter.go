package converter

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"md2html/internal/config"
	"md2html/internal/utils"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type Converter struct {
	config *config.Config
	tmpl   *template.Template
	policy *bluemonday.Policy
}

func NewConverter(config *config.Config) *Converter {
	tmpl, err := template.ParseFiles(
		config.Template.Layout,
		config.Template.Style,
	)
	if err != nil {
		tmpl = nil
	}

	return &Converter{
		config: config,
		tmpl:   tmpl,
		policy: bluemonday.UGCPolicy(),
	}
}

func (c *Converter) ConvertFile(mdPath string) error {
	mdContent, err := os.ReadFile(mdPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %s: %w", mdPath, err)
	}

	mdContent = []byte(utils.RewriteLinks(string(mdContent)))
	htmlContent := blackfriday.Run(mdContent)

	relPath, err := filepath.Rel(c.config.InputDir, mdPath)
	if err != nil {
		return fmt.Errorf("failed to obtain relative path: %w", err)
	}

	htmlFileName := strings.TrimSuffix(filepath.Base(relPath), filepath.Ext(relPath)) + ".html"
	htmlRelPath := filepath.Join(filepath.Dir(relPath), htmlFileName)
	htmlPath := filepath.Join(c.config.OutputDir, htmlRelPath)

	if !c.config.Overwrite {
		if _, err := os.Stat(htmlPath); err == nil {
			fmt.Printf("skip existing file:  %v\n", htmlPath)
			return nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(htmlPath), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	title := strings.TrimSuffix(filepath.Base(mdPath), filepath.Ext(mdPath))
	fullHTML := c.WrapHTML(string(htmlContent), title)

	if err := os.WriteFile(htmlPath, []byte(fullHTML), 0o644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Printf("Conversion successful! %s -> %s\n", mdPath, htmlPath)

	return nil
}

func (c *Converter) ConvertDirectory() error {
	if _, err := os.Stat(c.config.InputDir); err != nil {
		return fmt.Errorf("input directory not found: %w", err)
	}

	if err := os.MkdirAll(c.config.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", c.config.OutputDir, err)
	}

	extensions := c.GetSupportedExtensions()

	var files []string

	err := filepath.WalkDir(c.config.InputDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if !c.config.Recursive && path != c.config.InputDir {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if slices.Contains(extensions, ext) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("traversing directory failed: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No .md files found.")
		return nil
	}

	fmt.Printf("Found %d Markdown files, start converting.\n", len(files))

	successCount := c.RunConCurrnet(files)

	fmt.Printf("\nConversion completed! Success: %d/%d \n", successCount, len(files))

	return nil
}

func (c *Converter) GetSupportedExtensions() []string {
	defaultExtensions := []string{".md", ".markdown", ".mdown", ".mkd"}

	if c.config.Extensions != "" {
		customExts := strings.Split(c.config.Extensions, ",")
		for i, ext := range customExts {
			customExts[i] = strings.TrimSpace(ext)
			if !strings.HasPrefix(customExts[i], ".") {
				customExts[i] = "." + customExts[i]
			}
		}
		return append(defaultExtensions, customExts...)
	}
	return defaultExtensions
}

func (c *Converter) RunConCurrnet(files []string) int64 {
	var success int64
	workCount := min(8, runtime.NumCPU())

	jobs := make(chan string, len(files))
	var wg sync.WaitGroup

	for range workCount {
		wg.Go(func() {
			for f := range jobs {
				if c.ConvertFile(f) == nil {
					atomic.AddInt64(&success, 1)
				}
			}
		})
	}

	for _, f := range files {
		jobs <- f
	}

	close(jobs)
	wg.Wait()

	return success
}

func (c *Converter) defwrap(content, title string) string {
	template := `<!DOCTYPE html>
	<html lang="zh-CN">
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{TITLE}}</title>
	<style>
	body {
		max-width: 900px;
		margin: 0 auto;
		padding: 20px;
		font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
		line-height: 1.6;
		color: #333;
	}
	h1, h2, h3, h4, h5, h6 {
		margin-top: 24px;
		margin-bottom: 16px;
		font-weight: 600;
		line-height: 1.25;
	}
	h1 { font-size: 2em; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
	h2 { font-size: 1.5em; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
	code {
		padding: 0.2em 0.4em;
		background-color: #f6f8fa;
		border-radius: 3px;
		font-family: "SF Mono", Monaco, Consolas, "Courier New", monospace;
		font-size: 85%;
	}
	pre {
		padding: 16px;
		overflow: auto;
		background-color: #f6f8fa;
		border-radius: 3px;
		line-height: 1.45;
	}
	pre code {
		padding: 0;
		background-color: transparent;
	}
	blockquote {
		padding: 0 1em;
		color: #6a737d;
		border-left: 0.25em solid #dfe2e5;
		margin: 0;
	}
	table {
		border-collapse: collapse;
		width: 100%;
	}
	th, td {
		padding: 6px 13px;
		border: 1px solid #dfe2e5;
	}
	th {
		font-weight: 600;
		background-color: #f6f8fa;
	}
	img {
		max-width: 100%;
		box-sizing: content-box;
	}
	hr {
		height: 0.25em;
		padding: 0;
		margin: 24px 0;
		background-color: #e1e4e8;
		border: 0;
	}
	</style>
	</head>
	<body>
	{{CONTENT}}
	</body>
	</html>`
	result := strings.ReplaceAll(template, "{{TITLE}}", title)
	result = strings.ReplaceAll(result, "{{CONTENT}}", content)
	return result
}

func (c *Converter) WrapHTML(content, title string) string {
	if c.tmpl == nil {
		fmt.Println("Layout file and style file not found, using default layout and style.")
		return c.defwrap(content, title)
	}

	var buffer bytes.Buffer

	safeHTML := c.policy.SanitizeBytes([]byte(content))
	pageData := struct {
		Title   string
		Content template.HTML
	}{
		Title:   title,
		Content: template.HTML(safeHTML),
	}

	if err := c.tmpl.Execute(&buffer, pageData); err != nil {
		fmt.Println("template execute error:", err)
		fmt.Println("using default layout and style.")
		return c.defwrap(content, title)
	}
	return buffer.String()
}
