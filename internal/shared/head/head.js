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

