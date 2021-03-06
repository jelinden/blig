package db

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"encoding/json"

	"github.com/jelinden/blig/app/config"
	"github.com/jelinden/blig/app/domain"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var db *leveldb.DB

const userPrefix = "user-"
const blogPrefix = "blog-"
const commonPrefix = "common-"

func init() {
	var err error
	db, err = leveldb.OpenFile("db", nil)
	if err != nil {
		log.Fatal("connection to db failed", err.Error())
	}
}

func CheckUsers() bool {
	users := GetUsers()
	username := *config.Conf.AdminUsername
	passwd := *config.Conf.AdminPassword
	if len(users) > 0 {
		return true
	} else if len(users) == 0 && len(username) > 3 && len(passwd) > 3 {
		SaveUser(domain.User{Username: username, Password: hashPassword(passwd)})
		return true
	}
	return false
}

func GetDB() *leveldb.DB {
	return db
}

func GetUserWithID(id string) domain.User {
	user := domain.User{}
	json.Unmarshal(getWithID(userPrefix+id), &user)
	return user
}

func GetBlogWithID(id string) domain.BlogPost {
	blogPost := domain.BlogPost{}
	json.Unmarshal(getWithID(blogPrefix+id), &blogPost)
	return blogPost
}

func GetUsers() []domain.User {
	iter := db.NewIterator(util.BytesPrefix([]byte(userPrefix)), nil)
	users := []domain.User{}
	for iter.Next() {
		user := domain.User{Username: string(iter.Key()), Password: string(iter.Value())}
		users = append(users, user)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Println(err)
	}
	return users
}

func GetPublishedBlogs() []domain.BlogPost {
	iter := db.NewIterator(util.BytesPrefix([]byte(blogPrefix)), nil)
	blogs := []domain.BlogPost{}
	for iter.Next() {
		blogPost := domain.BlogPost{}
		err := json.Unmarshal(iter.Value(), &blogPost)
		if err != nil {
			log.Println("failed to unmarshal blog post", err.Error())
		}
		if blogPost.Published {
			blogs = append(blogs, blogPost)
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Println(err)
	}
	return blogs
}

func GetBlogs() []domain.BlogPost {
	iter := db.NewIterator(util.BytesPrefix([]byte(blogPrefix)), nil)
	blogs := []domain.BlogPost{}
	for iter.Next() {
		blogPost := domain.BlogPost{}
		err := json.Unmarshal(iter.Value(), &blogPost)
		if err != nil {
			log.Println("failed to unmarshal blog post", err.Error())
		}
		blogs = append(blogs, blogPost)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Println(err)
	}
	return blogs
}

func getWithID(id string) []byte {
	data, err := db.Get([]byte(id), nil)
	if err != nil {
		log.Println("getting with id", id, "failed", err.Error())
		return nil
	}
	return data
}

func SaveBlog(blogPost domain.BlogPost) {
	blogPost.Version++
	j, _ := json.Marshal(blogPost)
	err := db.Put([]byte(blogPrefix+blogPost.ID), j, nil)
	if err != nil {
		log.Println("saving blog post", blogPost.ID, "failed")
	}
}

func SaveUser(user domain.User) {
	j, _ := json.Marshal(user)
	err := db.Put([]byte(userPrefix+user.Username), j, nil)
	if err != nil {
		log.Println("saving blog post", user.Username, "failed")
	}
}

func SaveBlogName(blogName string) {
	err := db.Put([]byte(commonPrefix+"blogName"), []byte(blogName), nil)
	if err != nil {
		log.Println("saving blog name", blogName, "failed")
	}
}

func GetBlogName() string {
	name, err := db.Get([]byte(commonPrefix+"blogName"), nil)
	if err != nil {
		log.Println("getting blog name failed", err.Error())
	}
	return string(name)
}

func DeletePost(id string) {
	err := db.Delete([]byte(blogPrefix+id), nil)
	if err != nil {
		log.Println("deleting", id, "failed", err.Error())
	}
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println("hashing password failed", err.Error())
	}
	return string(bytes)
}
