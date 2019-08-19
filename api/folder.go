package api

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

type Folder struct {
	Up      string   `json:"up"`
	Path    string   `json:"path"`
	Folders []string `json:"folders"`
	Files   []File   `json:"files"`
}
type File struct {
	Name     string `json:"name"`
	Size     string `json:"size"`
	Path     string `json:"path"`
	ModTime  string `json:"modTime"`
	Type     string `json:"type"`
	Editable bool   `json:"editable"`
}

const (
	ROOT_PATH = "articles"
)

var EDITABLE_TYPE = []string{"txt", "md", "markdown", "h", "c", "cpp", "c++", "go", "xml", "json", "java", "conf", "ini", "css", "js", "sh", "py", "log"}

func initRoot() {
	if !CheckFileIsExist(ROOT_PATH) {
		err := os.MkdirAll(ROOT_PATH, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func GetFolder(path string) Folder {
	dir, err := ioutil.ReadDir(ROOT_PATH + path)
	if err != nil {
		initRoot()
		return Folder{Path: "/"}
	}
	folders := make([]string, 0, 10)
	files := make([]File, 0, 10)

	for _, fi := range dir {
		if fi.IsDir() {
			folders = append(folders, fi.Name()+"/")
		} else {
			fileType := getType(fi.Name())
			files = append(files, File{fi.Name(), formatSize(fi.Size()), path, fi.ModTime().String()[:16], fileType, isEditable(fileType)})
		}

	}
	return Folder{getParentDirectory(path), path, folders, files}
}

func formatSize(size int64) string {
	const len = 1024
	var b, kb, mb, gb, tb, n int64
	var result string

	if size < len {
		b = size
		n = 1
	} else if size/len < len {
		kb = size / len
		n = 2
	} else if kb/len < len {
		mb = size / (len * len)
		n = 3
	} else if mb/len < len {
		gb = size / (len * len * len)
		n = 4
	} else {
		tb = size / (len * len * len * len)
		n = 5
	}

	switch n {
	case 1:
		result = strconv.FormatInt(b, 10) + "B"
		break
	case 2:
		result = strconv.FormatInt(kb, 10) + "KB"
		break
	case 3:
		result = strconv.FormatInt(mb, 10) + "MB"
		break
	case 4:
		result = strconv.FormatInt(gb, 10) + "GB"
		break
	case 5:
		result = strconv.FormatInt(tb, 10) + "TB"
		break
	}
	return result
}

func getType(fileName string) string {
	ext := path.Ext(fileName)
	if len(ext) < 2 {
		return "file"
	}
	return strings.ToLower(ext[1:])
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func getParentDirectory(dirctory string) string {
	p := path.Dir(path.Dir(dirctory))
	if !strings.HasSuffix(p, "/") {
		p = p + "/"
	}
	return p
}

func isEditable(fileType string) bool {
	for _, v := range EDITABLE_TYPE {
		if v == fileType {
			return true
		}
	}
	return false
}
