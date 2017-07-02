package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jelinden/blig/app/db"
	"github.com/stretchr/testify/assert"
)

var client = &http.Client{}

const baseURL = "http://localhost" + port
const requestTooSlowInSeconds = 0.1

func init() {
	go main()
	waitForConnection()
}

func testPageLoad(t *testing.T, url string) []byte {
	beginning := time.Now()
	req, _ := http.NewRequest("GET", baseURL+url, nil)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.True(t, resp.StatusCode == 200)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	took := time.Now().Sub(beginning).Seconds()
	log.Println("statuscode", resp.StatusCode, "OK", "request took", took, "seconds", url)
	assert.True(t, took < requestTooSlowInSeconds, "request took too long, url: "+url)
	return body
}

func TestIndex(t *testing.T) {
	body := testPageLoad(t, "/")
	assert.True(t, strings.Contains(string(body), "title"))
}

func TestRSS(t *testing.T) {
	body := testPageLoad(t, "/rss")
	assert.True(t, strings.Contains(string(body), "xml"))
}

func TestSitemap(t *testing.T) {
	body := testPageLoad(t, "/sitemap.xml")
	assert.True(t, strings.Contains(string(body), "xml"))
}

func TestLoad(t *testing.T) {
	blogs := db.GetBlogs()
	beginning := time.Now()
	var i = 0
	const requests = 2000
	for i < requests {
		rand.Seed(time.Now().UnixNano())
		blog := blogs[rand.Intn(len(blogs))]
		body := testPageLoad(t, "/blog/"+blog.Slug+"/"+blog.ID)
		assert.True(t, strings.Contains(string(body), "title"))
		i++
	}
	took := time.Now().Sub(beginning).Seconds()
	amountRequests := math.Floor(requests * (1 / took))
	fmt.Println("\n\n ******** Made", requests, "requests in", took, "seconds (", amountRequests, "req/s ) ********\n")
}

var initConnection = 0

func waitForConnection() {
	maxInitConnection := 300
	if initConnection < maxInitConnection {
		resp, err := http.Get(baseURL + "/health")
		if err != nil || resp.StatusCode != 200 {
			time.Sleep(1000 * time.Millisecond)
			initConnection++
			waitForConnection()
		} else {
			log.Println("Health OK")
		}
	}
}
