# Blig

Blig is a markdown blog maker. Simple, with code syntax highlighting.

```
go get github.com/jelinden/blig
cd $GOPATH/src/github.com/jelinden/blig
go build
``` 

When starting for the first time, you need to give users credentials as parameters.

```
./blig -username username -password password
```

This is the only user the blog will have.

![](https://github.com/jelinden/blig/raw/master/blig.png)



## using

### LevelDB - key/value database

github.com/syndtr/goleveldb/leveldb

### Others

github.com/russross/blackfriday/v2 - markdown processor

github.com/julienschmidt/httprouter - fast router

github.com/microcosm-cc/bluemonday - html sanitizer

github.com/ventu-io/go-shortid - short ids




