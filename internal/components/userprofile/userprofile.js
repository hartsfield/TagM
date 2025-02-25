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
async function submitEdits(elid, cls) {
        const form = document.getElementById(elid);
        const data = new FormData(form);

        let response = await fetch("/edit", {
                method: "POST", 
                body: data
        });
        let res = await response.json();
        console.log(res);
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

        let response = await fetch("/edit", {
                method: "POST", 
                body: data
        });
        let res = await response.json();
        if (res.status == "success") {
                console.log("success");
        } else {
                document.getElementById("errorField").innerHTML = res.error;
        }
}
