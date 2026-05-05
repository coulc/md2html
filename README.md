[中文](./readme_zh.md)
# md2html - Markdown to HTML Converter

A lightweight, fast Markdown to HTML static conversion tool written in Go.

## Project Introduction

md2html is a command-line tool designed for batch conversion of Markdown files to HTML. It supports directory traversal, link rewriting, custom templates, and concurrent processing, making it ideal for building static websites, documentation sites, or personal blogs.

## Features

- **Markdown to HTML Conversion**: Uses `blackfriday/v2` for fast and reliable Markdown parsing
- **Link Rewriting**: Automatically converts relative `.md` links to `.html` links
- **Batch Processing**: Converts all Markdown files in a directory in one go
- **Recursive Traversal**: Optionally processes subdirectories recursively
- **File Conflict Prevention**: Handles same-name files with different extensions (e.g., `test.md` vs `test.txt`)
- **Concurrent Conversion**: Uses goroutines for parallel processing (up to 8 workers)
- **Custom Templates**: Supports custom HTML layout and CSS styles via Go templates
- **HTML Sanitization**: Uses `bluemonday` to sanitize output HTML for security
- **Configuration File**: Supports YAML configuration file for persistent settings

## Installation

### Prerequisites

- Go 1.21 or higher

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd md2html

# Install dependencies
go mod tidy

# Build
go build -o md2html

# (Optional) Install to PATH
go install
```

## Usage Examples

### Basic Usage

Convert all Markdown files in the default `./content` directory to `./html`:

```bash
md2html
```

### Specify Input and Output Directories

```bash
md2html -input ./docs -output ./public
```

### Process Subdirectories Recursively

```bash
md2html -input ./docs -output ./public -recursive
```

### Overwrite Existing Files

```bash
md2html -input ./docs -output ./public -overwrite
```

### Support Additional File Extensions

```bash
md2html -input ./docs -extensions .txt,.mdx,.markdown
```

### Full Example with All Options

```bash
md2html -input ./content -output ./html -recursive -overwrite -extensions .txt,.mkd
```

### Help Information

```bash
md2html -h
# or
md2html -help
```

## Project Structure

```
md2html/
├── main.go           # Main program entry
├── config.yaml       # Configuration file (optional)
├── go.mod            # Go module file
├── go.sum            # Dependencies checksum
├── templates/        # Custom templates directory (optional)
│   ├── layout.html   # HTML layout template
│   └── style.tmpl    # CSS styles template
├── content/          # Default input directory (Markdown files)
└── html/             # Default output directory (generated HTML)
```

## Configuration

### Command Line Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `-input` | string | `./content` | Input directory path containing Markdown files |
| `-output` | string | `./html` | Output directory path for generated HTML files |
| `-recursive` | bool | `false` | Process subdirectories recursively |
| `-overwrite` | bool | `false` | Overwrite existing HTML files |
| `-extensions` | string | (empty) | Additional file extensions, comma-separated |

### Configuration File (config.yaml)

Create a `config.yaml` file in the project root for persistent configuration:

```yaml
input: ./content
output: ./html
recursive: false
overwrite: false
extensions: .txt,.mdx
template:
  layout: ./templates/layout.html
  style: ./templates/style.tmpl
```

**Note**: Command line options override configuration file settings.

### Default Supported Extensions

- `.md`
- `.markdown`
- `.mdown`
- `.mkd`

## Custom Templates

### Layout Template (layout.html)

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        {{.Content}}
    </style>
</head>
<body>
    {{.Content}}
</body>
</html>
```

### Style Template (style.tmpl)

```css
body {
    max-width: 900px;
    margin: 0 auto;
    padding: 20px;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    line-height: 1.6;
    color: #333;
}

/* Add your custom styles here */
```

## Dependencies

- `github.com/russross/blackfriday/v2` - Markdown parser
- `github.com/microcosm-cc/bluemonday` - HTML sanitizer
- `gopkg.in/yaml.v3` - YAML configuration parser

## License

This project is open source. Please refer to the license file for details.
