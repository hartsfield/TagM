function jumpTo(eid) {
        var jump = document.getElementById(eid);
        jump.scrollIntoView({
                behavior: 'auto',
                block: 'center',
                inline: 'center'
        });
}
function toggleDisplay(elem) {
        let divs = document.getElementById("hiddenTop").children;
        let formDisplay = document.getElementById(elem);
        for (let i=0;i<divs.length;i++) {
                if (divs[i].id != formDisplay.id) {
                        divs[i].style.display = "none";
                }
        }
        if (formDisplay.style.display == "none" || formDisplay.style.display == "") {
                formDisplay.style.display = "unset";
        } else {
                formDisplay.style.display = "none";
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
window.onscroll = function(e) {
        // print "false" if direction is down and "true" if up
        if (this.oldScroll > this.scrollY) { if (!toggled) {toggleNew()} toggled = true }
        if (this.oldScroll < this.scrollY) {if (toggled) {toggleNew()} toggled = false}
        console.log(this.oldScroll > this.scrollY);
        this.oldScroll = this.scrollY;
}
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
