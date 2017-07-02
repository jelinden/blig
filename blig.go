package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jelinden/blig/app/config"
	"github.com/jelinden/blig/app/db"
	"github.com/jelinden/blig/app/routes"
	"github.com/jelinden/blig/app/util"
	"github.com/julienschmidt/httprouter"
)

func main() {
	configure()
	if !db.CheckUsers() {
		log.Fatal("You need to give proper username and password parameters './blig -h'")
	}
	router := httprouter.New()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true
	fs := util.JustFilesFilesystem{Fs: http.Dir("public/admin")}
	router.ServeFiles("/admin/static/*filepath", fs)
	fsStatic := util.JustFilesFilesystem{Fs: http.Dir("public/")}
	router.Handler("GET", "/static/*filepath", http.StripPrefix("/static", util.GH(http.FileServer(fsStatic))))
	router.GET("/admin/", routes.AuthHandler(http.HandlerFunc(routes.AdminRoot)))
	router.GET("/admin/post/new", routes.AuthHandler(http.HandlerFunc(routes.New)))
	router.GET("/admin/post/id/:id", routes.AuthHandler(http.HandlerFunc(routes.Index)))
	router.POST("/admin/file/post/:id", routes.AuthHandler(http.HandlerFunc(routes.FilePost)))
	router.POST("/admin/push/post/:id", routes.AuthHandler(http.HandlerFunc(routes.Post)))
	router.POST("/admin/push/publish/:id", routes.AuthHandler(http.HandlerFunc(routes.Publish)))
	router.GET("/admin/post/delete/:id", routes.AuthHandler(http.HandlerFunc(routes.DeletePost)))
	router.POST("/admin/save/blogname", routes.AuthHandler(http.HandlerFunc(routes.SaveBlogName)))
	router.GET("/admin/login", routes.Login)
	router.POST("/admin/login", routes.LoginPost)

	router.GET("/", util.MakeGzipHandler(routes.Root))
	router.GET("/blog/:slug/:id", util.MakeGzipHandler(routes.Blog))

	router.GET("/rss", util.MakeGzipHandler(routes.RSS))
	router.GET("/sitemap.xml", routes.Sitemap)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func configure() {
	username := flag.String("username", "foo", "a string")
	password := flag.String("password", "bar", "a string")
	flag.Parse()
	config.SetConfig(config.Config{AdminUsername: username, AdminPassword: password})
}
