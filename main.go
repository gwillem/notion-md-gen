package main

import (
	"log"
	"os"

	"github.com/gwillem/notion-md-gen/pkg/generator"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	// Verbose  []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	NotionKey   string `short:"k" description:"Notion API key" env:"NOTION_SECRET" required:"yes"`
	DatabaseID  string `short:"d" description:"Database ID" required:"yes"`
	PostPath    string `long:"post-path" default:"./posts"`
	ImgPath     string `long:"img-path" default:"./images"`
	ImgURL      string `long:"img-url" default:"/images/notion"`
	Development bool   `long:"dev" description:"Dump all articles"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	// fmt.Println("Starting notion export...")

	config := generator.Config{
		Notion: generator.Notion{
			Key:            opts.NotionKey,
			DatabaseID:     opts.DatabaseID,
			FilterArticles: !opts.Development,
			// FilterProp:     "Status",
			// FilterValue:    []string{"Finished", "Published"},
			// PublishedValue: "Published",
		},
		Markdown: generator.Markdown{
			ShortcodeSyntax: "vuepress",
			PostSavePath:    opts.PostPath,
			ImageSavePath:   opts.ImgPath,
			ImagePublicLink: opts.ImgURL,
		},
	}

	if err := generator.Run(config); err != nil {
		log.Fatal(err)
	}

}
