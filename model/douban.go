package model

type DoubanBookInfo struct {
	DoubanId    int64             `json:"id,string"`
	Title       string            `json:"title"`
	AltTitle    string            `json:"alt_title"`
	OriginTitle string            `json:"origin_title"`
	Author      []string          `json:"author"`
	Pubdate     string            `json:"pubdate"`
	Binding     string            `json:"binding"`
	Translator  []string          `json:"translator"`
	Catalog     string            `json:"catalog"`
	Pages       string            `json:"pages"`
	Images      map[string]string `json:"images"`
	Publisher   string            `json:"publisher"`
	Isbn10      string            `json:"isbn10"`
	Isbn13      string            `json:"isbn13"`
	AuthorIntro string            `json:"author_intro"`
	Summary     string            `json:"summary"`
	Price       string            `json:"price"`
}
