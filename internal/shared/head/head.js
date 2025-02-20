function jumpTo(eid) {
        var jump = document.getElementById(eid);
        jump.scrollIntoView({
                behavior: 'auto',
                block: 'center',
                inline: 'center'
        });
}
const myElement = document.getElementById('stream');

myElement.addEventListener('scroll', function() {
        // Get the current vertical scroll position
        const scrollTop = myElement.scrollTop;

        // Log the scroll position (for example)
        console.log('Scrolled to:', scrollTop);

        // Example: Change element style based on scroll position
        if (scrollTop > 100) {
                myElement.classList.add('scrolled');
                (function() {
                        window.scroll({
                                top: 50,
                                left: 0,
                                behavior: 'auto'
                        });
                        document.getElementById("stream").classList.add("stream-loaded");
                })();

        } else {
                myElement.classList.remove('scrolled');
        }
});
function toggleDisplay(elem) {
        let formDisplay = document.getElementById("item-controls_"+elem);
        let butt = document.getElementById("item-shr-"+elem);
        if (formDisplay.style.display == "none" || formDisplay.style.display == "") {
                formDisplay.style.display = "flex";
                butt.innerHTML = "<"
        } else {
                formDisplay.style.display = "none";
                butt.innerHTML = "+"
        }
}
async function getExample(view) {
        const response = await fetch("/api/getExample", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({special:view}),
        });
        let res = await response.json();
        if (res.success == "true") {
                // do stuff
        } else {
                console.log("error");
        }
}
let toggled = false;
{{ if .Credentials.IsLoggedIn }}
window.onscroll = function(e) {
        // print "false" if direction is down and "true" if up
        if (this.oldScroll > this.scrollY) { if (!toggled) {toggleNew()} toggled = true }
        if (this.oldScroll < this.scrollY) {if (toggled) {toggleNew()} toggled = false}
        console.log(this.oldScroll > this.scrollY);
        this.oldScroll = this.scrollY;
}
{{ end }}
//setInterval(autoReload, 500); // 5000 milliseconds = 5 seconds
//async function autoReload() {
//        const response = await fetch("/wasmodified", {
//                method: "GET",
//                //headers: { "Content-Type": "application/json" },
//                //body: JSON.stringify({"na":"na"}),
//        });
//        let res = await response.json();
//        if (res.modified == "true") {
//                location.reload();
//        } 
//}
//
//document.getElementsByTagName("img").onerror='this.style.display = "none"' 
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
        console.log(res);
        if (res.success == "true") {
                document.getElementById("like_"+postID).innerHTML = res.score
        } else {
                //document.getElementById("errorField").innerHTML = res.error;
        }


}
async function share(postID) {
        let response = await fetch("/share", {
                method: "POST",
                body: {"id": postID},
        });

        let res = await response.json();
        handleResponse(res);
}
function handleResponse(res) {
        if (res.success == "true") {
                //window.location = window.location.origin;
        } else {
                //document.getElementById("errorField").innerHTML = res.error;
        }
}

//function toggleReply(id) {
//        var elm = document.getElementById(id);
//        var butt = document.getElementById("item-comment-submit-"+id);
//        if (elm.classList.contains("item-comment-submit")) 
//                elm.classList.remove("item-comment-submit");
//        butt.innerHTML = "x"
//        return
//}
//elm.classList.add("item-comment-submit");
//butt.innerHTML = "+"
//}

