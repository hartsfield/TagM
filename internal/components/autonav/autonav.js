//let np = document.getElementById("nav-portrait");
//np.style.position = "absolute";
//np.style.right = "-" + np.offsetWidth + "px";
//function showNavPortrait() {
//        np.style.right = 0;
//        setTimeout(function () {
//                document.addEventListener('click', tf, false);
//        }, 50);
//}

var authToggled = false;

var newToggled = false;
function toggleNew() {
    let upForm = document.getElementById("form-wrapper");
    let tagdiv = document.getElementById("up-sym-wrap");
    let stream = document.getElementById("stream");
    if (!newToggled) {
        upForm.classList.add("upload-outer-toggle");
        tagdiv.classList.add("upload-outer-toggle");
        stream.classList.add("upload-outer-toggle-stream");
        newToggled= true;
    } else {
        upForm.classList.remove("upload-outer-toggle");
        tagdiv.classList.remove("upload-outer-toggle");
        stream.classList.remove("upload-outer-toggle-stream");
        newToggled= false;
    }
}
// auth is used for signing up and signing in/out. path could be:
//
//     /signup       /signin        /signout
//
{{ if eq .Credentials.User.ID .Profile.ID }}
let password = document.getElementById("chron-nav");
let username = document.getElementById("hot-nav");
let form = document.getElementById("authForm");
form.addEventListener("input", function () {
    username.classList.remove("false");
    password.classList.remove("false");
    validateFormData()
});
{{ end }}
function validateFormData() {
    let pv = (password.value.length > 6);
    let nv = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(username.value);
    if (pv && nv) { return true }
    password.classList.add(pv)
    username.classList.add(nv)
    return false
}
async function auth(path) {
    if (validateFormData()) {
        let response = await fetch("/"+path, {
            method: "POST",
            body: JSON.stringify({
                password: password.value,
                username: username.value,
            }),
        });
        let res = await response.json();
        if (res.status == "success") {
            location.reload();
            return
        } 
        document.getElementById("errorDiv").innerHTML = res.status;
    }
}
function tf() {
    np.style.right = "-" + np.offsetWidth + "px";
    document.removeEventListener('click', tf);
}
function toggleAuth() {
    let auth = document.getElementById("logo-nav");
    let togl = document.getElementById("nav-toggle-all");
    let togh = document.getElementById("nav-toggle-all-hid");
    let tog2 = document.getElementById("nav-toggle-all-hid2");
    let hotn = document.getElementById("hot-nav");
    let chrn = document.getElementById("chron-nav");
    if (!authToggled) {
        hotn.placeholder = "email";
        chrn.placeholder = "password";
        togl.innerHTML = "signup";
        togh.innerHTML = "signin";
        username.disabled = false;
        password.disabled = false;
        auth.classList.add("logo-shrink-toggle");
        chrn.classList.add("chron-nav-auth");
        hotn.classList.add("chron-nav-auth");
        togl.classList.add("nav-toggled");
        togh.classList.add("nta2-toggled");
        tog2.classList.add("nta3-toggled");
        togl.onclick = async function auth(path) {
            if (validateFormData()) {
                let response = await fetch("/signup", {
                    method: "POST",
                    body: JSON.stringify({
                        password: password.value,
                        username: username.value,
                    }),
                });
                let res = await response.json();
                if (res.status == "success") {
                    location.reload();
                    return
                } 
                document.getElementById("errorDiv").innerHTML = res.status;
            }
        }

        togh.onclick =  async function auth(path) {
            if (validateFormData()) {
                let response = await fetch("/signin", {
                    method: "POST",
                    body: JSON.stringify({
                        password: password.value,
                        username: username.value,
                    }),
                });
                let res = await response.json();
                if (res.status == "success") {
                    location.reload();
                    return
                } 
                document.getElementById("errorDiv").innerHTML = res.status;
            }
        }

        //tog2.onclick = toggleAuth();
        chrn.onclick = "";
        hotn.onclick = "";
        authToggled = true;
    } else {
        hotn.placeholder = ""
        chrn.placeholder = ""
        togl.innerHTML = ""
        togh.innerHTML = ""
        username.disabled = true;
        password.disabled = true;

        auth.classList.remove("logo-shrink-toggle");
        chrn.classList.remove("chron-nav-auth");
        hotn.classList.remove("chron-nav-auth");
        togl.classList.remove("nav-toggled");
        togh.classList.remove("nta2-toggled");
        tog2.classList.remove("nta3-toggled");

        authToggled = false;
    }
}
