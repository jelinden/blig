package routes

import (
	"net/http"
	"time"

	"sort"

	"github.com/jelinden/blig/app/db"
	"github.com/jelinden/blig/app/domain"
	"github.com/jelinden/blig/app/sitemap"
	"github.com/julienschmidt/httprouter"
)

const mainURL = "https://www.jelinden.fi"

func Sitemap(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s := sitemap.New()
	now := time.Now().UTC()
	sitemap.Add(s, &sitemap.URL{
		Loc:        mainURL,
		LastMod:    &now,
		ChangeFreq: sitemap.Daily,
	})
	blogs := db.GetPublishedBlogs()
	sort.Sort(domain.TimeSlice(blogs))
	for _, blog := range blogs {
		sitemap.Add(s, &sitemap.URL{
			Loc:        mainURL + "/blog/" + blog.Slug + "/" + blog.ID,
			LastMod:    &blog.Modified,
			ChangeFreq: sitemap.Daily,
		})
	}

	sitemap.Write(w, s)
}
