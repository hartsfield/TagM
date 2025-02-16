async function submitEdits(elid, cls) {
        const form = document.getElementById(elid);
        const data = new FormData(form);

        let response = await fetch("/edit", {method: "POST", body: data});
        let res = await response.json();
        if (res.status == "success") {
                document.getElementById(cls).style.backgroundImage = "url(/" + res.payload + ")";
        } else {
                document.getElementById("errorField").innerHTML = res.error;
        }
}
async function submitEditsInputs(elid) {
        const form = document.getElementById("inputs-form");
        const data = new FormData(form);
        console.log(data, form)

        let response = await fetch("/edit", {method: "POST", body: data});
        let res = await response.json();
        if (res.status == "success") {
                console.log("success");
        } else {
                document.getElementById("errorField").innerHTML = res.error;
        }
}
