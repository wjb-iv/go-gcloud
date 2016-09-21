package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

// Hello - used in test JSON api
type Hello struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	if tmpl := tc.Get(testTemplate); tmpl != nil {
		w.Header().Set("Content-Type", "text/html")
		model := Model{"Hello World", "Hello"} // Page model
		tmpl.Execute(w, model)                 // merge
	} else {
		log.Fatal("Unable to write response")
		http.Error(w, http.StatusText(500), 500)
	}
}

func helloToo(w http.ResponseWriter, r *http.Request) {
	if tmpl := tc.Get(testTemplate); tmpl != nil {
		w.Header().Set("Content-Type", "text/html")
		model := Model{"The quick brown fox...", "Brown Fox"} // Page model
		tmpl.Execute(w, model)                                // merge
	} else {
		log.Fatalf("Unable to write response")
		http.Error(w, http.StatusText(500), 500)
	}
}

func helloAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	reply := Hello{
		id,
		"Hello",
	}
	if err := json.NewEncoder(w).Encode(reply); err != nil {
		log.Fatalf("Unable to write response: %v", err)
		http.Error(w, http.StatusText(500), 500)
	}
}

func helloAPIAdd(w http.ResponseWriter, r *http.Request) {
	if body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)); err == nil {
		defer r.Body.Close()
		log.Printf("POST Body: %s ", body)
		var hello Hello
		if err := json.Unmarshal(body, &hello); err == nil {

			log.Printf("Do something with entity ID: %s; Message: %s", hello.ID, hello.Message)

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				log.Fatalf("%v", err)
			}
		}
	} else {
		log.Fatalf("Unable to process POST")
		http.Error(w, http.StatusText(500), 500)
	}
}

func main() {
	tc = templateutils.NewCache(bucketName)
	defer tc.Close()
	router := mux.NewRouter().StrictSlash(true)

	// Pages:
	router.Methods("GET").Path("/").Name("hello").Handler(http.HandlerFunc(hello))
	router.Methods("GET").Path("/too").Name("helloToo").Handler(http.HandlerFunc(helloToo))

	//JSON APIs:
	router.Methods("GET").Path("/api/hello/{id}").Name("getHello").Handler(http.HandlerFunc(helloAPI))
	router.Methods("POST").Path("/api/hello").Name("postHello").Handler(http.HandlerFunc(helloAPIAdd))

	log.Fatal(http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router)))
}
