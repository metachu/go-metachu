package main

type M map[string]interface{}

type File struct {
	Name             string
	MD5              string
	Size             int64
	Humansize        string
	Ahref            string
	Icon             string
	LastModified     string
	LastModifiedTime int64
	IsDir            bool
	RelativePath     string
}

type Configuration struct {
	RootPath string
}
