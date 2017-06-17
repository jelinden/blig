package routes

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"strings"

	"github.com/jelinden/blig/app/db"
	"github.com/julienschmidt/httprouter"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-+="
const cookieLength = 30

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var Loggedin string

func Auth(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderTemplateWithoutParams(w, "root")
}

func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderTemplateWithoutParams(w, "login")
}

func AuthHandler(next http.Handler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "params", ps)
		r = r.WithContext(ctx)
		cookie, err := r.Cookie("login")
		if err != nil || !strings.EqualFold(cookie.Value, Loggedin) {
			time.Sleep(2 * time.Second)
			http.Redirect(w, r, "/login?loginfailed=true", http.StatusFound)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func LoginPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := db.GetUserWithID(r.FormValue("username"))
	if checkPasswordHash(r.FormValue("password"), user.Password) {
		Loggedin = stringWithCharset()
		http.SetCookie(w, &http.Cookie{Name: "login", Value: Loggedin})
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login?loginfailed=true", http.StatusFound)
	}
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func stringWithCharset() string {
	b := make([]byte, cookieLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
