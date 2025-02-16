// auth is used for signing up and signing in/out. path could be:
//
//     /signup       /signin        /signout
//
//let password = document.getElementById("password");
//let username = document.getElementById("username");
//let form = document.getElementById("authForm");
//form.addEventListener("input", function () {
//        username.classList.remove("false");
//        password.classList.remove("false");
//        validateFormData()
//});
//function validateFormData() {
//        let pv = (password.value.length > 6);
//        let nv = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(username.value);
//        if (pv && nv) { return true }
//        password.classList.add(pv)
//        username.classList.add(nv)
//        return false
//}
//async function auth(path) {
//        if (validateFormData()) {
//                let response = await fetch("/"+path, {
//                        method: "POST",
//                        body: JSON.stringify({
//                                password: password.value,
//                                username: username.value,
//                        }),
//                });
//                let res = await response.json();
//                if (res.status == "success") {
//                        location.reload();
//                        return
//                } 
//                document.getElementById("errorDiv").innerHTML = res.status;
//        }
//}
