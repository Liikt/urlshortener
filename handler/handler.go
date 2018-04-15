package handler

import (
	"net/http"

	"github.com/boltdb/bolt"

	yaml "gopkg.in/yaml.v2"
)

type urlMap struct {
	paths map[string]string
}

type mapItem struct {
	Path string
	Url  string
}

var globalMap *urlMap = &urlMap{make(map[string]string)}

func (*urlMap) redirect(w http.ResponseWriter, r *http.Request) {
	if url, ok := globalMap.paths[r.URL.String()]; ok {
		http.Redirect(w, r, url, 307)
		return
	}
}

func YAMLHandler(yml []byte, fallback *http.ServeMux) error {
	var list []mapItem
	err := yaml.Unmarshal(yml, &list)
	if err != nil {
		return err
	}

	for _, item := range list {
		globalMap.paths[item.Path] = item.Url
		fallback.HandleFunc(item.Path, globalMap.redirect)
	}

	return nil
}

func BoltHandler(path string, fallback *http.ServeMux) error {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil
	}
	defer db.Close()

	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("paths"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			globalMap.paths[string(k)] = string(v)
			fallback.HandleFunc(string(k), globalMap.redirect)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
