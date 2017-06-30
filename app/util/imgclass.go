package util

import "strings"

func ImgClass(html string) string {
	count := strings.Count(html, "<img src=")
	for i := 0; i <= count; i++ {
		index := strings.Index(html, "<img src=") + 4
		html = html[:index] + " class=\"pure-img\"" + html[index:]
	}
	return html
}
