# md2html - Markdown 转 HTML 工具

一个轻量、快速的 Markdown 静态转换工具，使用 Go 语言开发。

## 项目简介

md2html 是一个命令行工具，用于批量将 Markdown 文件转换为 HTML。支持目录遍历、链接重写、自定义模板和并发处理，适合构建静态网站、文档站点或个人博客。

## 功能特性

- **Markdown 转 HTML**：使用 `blackfriday/v2` 进行快速可靠的 Markdown 解析
- **链接重写**：自动将相对路径的 `.md` 链接转换为 `.html` 链接
- **批量处理**：一次性转换目录下的所有 Markdown 文件
- **递归遍历**：可选择递归处理子目录
- **文件名冲突预防**：处理同名但不同扩展名的文件（如 `test.md` 和 `test.txt`）
- **并发转换**：使用 goroutine 并行处理（最多 8 个工作协程）
- **自定义模板**：通过 Go 模板支持自定义 HTML 布局和 CSS 样式
- **HTML 安全处理**：使用 `bluemonday` 净化输出 HTML 确保安全性
- **配置文件**：支持 YAML 配置文件持久化设置

## 安装方式

### 环境要求

- Go 1.21 或更高版本

### 从源码编译

```bash
# 克隆仓库
git clone <仓库地址>
cd md2html

# 安装依赖
go mod tidy

# 编译
go build -o md2html

# （可选）安装到 PATH
go install
```

## 使用示例

### 基本用法

将默认 `./content` 目录下的所有 Markdown 文件转换到 `./html` 目录：

```bash
md2html
```

### 指定输入和输出目录

```bash
md2html -input ./docs -output ./public
```

### 递归处理子目录

```bash
md2html -input ./docs -output ./public -recursive
```

### 覆盖已存在的文件

```bash
md2html -input ./docs -output ./public -overwrite
```

### 支持额外的文件扩展名

```bash
md2html -input ./docs -extensions .txt,.mdx,.markdown
```

### 包含所有选项的完整示例

```bash
md2html -input ./content -output ./html -recursive -overwrite -extensions .txt,.mkd
```

### 查看帮助信息

```bash
md2html -h
# 或
md2html -help
```

## 项目结构说明

```
md2html/
├── main.go           # 主程序入口
├── config.yaml       # 配置文件（可选）
├── go.mod            # Go 模块文件
├── go.sum            # 依赖校验和
├── templates/        # 自定义模板目录（可选）
│   ├── layout.html   # HTML 布局模板
│   └── style.css     # CSS 样式模板
├── content/          # 默认输入目录（存放 Markdown 文件）
└── html/             # 默认输出目录（生成的 HTML 文件）
```

## 配置说明

### 命令行选项

| 选项 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-input` | string | `./content` | 输入目录路径，存放 Markdown 文件 |
| `-output` | string | `./html` | 输出目录路径，存放生成的 HTML 文件 |
| `-recursive` | bool | `false` | 递归处理子目录 |
| `-overwrite` | bool | `false` | 覆盖已存在的 HTML 文件 |
| `-extensions` | string | (空) | 额外的文件扩展名，逗号分隔 |

### 配置文件 (config.yaml)

在项目根目录创建 `config.yaml` 文件进行持久化配置：

```yaml
input: ./content
output: ./html
recursive: false
overwrite: false
extensions: .txt,.mdx
template:
  layout: ./templates/layout.html
  style: ./templates/style.css
```

**注意**：命令行选项优先级高于配置文件设置。

### 默认支持的扩展名

- `.md`
- `.markdown`
- `.mdown`
- `.mkd`

## 自定义模板

### 布局模板 (layout.html)

```html
<!DOCTYPE html>
<html lang="zh-CN">
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

### 样式模板 (style.css)

```css
body {
    max-width: 900px;
    margin: 0 auto;
    padding: 20px;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    line-height: 1.6;
    color: #333;
}

/* 在此添加自定义样式 */
```

## 依赖项

- `github.com/russross/blackfriday/v2` - Markdown 解析器
- `github.com/microcosm-cc/bluemonday` - HTML 安全净化
- `gopkg.in/yaml.v3` - YAML 配置解析器

## 许可证

本项目为开源项目。详情请参考许可证文件。
