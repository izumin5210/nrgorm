# NrGorm
[![Build Status](https://travis-ci.org/izumin5210/nrgorm.svg?branch=master)](https://travis-ci.org/izumin5210/nrgorm)
[![GoDoc](https://godoc.org/github.com/izumin5210/nrgorm?status.svg)](https://godoc.org/github.com/izumin5210/nrgorm)
[![Go project version](https://badge.fury.io/go/github.com%2Fizumin5210%2Fnrgorm.svg)](https://badge.fury.io/go/github.com%2Fizumin5210%2Fnrgorm)
[![license](https://img.shields.io/github/license/izumin5210/nrgorm.svg)](./LICENSE)

## Example
### Simple app

```go
var db *gorm.DB
var nrapp newrelic.Application

func main() {
	db = initDB()
	defer db.Close()

	nrapp = initNrApp()

	http.HandleFunc("/", handler)
	http.ListenAndServve(":8080", nil)
}

func initNrApp() newrelic.Application {
	cfg := newrelic.NewConfig("Example App", "__YOUR_NEW_RELIC_LICENSE_KEY__")

	app, err := newrelic.NewApplication(cfg)
	if err != nil {
		panic(err)
	}

	return app
}

func initDB() *gorm.DB {
	db, err := gorm.Open("postgres", "postgres://postgres:@localhost/maindb?sslmode=disable")
	if err != nil {
		panic(err)
	}

	// Register callback funcs to *gorm.DB
	nrgorm.RegsiterCallbacks(db, "maindb")

	return db
}

func handler(w http.ResponseWriter, r *http.Request) {
	txn := nrapp.StartTransaction("/", w, r)
	defer txn.End()

	// Set newrelic.Transaction to *gorm.DB
	db = nrgorm.Wrap(txn, db)

	// do something...
}
```

### With [Gin Web Framework](https://github.com/gin-gonic/gin) and [nrgin](https://github.com/newrelic/go-agent/tree/master/_integrations/nrgin/v1)

```go
var db *gorm.DB

func WrapDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		db = nrgorm.Wrap(db, nrgin.Transcation(c))
	}
}

func main() {
	r := gin.Default()
	db = initDB()
	r.Use(nrgin.Middleware(initNrApp()))
	r.Use(WrapDB())
	r.GET("/", func(c *gin.Context) {
		// do something...
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
```
