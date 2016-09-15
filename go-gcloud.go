package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
)

var (
	testTemplate  = "example_template.html"
	storageClient *storage.Client
	ctx           = context.Background()
	bucketName    = "org_wjb_test_bucket"
)

// Model - used in test templates
type Model struct {
	Message string
	Title   string
}

func hello(w http.ResponseWriter, r *http.Request) {
	if tmpl := getObject(testTemplate); tmpl != "" {
		w.Header().Set("Content-Type", "text/html")
		t := template.New("test")              // Create a template
		t, _ = t.Parse(tmpl)                   // Parse template string
		model := Model{"Hello World", "Hello"} // Page model
		t.Execute(w, model)                    // merge
	} else {
		http.Error(w, http.StatusText(500), 500)
	}
}

func main() {
	storageClient = initStorageClient()
	if storageClient != nil {
		defer storageClient.Close()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	http.ListenAndServe(":8080", mux)
}

func initStorageClient() *storage.Client {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("Cannot create storage client")
	} else {
		return client
	}
	return nil
}

func getObject(objName string) string {
	if rc, err := storageClient.Bucket(bucketName).Object(objName).NewReader(ctx); err == nil {
		if slurp, err := ioutil.ReadAll(rc); err == nil {
			defer rc.Close()
			return fmt.Sprintf("%s", slurp)
		}
		log.Fatalf("Cannot read object - %s", err)
	}
	return ""
}
