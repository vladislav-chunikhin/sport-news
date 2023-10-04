package htafc

import "encoding/xml"

type NewsArticleInformation struct {
	XMLName        xml.Name    `xml:"NewsArticleInformation"`
	ClubName       string      `xml:"ClubName"`
	ClubWebsiteURL string      `xml:"ClubWebsiteURL"`
	NewsArticle    NewsArticle `xml:"NewsArticle"`
}

type NewsArticle struct {
	XMLName           xml.Name `xml:"NewsArticle"`
	ArticleURL        string   `xml:"ArticleURL"`
	NewsArticleID     int      `xml:"NewsArticleID"`
	PublishDate       string   `xml:"PublishDate"`
	Taxonomies        string   `xml:"Taxonomies"`
	TeaserText        string   `xml:"TeaserText"`
	Subtitle          string   `xml:"Subtitle"`
	ThumbnailImageURL string   `xml:"ThumbnailImageURL"`
	Title             string   `xml:"Title"`
	BodyText          string   `xml:"BodyText"`
	GalleryImageURLs  string   `xml:"GalleryImageURLs"`
	VideoURL          string   `xml:"VideoURL"`
	OptaMatchId       string   `xml:"OptaMatchId"`
	LastUpdateDate    string   `xml:"LastUpdateDate"`
	IsPublished       string   `xml:"IsPublished"`
}
