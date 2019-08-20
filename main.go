package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/newham/docer/api"
	"github.com/newham/hamgo"
)

func main() {
	server := hamgo.New(hamgo.Properties{})
	server.Static("public")
	server.Get("/", index)
	server.Handler("/article/=name", article, "GET,DELETE")
	server.Get("/edit/=name", edit)
	server.Get("/folder", folder)
	server.Post("/article", newArticle)
	server.RunAt("8089")
}

func index(ctx hamgo.Context) {
	ctx.HTML("view/index.html")
}

func article(ctx hamgo.Context) {
	name := ctx.PathParam("name")
	if ctx.Method() == http.MethodGet {
		b, err := ioutil.ReadFile("articles/" + name)
		if err != nil {
			ctx.JSONFrom(404, newMsg(404, ""))
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
