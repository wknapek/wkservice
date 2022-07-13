package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var conf = Configuration{}

func main() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&conf)
	if err != nil {
		log.Fatalln(err.Error())
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", ping)
	http.ListenAndServe(":3000", r)

}

type Configuration struct {
	versionBuild string
	env          string
}

type answer struct {
	Args struct {
		Foo1 string `json:"foo1"`
		Foo2 string `json:"foo2"`
	} `json:"args"`
	Headers struct {
		XForwardedProto string `json:"x-forwarded-proto"`
		Host            string `json:"host"`
		Accept          string `json:"accept"`
		AcceptEncoding  string `json:"accept-encoding"`
		CacheControl    string `json:"cache-control"`
		Cookie          string `json:"cookie"`
		PostmanToken    string `json:"postman-token"`
		UserAgent       string `json:"user-agent"`
		XForwardedPort  string `json:"x-forwarded-port"`
	} `json:"headers"`
	URL       string    `json:"url"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	env       string    `json:"env"`
}

func ping(w http.ResponseWriter, r *http.Request) {
	_, err := r.URL.Query()["message"]
	if !err {
		http.Error(w, "param message is missing", http.StatusBadRequest)
	}
	ans := answer{}
	ans.Timestamp = time.Now()
	ans.Headers.XForwardedProto = "https"
	ans.Headers.Host = "wkservice"
	ans.Headers.Accept = "*/*"
	ans.Headers.AcceptEncoding = "gzip, deflate"
	ans.Headers.CacheControl = "no-cache"
	ans.Headers.Cookie = ""
	ans.Headers.PostmanToken = ""
	ans.Headers.UserAgent = ""
	ans.Headers.XForwardedPort = "3000"
	ans.URL = "/"
	ans.env = conf.env
	ans.Version = conf.versionBuild
	w.Write([]byte(fmt.Sprintf("%v", ans)))
}
