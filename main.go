package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/newham/docer/api"
	"github.com/newham/hamgo"
)

func main() {
	server := hamgo.New(hamgo.Properties{SessionMaxLifeTime: 3600 * 24})
	server.AddFilter(func(ctx hamgo.Context) bool {
		ctx.GetSession()
		if ctx.GetSession().Get(SESSION_NAME) == nil {
			if ctx.Method() == http.MethodGet {
				ctx.HTML("view/login.html")
			} else {
				ctx.JSONString(403, "403")
			}
			return false
		}
		return true
	}).AddAnnoURL("/login", "POST")
	server.Static("public")
	server.Get("/", index)
	server.Handler("/article/=name", article, "GET,DELETE")
	server.Get("/download/=name", download)
	server.Get("/edit/=name", edit)
	server.Get("/folder", folder)
	server.Post("/article", newArticle)
	server.Post("/upload", upload)
	server.Post("/login", login)
	server.Post("/logout", logout)
	server.RunAt("8089")
}

const (
	SESSION_NAME = "username"
)

func index(ctx hamgo.Context) {
	ctx.HTML("view/index.html")
}

func download(ctx hamgo.Context) {
	name := ctx.PathParam("name")

	filepath := "articles/" + name
	if !api.CheckFileIsExist(filepath) {
		ctx.JSONString(404, "file not existed")
		return
	}

	ctx.Attachment(filepath)
	return
}

func article(ctx hamgo.Context) {
	name := ctx.PathParam("name")
	if ctx.Method() == http.MethodGet {
		b, err := ioutil.ReadFile("articles/" + name)
		if err != nil {
			ctx.JSONFrom(404, newMsg(404, ""))
			return
		}
		a := Article{}
		if json.Unmarshal(b, &a) != nil {
			ctx.JSONFrom(400, newMsg(400, "打开文件错误，请检查文件格式"))
			return
		}
		ctx.JSON(200, b)
	} else if ctx.Method() == http.MethodDelete {
		err := os.Remove("articles/" + name)
		if err != nil {
			ctx.JSONFrom(500, newMsg(500, err.Error()))
			return
		}
		ctx.JSON(200, nil)
	}

}

func edit(ctx hamgo.Context) {
	create := 0
	if ctx.FormValue("create") != "" {
		create = 1
	}
	name := ctx.PathParam("name")

	if create != 1 && !api.CheckFileIsExist(api.ROOT_PATH+"/"+name) {
		ctx.WriteString("404")
		ctx.Text(404)
		return
	}
	ctx.PutData("filename", name)
	ctx.PutData("create", create)
	ctx.HTML("view/doc.html")
}

/**
this.article = {
                    file: "new.json",
                    title: "",
                    words: 0,
                    lines: 0,
                    chapters: [
                        {
                            level: 1,
                            title: "",
                            lines: [
                                {
                                    content: "",
                                    translation: "",
                                    reference: -1
                                },
                            ]
                        }
                    ],
                    references: [""]
                }
*/

type Article struct {
	File       string      `json:"file"`
	Title      string      `json:"title"`
	Words      int         `json:"words"`
	Lines      int         `json:"lines"`
	Chapters   []Chapters  `json:"chapters"`
	References []Reference `json:"references"`
}

type Chapters struct {
	Level int    `json:"level"`
	Title string `json:"title"`
	Lines []Line `json:"lines"`
}

type Line struct {
	Content     string `json:"content"`
	Translation string `json:"translation"`
	Reference   int    `json:"reference"`
}

type Reference struct {
	Id   int64  `json:"id"`
	Text string `json:"text"`
}

func newArticle(ctx hamgo.Context) {

	a := Article{}
	err := ctx.BindJSON(&a)
	if err != nil {
		println(err.Error())
		ctx.JSONString(400, err.Error())
		return
	}

	filepath := api.ROOT_PATH + "/" + a.File

	// if api.CheckFileIsExist(filepath) {
	// 	ctx.JSONString(400, "file existed")
	// 	return
	// }

	b, _ := json.Marshal(a)
	err = ioutil.WriteFile(filepath, b, os.ModePerm)
	if err != nil {
		ctx.WriteString(err.Error())
		ctx.Text(500)
		return
	}
	ctx.JSONString(200, "")
}

func newMsg(code int, msg string) map[string]interface{} {
	return map[string]interface{}{"code": code, "msg": msg}
}

func folder(ctx hamgo.Context) {
	ctx.JSONFrom(200, api.GetFolder("/"))
}

func upload(ctx hamgo.Context) {
	file, fileHeader, err := ctx.FormFile("file")
	defer file.Close()
	if err != nil {
		ctx.JSONString(500, err.Error())
		return
	}
	//2.create local file
	f, err := os.OpenFile(api.ROOT_PATH+"/"+fileHeader.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		ctx.JSONString(500, err.Error())
		return
	}
	//3.copy uploadfile to localfile
	io.Copy(f, file)
	ctx.Redirect("/")
}

func login(ctx hamgo.Context) {
	username := ctx.FormValue("username")
	pwd := ctx.FormValue("pwd")
	b, err := ioutil.ReadFile("USER")
	//println(string(b),api.Base64Encode(username+","+pwd))
	if err == nil && b != nil && strings.Contains(string(b), api.Base64Encode(username+","+pwd)) {
		ctx.GetSession().Set(SESSION_NAME, username)
	}
	ctx.Redirect("/")
}

func logout(ctx hamgo.Context) {
	//ctx.GetSession().Delete(SESSION_NAME)
	ctx.DeleteSession()
	ctx.JSONString(200, "")
}
