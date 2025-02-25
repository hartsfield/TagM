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
function toggleDisplay(elem) {
        let formDisplay = document.getElementById("item-controls_"+elem);
        let butt = document.getElementById("item-shr-"+elem);

        if   (  formDisplay.style.display == "none" ||
                formDisplay.style.display == "")        
             {  formDisplay.style.display = "flex"; butt.innerHTML = "<";} 
        else {  formDisplay.style.display = "none"; butt.innerHTML = "+";}
}
let did_submit_reply = false
async function submitReply(parent) {
        if (!did_submit_reply) {
                let txt = document.getElementById("uptext_"+parent).value
                let response = await fetch("/reply", {
                        method: "POST",
                        body: JSON.stringify({"parent": parent, "uptext": txt}),
                });
                let res = await response.json();
                if (res.status == "success") {
                        window.location = window.location.origin + "/view/"+res.ID;
                }
        }
}
async function like(postID) {
        let response = await fetch("/like/"+postID, {
                method: "POST",
                body: {"id": postID},
        });

        let res = await response.json();
        if (res.success == "true") {
                document.getElementById("like_"+postID).innerHTML = res.score
        } else {
                document.getElementById("errorField").innerHTML = res.error;
        }
}
async function share(postID) {
        let response = await fetch("/share", {method: "POST", body: {"id": postID}});
        let res      = await response.json();

        if (res.success == "true") {window.location = window.location.origin;} 
        else {document.getElementById("errorField").innerHTML = res.error;}
}
//let toggled = false;
//{{ if .Credentials.IsLoggedIn }}
//window.onscroll = function(e) {
//        // print "false" if direction is down and "true" if up
//        if (this.oldScroll > this.scrollY) { if (!toggled) {toggleNew()} toggled = true }
//        if (this.oldScroll < this.scrollY) {if (toggled) {toggleNew()} toggled = false}
//        console.log(this.oldScroll > this.scrollY);
//        this.oldScroll = this.scrollY;
//}
//{{ end }}

