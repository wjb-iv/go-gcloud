package templateutils

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"sync"

	"cloud.google.com/go/storage"
)

// Private, synchronized map to store cached templates ---
type syncMap struct {
	lock sync.RWMutex
	data map[string]*template.Template
}

func (sc *syncMap) get(key string) (*template.Template, bool) {
	sc.lock.RLock()
	defer sc.lock.RUnlock()
	value, ok := sc.data[key]
	return value, ok
}

func (sc *syncMap) set(key string, value *template.Template) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	sc.data[key] = value
} // ---- End of private synchronized map object ----

// Cache - the main interface for working with the cache
type Cache struct {
	cacheData     *syncMap
	storageClient *storage.Client
	bucketName    string
	ctx           context.Context
}

// NewCache - creates a reference to a cache to hold onto
func NewCache(templateBucket string) *Cache {
	return &Cache{
		&syncMap{data: make(map[string]*template.Template)},
		initStorageClient(),
		templateBucket,
		context.Background(),
	}
}

// Get - retrieves the named template from cache, or if not found
//       attempts to load it from a Google storage bucket
func (c *Cache) Get(templateName string) *template.Template {
	log.Printf("Looking for template name: %s...", templateName)
	if val, ok := c.cacheData.get(templateName); ok {
		log.Printf("Found template: %s in template cache.", templateName)
		return val
	}
	log.Printf("Template: %s NOT in template cache, loading from bucket: %s", templateName, c.bucketName)
	if tmplStr := c.getObject(templateName + ".html"); tmplStr != "" {
		t := template.New(templateName)
		t, _ = t.Parse(tmplStr)
		c.cacheData.set(templateName, t)
		return t
	}
	return nil
}

// Close - should be called to clean up resources
func (c *Cache) Close() {
	c.storageClient.Close()
}

func (c *Cache) getObject(objName string) string {
	if rc, err := c.storageClient.Bucket(c.bucketName).Object(objName).NewReader(c.ctx); err == nil {
		if slurp, err := ioutil.ReadAll(rc); err == nil {
			defer rc.Close()
			return fmt.Sprintf("%s", slurp)
		}
		log.Fatalf("Cannot read object - %s", err)
	}
	return ""
}

func initStorageClient() *storage.Client {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatal("Cannot create storage client")
	} else {
		return client
	}
	return nil
}
