package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
)

type M map[string]interface{}

type afile struct {
	Name             string
	Size             int64
	Humansize        string
	Ahref            string
	Icon             string
	LastModified     string
	LastModifiedTime int64
	IsDir            bool
}

type Configuration struct {
	RootPath string
}

var configuration = Configuration{RootPath: "/home/vagrant"}

func main() {
	loadConfig()
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Use(martini.Static(configuration.RootPath, martini.StaticOptions{Prefix: "download/files/browse"}))

	m.Get("/", func() string {
		return "Hello World!"
	})

	m.Get("/files", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/files/browse/", 302)
	})

	m.Get("/files/browse/**", func(params martini.Params, r render.Render, req *http.Request) {
		raw_path, _ := filepath.Abs(filepath.Join(configuration.RootPath, params["_1"]))
		raw_path = filepath.Clean(raw_path)
		cleanurl := filepath.Clean(req.URL.String())
		file, err := os.Stat(raw_path)
		if err != nil {
			// file DOES NOT EXISTS
			fmt.Println("We have encountered an error!")
			context := M{"title": "401 - Path not found", "flash": "401 - Path not found. You encountered an error!"}
			r.HTML(401, "401", context)
			return
		}

		switch mode := file.Mode(); {

		case mode.IsDir():
			rawfiles, _ := ioutil.ReadDir(raw_path)
			afiles := make([]afile, len(rawfiles))

			for i, f := range rawfiles {
				afiles[i].Name = truncate(f.Name(), 20)
				afiles[i].Size = f.Size()
				afiles[i].Humansize = humanSize(f.Size())
				afiles[i].LastModified = f.ModTime().Format("2006-01-02 03:04:05 PM ")
				afiles[i].LastModifiedTime = f.ModTime().Unix()
				afiles[i].IsDir = f.IsDir()
				afiles[i].Ahref = filepath.Join(req.URL.String(), f.Name())
				if f.IsDir() {

					afiles[i].Icon = "mdi-file-folder-open"
				} else {
					//setting download url
					afiles[i].Ahref = filepath.Join("/download", afiles[i].Ahref)

					//setting default icon
					afiles[i].Icon = "mdi-file-file-download"
					//non default icons
					switch filepath.Ext(filepath.Join(raw_path, f.Name())) {
					case ".txt":
						afiles[i].Icon = "mdi-content-text-format"
					case ".pdf", ".mobi", ".epub":
						afiles[i].Icon = "mdi-av-my-library-books"
					case ".mp4", ".avi", ".ogg", ".wmv", ".flv":
						afiles[i].Icon = "mdi-maps-local-movies"
					}

				}
			}
			context := M{"path": cleanurl, "files": afiles}
			r.HTML(200, "index", context)

		case mode.IsRegular():
			context := M{"path": cleanurl}
			r.HTML(200, "index", context)
		}

	})

	m.RunOnAddr("0.0.0.0:3000")
}

func humanSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	var exp int
	var pre string
	pre = "KMGTPE"
	exp = int(math.Log(float64(bytes)) / math.Log(float64(1024)))

	return fmt.Sprintf("%.2f %sB", float64(bytes)/math.Pow(float64(1024), float64(exp)), string(pre[exp-1]))
}

func truncate(s string, i int) string {
	runes := []rune(s)
	if len(runes) > i {
		return string(runes[:i]) + "..."
	}
	return s
}

func loadConfig() {
	file, err := os.Open("conf.json")
	if err != nil {
		fmt.Println("Error! Config conf.json could not be found so defaulting root_path to /home/vagrant")
		return
	}
	decoder := json.NewDecoder(file)
	new_configuration := Configuration{}
	err = decoder.Decode(&new_configuration)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	configuration = new_configuration
	fmt.Println("Configuration successfully loaded! path is:", configuration.RootPath)
}
