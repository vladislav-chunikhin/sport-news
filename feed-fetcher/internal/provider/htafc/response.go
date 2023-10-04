package htafc

import "encoding/xml"

type NewListInformation struct {
	XMLName        xml.Name           `xml:"NewListInformation"`
	ClubName       string             `xml:"ClubName"`
	ClubWebsiteURL string             `xml:"ClubWebsiteURL"`
	NewsletterNews NewsletterNewsList `xml:"NewsletterNewsItems"`
}

type NewsletterNewsList struct {
	XMLName        xml.Name             `xml:"NewsletterNewsItems"`
	NewsletterNews []NewsletterNewsItem `xml:"NewsletterNewsItem"`
}

type NewsletterNewsItem struct {
	XMLName           xml.Name `xml:"NewsletterNewsItem"`
	ArticleURL        string   `xml:"ArticleURL"`
	NewsArticleID     int      `xml:"NewsArticleID"`
	PublishDate       string   `xml:"PublishDate"`
	Taxonomies        string   `xml:"Taxonomies"`
	TeaserText        string   `xml:"TeaserText"`
	ThumbnailImageURL string   `xml:"ThumbnailImageURL"`
	Title             string   `xml:"Title"`
	OptaMatchId       string   `xml:"OptaMatchId"`
	LastUpdateDate    string   `xml:"LastUpdateDate"`
	IsPublished       string   `xml:"IsPublished"`
}
