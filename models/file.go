package models

import (
	"io/ioutil"
	"os"
	"time"
)

type FileStat struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Mode    string    `json:"mode"`
	Dir     bool      `json:"dir"`
	ModTime time.Time `json:"modtime"`
}

func GetFileStat(filename string) (fs FileStat, err error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return
	}
	fs = getFileStatFromInfo(fi)
	return
}

func GetDirList(filename string) (fsl []FileStat, err error) {
	fileinfos, err := ioutil.ReadDir(filename)
	if err != nil {
		return
	}
	fsl = make([]FileStat, 0)
	for _, fileinfoi := range fileinfos {
		fsl = append(fsl, getFileStatFromInfo(fileinfoi))
	}
	return
}

func getFileStatFromInfo(fi os.FileInfo) (fs FileStat) {
	fs.Name = fi.Name()
	fs.Dir = fi.IsDir()
	fs.Mode = fi.Mode().String()
	fs.ModTime = fi.ModTime()
	return fs
}
