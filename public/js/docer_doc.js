var current_chapter_index = 0, current_chapter_line_index = 0, current_reference_index = 0;
var add_type = -1;
// var show_type = 1;
var vm = new Vue({
    delimiters: ['${', '}'],
    el: '#docer',
    data: {
        newFilename: "新建文档.json",
        isEdited: "",
        content: "",
        article: {},
        folder: {}
    },
    methods: {
        newArticle: function (name) {
            if (name == "") {
                alert("file name is null");
                return
            }
            this.article = {
                file: name,
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
                references: [
                    {
                        id: (new Date()).valueOf(),
                        text: ""
                    }
                ]
            }
        },
        focusTitle: function (i) {
            document.getElementById("title").focus();
        },
        downTitle: function (i) {
            this.focusChapter(0);
        },
        addChapter: function (i, step) {
            // if(this.article.chapters[i].title==""){ //本行为空，禁止增加新行
            //     return
            // }
            i += step;
            console.log("add", i);
            this.article.chapters.splice(i, 0, {
                title: "",
                level: 1,
                lines: [
                    {content: "", translation: "", reference: -1},
                ]
            });
            current_chapter_index = i;
            add_type = 1;
        },
        delChapter: function (i) {
            if (this.article.chapters.length == 1) {
                return;
            }
            this.article.chapters.splice(i, 1);
            console.log("del", i);
            //foucs up item
            current_chapter_index = i - 1;
            add_type = 1;
        },
        focusChapter: function (i) {
            if (i < this.article.chapters.length && i >= 0) {
                document.getElementById("c-" + i).focus();
            }
        },
        upChapter: function (i) {
            if (i >= 1) {
                this.focusChapter(i - 1);
            }

            if (i == 0) {
                this.focusTitle();
            }
        },
        downChapter: function (i) {
            if (i <= this.article.chapters.length - 2) {
                this.focusChapter(i + 1);
            }
        },
        rightChapter: function (i) {
            // console.log("right", i);
            this.focusContent(i, 0);
        },
        addContent: function (i, j, step) {
            j += step;
            console.log("add", i, j);
            this.article.chapters[i].lines.splice(j, 0, {
                content: "",
                translation: "",
                reference: this.article.chapters[i].lines[j - step].reference
            })
            current_chapter_index = i;
            current_chapter_line_index = j;
            add_type = 2;
        },
        delContent: function (i, j) {
            // if (confirm("sure to delete:"+this.article.chapters[i].lines[i].content+"?")==false){
            //     return
            // }
            console.log("del", i, j);
            if (this.article.chapters[i].lines.length == 1) {
                return;
            }
            this.article.chapters[i].lines.splice(j, 1);
            //foucs up item
            current_chapter_index = i;
            current_chapter_line_index = j - 1;
            add_type = 2;
        },
        upContent: function (i, j) {
            // console.log("up", i, j);
            if (j >= 1) {
                this.focusContent(i, j - 1);
            }
            if (j == 0 && i - 1 >= 0) {
                this.focusContent(i - 1, this.article.chapters[i - 1].lines.length - 1)
            }

        },
        downContent: function (i, j) {
            // console.log("down", i, j);
            if (j <= this.article.chapters[i].lines.length - 2) {
                this.focusContent(i, j + 1);
            }
            if (j == this.article.chapters[i].lines.length - 1 && i + 1 <= this.article.chapters.length - 1) {
                this.focusContent(i + 1, 0);
            }

        },
        leftContent: function (i, j) {
            // console.log("left", i, j);
            this.focusChapter(i);
        },
        focusContent: function (i, j) {
            if (i >= 0 && j >= 0 && j < this.article.chapters[i].lines.length) {
                document.getElementById("c-" + i + "-c-" + j).focus();

                current_chapter_index = i;
                current_chapter_line_index = j;
                this.content = this.article.chapters[i].lines[j].content;

            }
        },
        onfocusContent: function (i, j) {
            console.log("onfocusContent", i, j)
        },
        addReference: function (i) {
            i++;
            console.log("add reference", i);
            this.article.references.splice(i, 0, {id: (new Date()).valueOf(), text: ""})
            //focus
            current_reference_index = i;
            add_type = 3;
        },
        delReference: function (i) {
            if (this.article.references.length == 1) {
                return
            }
            console.log("del reference", i);
            this.article.references.splice(i, 1)
        },
        focusReference: function (i) {
            if (i < this.article.references.length && i >= 0) {
                document.getElementById("r-" + i).focus();
            }
        },
        setFooter: function () {
            if (this.article == null || this.article.chapters == null) {
                return
            }
            var words = 0;
            var lines = 0;
            for (i = 0; i < this.article.chapters.length; i++) {
                words += this.article.chapters[i].title.length;
                for (j = 0; j < this.article.chapters[i].lines.length; j++) {
                    words += this.article.chapters[i].lines[j].content.length;
                    if (this.article.chapters[i].lines[j].content != "") {
                        lines++;
                    }
                }
            }
            // console.log("words", total);
            if (this.article.words != words) {
                this.article.words = words;
            }
            if (this.article.lines != lines) {
                this.article.lines = lines;
            }
        },
        postArticle: function () {
            console.log("save", this.article.file);
            // return
            axios.post('/article', this.article, {
                headers:
                    {
                        'Content-Type': 'application/json'
                    }
            })
                .then(function (response) {
                    console.log(response);
                    alert("保存成功！");
                })
                .catch(function (error) {
                    console.log(error);
                    alert("保存失败,文件已经存在！");
                });
        },
        getArticle: function (name) {
            axios.get('/article/' + name)
                .then(function (response) {
                    // console.log(response.data.title);
                    console.log(response.status);
                    // console.log(response.statusText);
                    // console.log(response.headers);
                    // console.log(response.config);
                    vm.article = response.data; // 这里不能用this，否则无法跟新
                })
                .catch(function (error) {
                    console.log(error);
                    alert("打开文件错误");
                    window.location.href = "/";
                });
        },
        showArticles: function () {
            axios.get('/folder')
                .then(function (response) {
                    // console.log(response.data.title);
                    // console.log(response.status);
                    // console.log(response.statusText);
                    // console.log(response.headers);
                    // console.log(response.config);
                    vm.folder = response.data; // 这里不能用this，否则无法跟新
                });
            $('#folderModal').modal('show');
        },
        confirmArticle: function () {
            $('#confirmModal').modal('show');
        },
        openArticles: function (name) {
            window.open("/edit/" + name);
        }

    },
    updated: function () {
        // console.log("updated", add_type, current_chapter_index);
        if (add_type == 1) {
            this.focusChapter(current_chapter_index);
        } else if (add_type == 2) {
            this.focusContent(current_chapter_index, current_chapter_line_index);
        } else if (add_type == 3) {
            this.focusReference(current_reference_index);
        }
        add_type = -1;

        this.setFooter();

        // newContent = this.article.chapters[current_chapter_index].lines[current_chapter_line_index].content;
        // if (this.content != newContent) {
        //     this.content = newContent;
        // }
    }
})
