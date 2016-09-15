package main

import (
	"net/http"

	"github.com/wjb-iv/go-gcloud/templateutils"
)

var (
	testTemplate = "example_template"
	bucketName   = "org_wjb_test_bucket"
	tc           *templateutils.Cache
)

// Model - used in test templates
type Model struct {
	Message string
	Title   string
}

func hello(w http.ResponseWriter, r *http.Request) {
	if tmpl := tc.Get(testTemplate); tmpl != nil {
		w.Header().Set("Content-Type", "text/html")
		model := Model{"Hello World", "Hello"} // Page model
		tmpl.Execute(w, model)                 // merge
	} else {
		http.Error(w, http.StatusText(500), 500)
	}
}

func helloToo(w http.ResponseWriter, r *http.Request) {
	if tmpl := tc.Get(testTemplate); tmpl != nil {
		w.Header().Set("Content-Type", "text/html")
		model := Model{"The quick brown fox...", "Brown Fox"} // Page model
		tmpl.Execute(w, model)                                // merge
	} else {
		http.Error(w, http.StatusText(500), 500)
	}
}

func main() {
	tc = templateutils.NewCache(bucketName)
	defer tc.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/too", helloToo)
	http.ListenAndServe(":8080", mux)
}
