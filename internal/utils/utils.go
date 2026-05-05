package utils

import (
	"flag"
	"fmt"
	"md2html/internal/config"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func GetFlags(config *config.Config) {
	inputDir := flag.String("input", "", "enter directory path(default: ./content)")
	outputDir := flag.String("output", "", "output directory path(default: ./html)")
	recursive := flag.Bool("recursive", false, "whether to process subdirectories recursively")
	overwrite := flag.Bool("overwrite", false, "overwrite existing files")
	extensions := flag.String("extensions", "", "additional file extensions,separated by commas(e.g.,.txt,.mdx)")

	flag.Usage = HelpInformation

	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "input":
			config.InputDir = *inputDir
		case "output":
			config.OutputDir = *outputDir
		case "recursive":
			config.Recursive = *recursive
		case "overwrite":
			config.Overwrite = *overwrite
		case "extensions":
			config.Extensions = *extensions
		}
	})

	if config.InputDir == "" {
		config.InputDir = "./content"
	}
	if config.OutputDir == "" {
		config.OutputDir = "./html"
	}
}

func HelpInformation() {
	fmt.Println("Markdown to HTML")
	fmt.Println("\nUsage:")
	fmt.Println("  md2html -input <target-path> [options]")
	fmt.Println()

	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()

	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  md2html -input ./docs")
	fmt.Println("  md2html -input ./docs -output ./html -recursive")
	fmt.Println("  md2html -input ./docs -output ./html -overwrite -extensions .txt,.mkd")
}

func RewriteLinks(content string) string {
	// match: ](xxx.md)
	mdLinkRegexp := regexp.MustCompile(`\]\(([^)]+\.md)\)`)

	return mdLinkRegexp.ReplaceAllStringFunc(content, func(match string) string {
		start := strings.Index(match, "(")
		end := strings.Index(match, ")")

		if start == -1 || end == -1 {
			return match
		}

		link := match[start+1 : end]

		if base, ok := strings.CutSuffix(link, ".md"); ok {
			return "](" + base + ".html" + ")"
		}

		return match
	})
}

func LoadConfig() (*config.Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var config config.Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
