package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var Config = Configuration{RootPath: "/home/vagrant"}

func main() {
	loadConfig()
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Static(Config.RootPath, martini.StaticOptions{Prefix: "download/files/browse"}))
	m.Use(martini.Static("public", martini.StaticOptions{Prefix: "public"}))
	m.Get("/", func() string {
		return "Hello World!"
	})

	m.Get("/files", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/files/browse/", 302)
	})

	m.Get("/files/browse", FileBrowserHandler)
	m.Get("/files/browse/**", FileBrowserHandler)

	m.Post("/files/json/zip", binding.Form(FileActionJson{}), FileZipHandler)
	m.Post("/files/json/zip/", binding.Form(FileActionJson{}), FileZipHandler)

	m.Post("/files/json/delete", binding.Form(FileActionJson{}), FileDeleteHandler)
	m.Post("/files/json/delete/", binding.Form(FileActionJson{}), FileDeleteHandler)

	m.RunOnAddr("0.0.0.0:3000")

}

func FileDeleteHandler(r render.Render, fileAction FileActionJson, req *http.Request) {
	fmt.Println(fileAction)
	raw_path, rel_path, err := ValidatePath(fileAction.Filepath)
	fmt.Println("Delete: ", raw_path, rel_path)
	if err != nil {
		r.JSON(400, ERROR_INVALID_PATH)
		return
	}
	err = os.Remove(raw_path)
	if err != nil {
		r.JSON(400, ERROR_NO_OS_PERMISSION)
		return
	}
	r.JSON(200, M{"data": M{"title": "Success", "detail": "Deleted the file " + rel_path}})

}
func FileZipHandler(r render.Render, fileAction FileActionJson, req *http.Request) {
	raw_path, rel_path, err := ValidatePath(fileAction.Filepath)
	fmt.Println("Zip: ", raw_path, rel_path)
	if err != nil {
		r.JSON(400, ERROR_INVALID_PATH)
		return
	}
	fmt.Println(fileAction.Newname, raw_path)
	cmd := exec.Command("zip", "-r", fileAction.Newname, raw_path)
	cmd.Dir = filepath.Dir(raw_path)
	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
		r.JSON(400, ERROR_ZIP_COMMAND)
		return
	}
	if fileAction.Action == "zip" {
		r.JSON(200, M{"data": M{"title": "Zip Started", "detail": "The zip archive was created. Refreshing page in a moment"}})
		return
	}
	if fileAction.Action == "zipAndDelete" {
		err = os.Remove(raw_path)
		if err != nil {
			r.JSON(400, ERROR_NO_OS_PERMISSION)
			return
		}
		r.JSON(200, M{"data": M{"title": "Success.", "detail": "Zip archive was created and old file was successfully deleted. Refreshing."}})
		return
	}
}

func FileBrowserHandler(params martini.Params, r render.Render, req *http.Request) {
	raw_path, rel_path, err := ValidatePath(params["_1"])
	raw_path = filepath.Clean(raw_path)

	rawfiles, _ := ioutil.ReadDir(raw_path)
	fmt.Printf("rel_path %s\n", rel_path+"\t"+raw_path)
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

		files := make([]File, len(rawfiles))

		for i, f := range rawfiles {
			files[i].Name = truncate(f.Name(), 20)
			hash := md5.Sum([]byte(f.Name()))
			files[i].MD5 = hex.EncodeToString(hash[:])
			files[i].Size = f.Size()
			files[i].Humansize = humanSize(f.Size())
			files[i].LastModified = f.ModTime().Format("2006-01-02 03:04:05 PM ")
			files[i].LastModifiedTime = f.ModTime().Unix()
			files[i].IsDir = f.IsDir()
			files[i].Ahref = filepath.Join(req.URL.String(), f.Name())
			files[i].RelativePath = filepath.Join(rel_path, f.Name())

			if f.IsDir() {
				files[i].Icon = "mdi-file-folder-open"
			} else {
				//setting download url
				files[i].Ahref = filepath.Join("/download", files[i].Ahref)

				//setting default icon
				files[i].Icon = "mdi-file-file-download"
				//non default icons
				switch filepath.Ext(filepath.Join(raw_path, f.Name())) {
				case ".txt":
					files[i].Icon = "mdi-content-text-format"
				case ".pdf", ".mobi", ".epub":
					files[i].Icon = "mdi-av-my-library-books"
				case ".mp4", ".avi", ".ogg", ".wmv", ".flv":
					files[i].Icon = "mdi-maps-local-movies"
				}

			}
		}
		context := M{"path": raw_path, "files": files}
		r.HTML(200, "index", context)

	case mode.IsRegular():
		context := M{"path": raw_path}
		r.HTML(200, "index", context)
	}

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
	new_Config := Configuration{}
	err = decoder.Decode(&new_Config)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	Config = new_Config

	Config.RootPath, _ = filepath.Abs(Config.RootPath) // stripping trailing '/'

	fmt.Println("Config successfully loaded! path is:", Config.RootPath)
}
