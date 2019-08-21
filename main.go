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

func init() {
	api.ROOT_PATH = "articles"
}

func main() {
	server := hamgo.New(hamgo.Properties{SessionMaxLifeTime: 3600 * 24})
	server.AddFilter(func(ctx hamgo.Context) bool {
		ctx.GetSession()
		if ctx.GetSession().Get(USER_NAME) == nil {
			if ctx.Method() == http.MethodGet {
				Html(ctx, "view/login.html")
			} else {
				ctx.JSONString(403, "403")
			}
			return false
		}
		return true
	}).AddAnnoURL("/login", "POST").AddAnnoURL("/about", "GET").AddAnnoURL("/help", "GET").AddAnnoURL("/register", "POST")

	server.Static("public")
	server.Get("/", index)
	server.Handler("/article/=name", article, "GET,DELETE")
	server.Get("/download/=name", download)
	server.Get("/edit/=name", edit)
	server.Get("/preview/=name", preview)
	server.Get("/folder", folder)
	server.Post("/article", newArticle)
	server.Post("/upload", upload)
	server.Post("/login", login)
	server.Post("/logout", logout)
	server.Get("/about", about)
	server.Get("/help", help)
	server.Post("/register", register)
	server.RunAt("8089")
}

func preview(ctx hamgo.Context) {
	name := ctx.PathParam("name")
	if !api.CheckFileIsExist(getArticle(ctx, name)) {
		ctx.WriteString("404")
		ctx.Text(404)
		return
	}
	ctx.PutData("filename", name)
	Html(ctx, "view/preview.html")
}

func register(ctx hamgo.Context) {
	username := ctx.FormValue("username")
	pwd := ctx.FormValue("pwd")
	confirmPwd := ctx.FormValue("confirmPwd")
	token := ctx.FormValue("token")

	if username != "" && pwd != "" && pwd == confirmPwd && token == getToken() {
		f, err := os.OpenFile("USER", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			hamgo.Log.Error(err.Error())
			return
		}
		_, err = f.WriteString("\n" + getBase64(username, pwd))
		if err != nil {
			hamgo.Log.Error(err.Error())
		}
	}
	ctx.Redirect("/")
}

func getToken() string {
	b, err := ioutil.ReadFile("TOKEN")
	if err != nil {
		return ""
	}
	return string(b)
}

const (
	USER_NAME = "username"
)

func index(ctx hamgo.Context) {
	Html(ctx, "view/index.html")
}

func about(ctx hamgo.Context) {
	Html(ctx, "view/about.html")
}

func help(ctx hamgo.Context) {
	Html(ctx, "view/help.html")
}

func download(ctx hamgo.Context) {
	name := ctx.PathParam("name")

	filepath := getArticle(ctx, name)
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
		b, err := ioutil.ReadFile(getArticle(ctx, name))
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
		err := os.Remove(getArticle(ctx, name))
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

	if create != 1 && !api.CheckFileIsExist(getArticle(ctx, name)) {
		ctx.WriteString("404")
		ctx.Text(404)
		return
	}
	ctx.PutData("filename", name)
	ctx.PutData("create", create)
	Html(ctx, "view/doc.html")
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

	filepath := getArticle(ctx, a.File)

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
	ctx.JSONFrom(200, api.GetFolder(getHome(ctx)))
}

func upload(ctx hamgo.Context) {
	file, fileHeader, err := ctx.FormFile("file")
	defer file.Close()
	if err != nil {
		ctx.JSONString(500, err.Error())
		return
	}
	//2.create local file
	f, err := os.OpenFile(getArticle(ctx, fileHeader.Filename), os.O_WRONLY|os.O_CREATE, 0666)
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
	if err == nil && b != nil && strings.Contains(string(b), getBase64(username, pwd)) {
		if err = api.MkHome(username); err != nil {
			hamgo.Log.Error(err.Error())
			ctx.Redirect("/")
			return
		}
		ctx.GetSession().Set(USER_NAME, username)
	}
	ctx.Redirect("/")
}

func logout(ctx hamgo.Context) {
	//ctx.GetSession().Delete(USER_NAME)
	ctx.DeleteSession()
	ctx.JSONString(200, "")
}

func Html(ctx hamgo.Context, html string) {
	ctx.HTML(html, "view/head.html", "view/title.html", "view/tool.html")
}

func getUsername(ctx hamgo.Context) string {
	return ctx.GetSession().Get(USER_NAME).(string)
}

func getHome(ctx hamgo.Context) string {
	return api.GetHome(getUsername(ctx))
}

func getArticle(ctx hamgo.Context, filename string) string {
	return getHome(ctx) + filename
}

func getBase64(username, pwd string) string {
	return api.Base64Encode(username + "," + pwd)
}
