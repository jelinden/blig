package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jelinden/blig/app/config"
	"github.com/jelinden/blig/app/db"
	"github.com/jelinden/blig/app/routes"
	"github.com/julienschmidt/httprouter"
)

func main() {
	configure()
	if !db.CheckUsers() {
		log.Fatal("You need to give proper username and password parameters './blig -h'")
	}
	router := httprouter.New()
	fs := justFilesFilesystem{http.Dir("public/admin")}
	router.ServeFiles("/admin/static/*filepath", fs)
	router.GET("/admin/", routes.AuthHandler(http.HandlerFunc(routes.Root)))
	router.GET("/admin/post/new", routes.AuthHandler(http.HandlerFunc(routes.New)))
	router.GET("/admin/post/id/:id", routes.AuthHandler(http.HandlerFunc(routes.Index)))
	router.POST("/admin/file/post/:id", routes.AuthHandler(http.HandlerFunc(routes.FilePost)))
	router.POST("/admin/push/post/:id", routes.AuthHandler(http.HandlerFunc(routes.Post)))
	router.POST("/admin/push/publish/:id", routes.AuthHandler(http.HandlerFunc(routes.Publish)))
	router.GET("/admin/post/delete/:id", routes.AuthHandler(http.HandlerFunc(routes.DeletePost)))
	router.POST("/admin/save/blogname", routes.AuthHandler(http.HandlerFunc(routes.SaveBlogName)))
	router.GET("/admin/login", routes.Login)
	router.POST("/admin/login", routes.LoginPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func configure() {
	username := flag.String("username", "foo", "a string")
	password := flag.String("password", "bar", "a string")
	flag.Parse()
	config.SetConfig(config.Config{AdminUsername: username, AdminPassword: password})
}

type justFilesFilesystem struct {
	Fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.Fs.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}
