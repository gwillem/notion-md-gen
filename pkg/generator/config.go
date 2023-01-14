package generator

type Notion struct {
	Key        string
	DatabaseID string `yaml:"databaseId"`
}

type Markdown struct {
	ShortcodeSyntax string `yaml:"shortcodeSyntax"` // hugo,hexo,vuepress
	PostSavePath    string `yaml:"postSavePath"`
	ImageSavePath   string `yaml:"imageSavePath"`
	ImagePublicLink string `yaml:"imagePublicLink"`

	// Optional:
	GroupByMonth bool   `yaml:"groupByMonth,omitempty"`
	Template     string `yaml:"template,omitempty"`
}

type Config struct {
	Notion   `yaml:"notion"`
	Markdown `yaml:"markdown"`
}
