package markdown

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"regexp"

	"github.com/microcosm-cc/bluemonday"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	ErrFileNotFound = fmt.Errorf("file not found")
	ErrFileEmpty    = fmt.Errorf("file is empty")
)

var fileModTimes map[string]int64 = make(map[string]int64)

var parsedFiles map[string]template.HTML = make(map[string]template.HTML)

// p is a bluemonday policy that allows iframes with the class "enclave-object"
// and adds target="_blank" to fully qualified links.
// This enables youtube embeds and other iframes, while keeping the site secure.
var p = bluemonday.
	UGCPolicy().
	AddTargetBlankToFullyQualifiedLinks(true).
	// Allow iframe with class "enclave-object"
	AllowElementsMatching(regexp.MustCompile(`^iframe$`)).
	AllowAttrs("class").Matching(regexp.MustCompile(`\benclave-object\b`)).OnElements("iframe").
	AllowAttrs("src", "width", "height", "allow", "allowfullscreen", "frameborder").
	OnElements("iframe").
	AllowAttrs("role").OnElements("a")

func sanitizeHTML(input []byte) []byte {
	return p.SanitizeBytes(input)
}

func getfileModTime(filename string) (int64, error) {
	info, err := os.Stat(fmt.Sprintf("content/%s.md", filename))
	if err != nil {
		return 0, err
	}
	return info.ModTime().Unix(), nil
}

// RenderFromFile renders a markdown file to HTML.
// Files are read from /content and must not include the .md extension.
func RenderFromFile(filename string) (template.HTML, error) {
	// Open the file at /content/filename.md
	filepath := fmt.Sprintf("content/%s.md", filename)
	modtime, err := getfileModTime(filename)
	if err != nil {
		return "", ErrFileNotFound
	}

	if lastModTime, exists := fileModTimes[filename]; exists && lastModTime >= modtime {
		if lastModTime == modtime {
			return parsedFiles[filename], nil
		}
	}

	file, err := os.ReadFile(filepath)
	if err != nil {
		return "", ErrFileNotFound
	}
	if len(file) == 0 {
		return "", ErrFileEmpty
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Strikethrough,
			extension.Linkify,
			extension.TaskList,
			extension.Typographer,
			enclave.New(
				&enclave.Config{},
			),
		),
		goldmark.WithParserOptions(),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(file, &buf); err != nil {
		return "", err
	} else {
		fileModTimes[filename] = modtime
		parsedFiles[filename] = template.HTML(sanitizeHTML(buf.Bytes()))
		return template.HTML(sanitizeHTML(buf.Bytes())), nil
	}

}
