package sitemap

import (
	"encoding/xml"
	"io"
	"log"
	"time"
)

type changeFreq string

const (
	Always  changeFreq = "always"
	Hourly  changeFreq = "hourly"
	Daily   changeFreq = "daily"
	Weekly  changeFreq = "weekly"
	Monthly changeFreq = "monthly"
	Yearly  changeFreq = "yearly"
	Never   changeFreq = "never"
)

type URL struct {
	Loc        string     `xml:"loc"`
	LastMod    *time.Time `xml:"lastmod,omitempty"`
	ChangeFreq changeFreq `xml:"changefreq,omitempty"`
	Priority   float32    `xml:"priority,omitempty"`
}

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`

	URLs []*URL `xml:"url"`
}

func New() *Sitemap {
	return &Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]*URL, 0),
	}
}

func Add(s *Sitemap, u *URL) {
	s.URLs = append(s.URLs, u)
}

func Write(w io.Writer, s *Sitemap) {
	sitemap, err := xml.Marshal(s)
	if err != nil {
		log.Println("Marshallin sitemap failed", err.Error())
	}
	w.Write(append([]byte(xml.Header)[:], sitemap[:]...))
}
