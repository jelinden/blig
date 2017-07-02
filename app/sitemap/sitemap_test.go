package sitemap

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSitemap(t *testing.T) {
	s := New()
	now := time.Unix(0, 0).UTC()
	Add(s, &URL{
		Loc:        "http://example.com/",
		LastMod:    &now,
		ChangeFreq: Daily,
	})
	Add(s, &URL{
		Loc:        "http://example2.com/",
		LastMod:    &now,
		ChangeFreq: Always,
	})
	buf := new(bytes.Buffer)
	Write(buf, s)
	assert.True(t, strings.EqualFold(buf.String(), testXML), "")
}

const testXML = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>http://example.com/</loc><lastmod>1970-01-01T00:00:00Z</lastmod><changefreq>daily</changefreq></url><url><loc>http://example2.com/</loc><lastmod>1970-01-01T00:00:00Z</lastmod><changefreq>always</changefreq></url></urlset>`
