package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gwillem/notion-md-gen/pkg/tomarkdown"

	"github.com/dstotijn/go-notion"
)

func Run(config Config) error {
	if err := os.MkdirAll(config.Markdown.PostSavePath, 0755); err != nil {
		return fmt.Errorf("couldn't create content folder: %s", err)
	}

	// find database page
	client := notion.NewClient(config.Notion.Key)
	q, err := queryDatabase(client, config.Notion)
	if err != nil {
		return fmt.Errorf("❌ Querying Notion database: %s", err)
	}

	// fmt.Println("Found", len(q.Results), "articles.")
	// fetch page children
	for i, page := range q.Results {
		// fmt.Println(page)
		// Get page blocks tree
		fmt.Printf("[%d/%d] %s\n", i+1, len(q.Results), getSlug(page))

		// skip?
		if cb, ok := page.Properties.(notion.DatabasePageProperties); ok {
			if !*cb["Published"].Checkbox {
				fmt.Println("SKIP")
				continue
			}
		}

		blocks, err := queryBlockChildren(client, page.ID)
		if err != nil {
			fmt.Println("Block error:", err)
			continue
		}
		if len(blocks) == 0 {
			fmt.Println("Empty page, skipping")
			continue
		}
		// fmt.Println("Got:", page.Properties.(notion.DatabasePageProperties)["Published"])
		// if config.FilterArticles && !*page.Properties.(notion.DatabasePageProperties)["Published"].Checkbox {
		// 	fmt.Println("Not marked for publication, skipping")
		// 	continue
		// }

		// if prop := page.Properties.(notion.DatabasePageProperties); prop != nil {
		// 	if !*prop["Published"].Checkbox {
		// 		fmt.Println("Skipping, unpublished")
		// 		continue
		// 	}
		// }

		// *page.Properties.(notion.DatabasePageProperties)["Published"].Checkbox = true

		// fmt.Println("Converting", len(blocks), "blocks")
		// Generate content to file
		if err := generate(page, blocks, config.Markdown); err != nil {
			fmt.Println("generate err, skipping:", err)
			continue
		}

	}

	return nil
}

func getTitle(page notion.Page) string {
	for _, x := range []string{"Title", "Name"} {
		if t := tomarkdown.ConvertRichText(page.Properties.(notion.DatabasePageProperties)[x].Title); t != "" {
			return t
		}
	}
	return ""
}

func getSlug(page notion.Page) string {
	shortTitle := page.Properties.(notion.DatabasePageProperties)["ShortTitle"]
	title := tomarkdown.ConvertRichText(shortTitle.RichText)

	if title == "" {
		title = getTitle(page)
	}

	escapedTitle := strings.ReplaceAll(
		strings.ToValidUTF8(
			strings.ToLower(title),
			"",
		),
		" ", "-",
	)

	escapedTitle = regexp.MustCompile(`[^\w\d\-]`).ReplaceAllString(escapedTitle, "")

	if date := prop(page, "Date"); date != "" {
		escapedTitle = date + "-" + escapedTitle
	}
	return escapedTitle
}

func prop(page notion.Page, key string) string {
	props := page.Properties.(notion.DatabasePageProperties)
	// pp.Print(props)
	if props == nil {
		return ""
	}
	if val, ok := props[key]; ok {
		if val.Date != nil {
			return val.Date.Start.Format("2006-01-02")
		}
		return tomarkdown.ConvertRichText(val.RichText)
	}
	return ""
}

func generate(page notion.Page, blocks []notion.Block, config Markdown) error {
	// Create file
	title := getTitle(page)
	if title == "" {
		return fmt.Errorf("empty page, skipping")
	}

	slug := getSlug(page)
	if slug == "" {
		return fmt.Errorf("cannot construct slug, skipping: %s", title)
	}

	fname := filepath.Join(config.PostSavePath, slug+".md")
	fmt.Println(fname)

	f, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("error create file: %s", err)
	}

	// Generate markdown content to the file
	tm := tomarkdown.New()
	tm.ImgSavePath = filepath.Join(config.ImageSavePath, slug)
	tm.ImgVisitPath = filepath.Join(config.ImagePublicLink, slug)
	tm.ContentTemplate = config.Template
	tm.WithFrontMatter(page)
	if config.ShortcodeSyntax != "" {
		tm.EnableExtendedSyntax(config.ShortcodeSyntax)
	}

	return tm.GenerateTo(blocks, f)
}
