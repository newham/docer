function logout() {
    axios.post('/logout')
        .then(function (response) {
            window.location.href = "/";
        }).catch(function (error) {
            window.location.href = "/";
        });
}