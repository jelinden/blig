package routes

import (
	"net/http"
	"sort"

	"github.com/jelinden/blig/app/db"
	"github.com/jelinden/blig/app/domain"
	"github.com/julienschmidt/httprouter"
)

func Root(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	blogs := db.GetPublishedBlogs()
	sort.Sort(domain.TimeSlice(blogs))
	renderTemplateRoot(w, "root",
		domain.Blog{
			BlogName:  db.GetBlogName(),
			BlogPosts: blogs,
		})
}

func Blog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if ps.ByName("id") != "" {
		blog := db.GetBlogWithID(ps.ByName("id"))
		renderTemplate(w, "blog", domain.BlogItem{
			BlogName: db.GetBlogName(),
			BlogItem: blog,
		})
	} else {
		renderTemplateWithoutParams(w, "notfound")
	}
}
