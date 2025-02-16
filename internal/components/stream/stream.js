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
