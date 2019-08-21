function logout() {
    axios.post('/logout')
        .then(function (response) {
            window.location.href = "/";
        }).catch(function (error) {
        window.location.href = "/";
    });
}

function showAlert(name,t) {
    //动画
    alt = $(name);
    if(t==0){
        alt.removeClass("alert-dander")
        alt.addClass("alert-success")
    }else {
        alt.removeClass("alert-success")
        alt.addClass("alert-danger")
    }
    if (alt.css("display") == 'none') {
        alt.slideDown();
        alt.slideUp(400);
    }
}