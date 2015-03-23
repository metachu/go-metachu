package main

type FileActionJson struct {
	Filename    string `form:"filename" binding:"required"`
	Filemd5     string `form:"filemd5"  binding:"required"`
	Action      string `form:"action"`
	Filepath    string `form:"filepath"`
	Newfilepath string `form:"newfilepath"`
	Newname     string `form:"newname"`
}
