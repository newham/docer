package main

import (
	"encoding/json"
	"github.com/newham/hamgo"
	"io/ioutil"
	"os"
)

func main() {
	server := hamgo.New(hamgo.Properties{})
	server.Static("public")
	server.Get("/", index)
	server.Get("/article/=name", article)
	server.Get("/articles", articles)
	server.Post("/article", newArticle)
	server.RunAt("8089")
}

func index(ctx hamgo.Context) {
	ctx.HTML("view/index.html")
}

func article(ctx hamgo.Context) {
	name := ctx.PathParam("name")
	b, err := ioutil.ReadFile("articles/" + name)
	if err != nil {
		ctx.JSONFrom(404, newMsg(404, ""))
		return
	}
	ctx.JSON(200, b)
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

	b, _ := json.Marshal(a)
	err = ioutil.WriteFile("articles/"+a.File, b, os.ModePerm)
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

func articles(ctx hamgo.Context) {

}
