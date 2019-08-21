// var show_type = 1;
var vm = new Vue({
    delimiters: ['${', '}'],
    el: '#docer',
    data: {
        level: [],
        article: {},
    },
    methods: {
        getArticle: function (name) {
            axios.get('/article/' + name)
                .then(function (response) {
                    // console.log(response.data.title);
                    console.log(response.status);
                    // console.log(response.statusText);
                    // console.log(response.headers);
                    // console.log(response.config);
                    vm.article = response.data; // 这里不能用this，否则无法跟新

                    for (i = 0; i < vm.article.chapters.length; i++) {
                        if (vm.article.chapters[i].level == 1) {
                            vm.level.push([0, 0, 0, 0, 0]);
                            console.log(vm.level[0]);
                        }
                    }
                })
                .catch(function (error) {
                    console.log(error);
                    alert("打开文件错误");
                    window.location.href = "/";
                });
        },
        showPreview: function () {

        }
    },
    updated: function () {

    }
})