package routes

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jelinden/blig/app/db"
	"github.com/jelinden/blig/app/domain"
	"github.com/jelinden/blig/app/util"
	"github.com/julienschmidt/httprouter"
	"github.com/microcosm-cc/bluemonday"
	blackfriday "github.com/russross/blackfriday/v2"
	"github.com/teris-io/shortid"
)

var templates = template.Must(template.ParseGlob("public/tmpl/*"))
var p = bluemonday.UGCPolicy()

func init() {
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	p.AllowAttrs("width", "height", "autoplay", "preload", "loop", "controls").OnElements("video")
	p.AllowAttrs("src", "type").OnElements("source")
}

func New(w http.ResponseWriter, r *http.Request) {
	id, err := shortid.Generate()
	if err != nil {
		log.Println("generating id failed", err.Error())
	}
	http.Redirect(w, r, "/admin/post/id/"+id, http.StatusFound)
}

func AdminRoot(w http.ResponseWriter, r *http.Request) {
	blogs := db.GetBlogs()
	sort.Sort(domain.TimeSlice(blogs))
	blog := domain.Blog{BlogName: db.GetBlogName(), BlogPosts: blogs}
	renderTemplateRoot(w, "adminroot", blog)
}

func Index(w http.ResponseWriter, r *http.Request) {
	ps, _ := r.Context().Value("params").(httprouter.Params)
	if ps.ByName("id") != "" {
		blog := db.GetBlogWithID(ps.ByName("id"))
		renderTemplate(w, "post", blog)
	} else {
		renderTemplateWithoutParams(w, "post")
	}
}

func Post(w http.ResponseWriter, r *http.Request) {
	blogPost := sanitizeBlogPost(r)
	if len(blogPost.Title) > 5 {
		db.SaveBlog(blogPost)
	}
	w.Write([]byte(blogPost.Post))
}

func Publish(w http.ResponseWriter, r *http.Request) {
	blogPost := sanitizeBlogPost(r)
	if blogPost.Date.IsZero() {
		blogPost.Date = time.Now().UTC()
	}
	blogPost.Published = true
	if len(blogPost.Title) < 5 {
		w.Write([]byte("Title was not long enough"))
	} else {
		db.SaveBlog(blogPost)
		w.Write([]byte("published"))
	}
}

func FilePost(w http.ResponseWriter, r *http.Request) {
	ps, _ := r.Context().Value("params").(httprouter.Params)
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("file upload failed", err.Error())
		uploadFailed(w, r, err.Error())
		return
	}
	defer file.Close()
	makeDirs(ps)
	filePath := "./public/images/" + ps.ByName("id") + "/" + handler.Filename
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("saving file failed", err.Error())
		uploadFailed(w, r, err.Error())
		return
	}
	defer f.Close()
	io.Copy(f, file)
	w.WriteHeader(200)
	w.Write([]byte("{\"fileName\":\"" + "/static/images/" + ps.ByName("id") + "/" + handler.Filename + "\"}"))
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	ps, _ := r.Context().Value("params").(httprouter.Params)
	id := ps.ByName("id")
	if id != "" {
		db.DeletePost(id)
	}
	http.Redirect(w, r, "/admin/", http.StatusFound)
}

func makeDirs(ps httprouter.Params) {
	if _, err := os.Stat("./public/images"); os.IsNotExist(err) {
		os.Mkdir("./public/images", 0744)
	}
	if _, err := os.Stat("./public/images/" + ps.ByName("id")); os.IsNotExist(err) {
		os.Mkdir("./public/images/"+ps.ByName("id"), 0744)
	}
}

func SaveBlogName(w http.ResponseWriter, r *http.Request) {
	db.SaveBlogName(r.FormValue("blogName"))
	http.Redirect(w, r, "/admin/", http.StatusOK)
}

func uploadFailed(w http.ResponseWriter, r *http.Request, err string) {
	w.WriteHeader(400)
	w.Write([]byte(err))
}

func renderTemplate(w http.ResponseWriter, tmpl string, blogPost interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, blogPost)
	if err != nil {
		log.Println("executing template failed", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateRoot(w http.ResponseWriter, tmpl string, blog domain.Blog) {
	err := templates.ExecuteTemplate(w, tmpl, blog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateWithoutParams(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sanitizeBlogPost(r *http.Request) domain.BlogPost {
	blogTitle := r.FormValue("blogTitle")
	blogText := r.FormValue("blogText")
	blogText = strings.Replace(blogText, "\t", "    ", -1)
	unsafeText := blackfriday.Run([]byte(blogText))
	sanitized := p.SanitizeBytes(unsafeText)

	html := []byte(util.ImgClass(string(sanitized)))
	var id = r.FormValue("blogId")
	if id == "" {
		id, _ = shortid.Generate()
	}
	oldPost := db.GetBlogWithID(id)
	blogPost := domain.BlogPost{
		ID:        id,
		Title:     string(p.SanitizeBytes([]byte(blogTitle))),
		Slug:      util.Slugify(blogTitle),
		Post:      template.HTML(html),
		Markdown:  blogText,
		Modified:  time.Now().UTC(),
		Published: oldPost.Published,
	}
	if !oldPost.Date.IsZero() {
		blogPost.Date = oldPost.Date
	}
	return blogPost
}
