package server

import (
	"github.com/gofiber/template/django/v3"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

func initViews() *django.Engine {
	engine := django.New("./views", ".html")
	engine.Reload(true)
	engine.AddFunc("css", func(name string) (res template.HTML) {
		err := filepath.Walk("public/assets", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == name {
				res = template.HTML("<link rel=\"stylesheet\" href=\"/" + path + "\">")
			}
			return nil
		})
		if err != nil {
			log.Fatalf("failed to create django engine, unable to walk puiblic/assets folder. Error: %v", err)
		}
		return
	})
	return engine
}
