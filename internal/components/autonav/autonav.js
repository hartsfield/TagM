// Provided Under BSD (2 Clause)
//
// Copyright 2025 Johnathan A. Hartsfield
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS “AS IS”
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.
//
// ////////////////////////////////////////////////////////////////////////////
//
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
        hotn.placeholder = "";
        chrn.placeholder = "";
        togl.innerHTML = "";
        togh.innerHTML = "";
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
