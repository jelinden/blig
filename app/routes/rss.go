package routes

import (
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/feeds"
	"github.com/jelinden/blig/app/db"
	"github.com/jelinden/blig/app/domain"
	"github.com/julienschmidt/httprouter"
)

func RSS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	feed := &feeds.Feed{
		Title:       db.GetBlogName(),
		Link:        &feeds.Link{Href: "https://www.jelinden.fi"},
		Description: "Programming blog",
		Created:     time.Now(),
	}

	blogs := db.GetPublishedBlogs()
	sort.Sort(domain.TimeSlice(blogs))
	feed.Items = []*feeds.Item{}
	for i, blog := range blogs {
		if i < 5 {
			feed.Items = append(feed.Items, &feeds.Item{
				Title:       blog.Title,
				Link:        &feeds.Link{Href: "https://www.jelinden.fi/blog/" + blog.Slug + "/" + blog.ID},
				Description: blog.Title,
				Created:     blog.Date,
				Id:          blog.ID,
			})
		}
	}
	rss, err := feed.ToRss()
	if err != nil {
		log.Println("making rss failed", err)
	}
	w.Header().Set("Content-Type", "application/rss+xml")
	w.Write([]byte(rss))
}
