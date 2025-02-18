let did_submit = false
async function submitPost() {
        if (!did_submit) {
                did_submit = true;
                const form = document.getElementById("uploadForm");
                const data = new FormData(form);
                let response = await fetch("/uploadItem", {
                        method: "POST",
                        body: data,
                });
                let res = await response.json();
                handleResponse(res);
        }
}
function handleResponse(res) {
        console.log(res);
        if (res.status == "success") {
                location.reload()
        } else {
                document.getElementById("errorField").innerHTML = res.error;
        }
}
async function submitImg() {
        const form = document.getElementById("img-uploadForm");
        const data = new FormData(form);
        let response = await fetch("/uploadItem", {
                method: "POST",
                body: data,
        });
        let res = await response.json();
        console.log(res);
}

// !bang @mention #hash_tag $cash %evil ^good &pointer *star 
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
//                _ _________               _____________________
// ______________ # TAGRippeR _____ - _____ A sigma-firma product _____________
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// This is something like a "keyword highlighter" for html/text fields, but 
// its designed originally and foremost to extract #hash_tags from user input 
// as they type. It is pure JavaScript, no third party dependencies. 
//
// We use a special type of html text input/field that's just a regular element 
// with the attribute "contenteditable" set like so:
//
// <div contenteditable id="MyText">dream eternal</div>
//
// ----------------------------------------------------------------------------
// -----------------------        SET UP      ---------------------------------
// ----------------------------------------------------------------------------
// Place the 'id' of the html element you've made 'contenteditable' here, in
// place of MyText:
let editable = "uptext"
//
// The text is used as place holder text. When a user clicks/taps the above 
// <div>, its deleted. This is done only once, and is the purpose of the 
// following:
let unsealed = false
document.getElementById(editable).addEventListener("click", 
        function () { if(!unsealed) {
                let element         = document.getElementById(editable);
                element.innerText   = "";
                element.style.color = "black";
                unsealed            = true;}});
// Here, we listen for 'input' on the element with id="uploadForm", running the 
// tagification function (tagify()), on every input event:
document.getElementById("uploadForm").addEventListener("input", function () {
        tagify();
});
// ----------------------------------------------------------------------------
// -----------------------     SET UP (II)    ---------------------------------
// ----------------------------------------------------------------------------
// tagify() is the tag processor. It takes the text from 'uploadForm' after
// every 'input' event, splits it into an array of words, and checks each word
// for each prefix, adding it to that prefixes array as necessary:
function tagify() {
        // Delegate arrays to each type of prefix we wish to extract and/or
        // modify:
        //_______________________________________
        //__NAME/TYPE___________SYMBOL___________
        let tags     = [];  //    #     #tags
        let mentions = [];  //    @     @mentions
        let pointers = [];  //    &     &pointers 
        let stars    = [];  //    *     *stars   
        let cash     = [];  //    $     $cash
        let bangs    = [];  //    !     !bangs
        let evils    = [];  //    %     %evils
        let goods    = [];  //    ^     ^goods
        let wwwwwh   = [];  //    ?     ^wwwwwh
        // Grab the element we are typing in, and split the sentence string into 
        // an array words:
        let element  = document.getElementById(editable);
        let words    = element.innerText.split(" ");
        // Cycle through the words every time the user types and check for
        // keywords. If you don't cycle on every input event, you won't get the 
        // 'live' effect...:
        for (i=0; i<=words.length; i++) { if (sane(words[i])) { // Sanity check
                // Get the first character from each word and if it matches a symbol, 
                // push it to that symbols array (in Golang one might use a map, not 
                // sure how they work in js. 
                switch (words[i].charAt(0)) {
                        case "#": tags    .push(words[i].trim()); 
                        case "@": mentions.push(words[i].trim());
                        case "&": pointers.push(words[i].trim());
                        case "*": stars   .push(words[i].trim());
                        case "$": cash    .push(words[i].trim());
                        case "!": bangs   .push(words[i].trim());
                        case "%": evils   .push(words[i].trim());
                        case "^": goods   .push(words[i].trim());
                        case "?": wwwwwh  .push(words[i].trim());
                }}
        };
        // ________________________________________________________________________
        // ________________________________________________________________________
        // 'p' = "processing" - the string each word is added to after its
        // processed. We cycle through the [words] array again, checking each one 
        // to see which array, if any, contains them, and dress them in custom html
        // suits accordingly. Everything must be, and is, added back to `p`, in 
        // order:
        let p = ""; 
        // 'o' = opening section of html, added in front of each symbol:
        let o =` <span spellcheck="false" onclick="getTag('`;
        let z = `')" class="isym inline-`;
        // 'c' = closing section of html, added after each symbol:
        let c ='</span>'; 
        // Cycle through the words, check each symbols array for the word, and add
        // the appropriate html if the array includes the word:
        for (r=0; r<= words.length; r++) { if (sane(words[r])) { // Sanity check
                let w_ = words[r];       let w = w_.trim();
                if          (tags.includes(w)){p = p + o + w + z +     'tag"> ' + w_+ c;}
                else if (mentions.includes(w)){p = p + o + w + z + 'mention"> ' + w_+ c;}
                else if (pointers.includes(w)){p = p + o + w + z + 'pointer"> ' + w_+ c;}
                else if    (stars.includes(w)){p = p + o + w + z +    'star"> ' + w_+ c;}
                else if     (cash.includes(w)){p = p + o + w + z +    'cash"> ' + w_+ c;}
                else if    (bangs.includes(w)){p = p + o + w + z +    'bang"> ' + w_+ c;}
                else if    (evils.includes(w)){p = p + o + w + z +    'evil"> ' + w_+ c;}
                else if    (goods.includes(w)){p = p + o + w + z +    'good"> ' + w_+ c;}
                else if   (wwwwwh.includes(w)){p = p + o + w + z +  'wwwwwh"> ' + w_+ c;}
                else  /*    non-match     */  {p = p +              ' ' + w_   ;}}
        }
        // Add the newly processed string of html back to  the element. We do this
        // on every input, so we check for the empty string in case the user as
        // deleted something. I don't think it works properly if we don't wrap it 
        // in some other element:
        if (p.length > 0) { element.innerHTML = "<div>" + p + "</div>" }
        // and finally, move the cursor back to the end of the line (it moves it to
        // the beginning).
        //
        // TODO: When someone manually moves the cursor to the middle of the 
        // text and begins typing, the pointer is moved to the end of the text
        // after every keystroke. 
        moveCursorToEnd(element)
        let kws = "";
        // Cycle through the words, check each symbols array for the word, and add
        // the appropriate html if the array includes the word:
        for (r=0; r<= words.length; r++) { if (sane(words[r])) { // Sanity check
                let w_ = words[r];       let w = w_.trim();
                if          (tags.includes(w)){kws = kws + o + w + z +     'tag"> ' + w_ + c;}
                else if (mentions.includes(w)){kws = kws + o + w + z + 'mention"> ' + w_ + c;}
                else if (pointers.includes(w)){kws = kws + o + w + z + 'pointer"> ' + w_ + c;}
                else if    (stars.includes(w)){kws = kws + o + w + z +    'star"> ' + w_ + c;}
                else if     (cash.includes(w)){kws = kws + o + w + z +    'cash"> ' + w_ + c;}
                else if    (bangs.includes(w)){kws = kws + o + w + z +    'bang"> ' + w_ + c;}
                else if    (evils.includes(w)){kws = kws + o + w + z +    'evil"> ' + w_ + c;}
                else if    (goods.includes(w)){kws = kws + o + w + z +    'good"> ' + w_ + c;}
                else if   (wwwwwh.includes(w)){kws = kws + o + w + z +  'wwwwwh"> ' + w_ + c }};
        }
        if (kws.length > -1) {document.getElementById('upload-symbols')
                        .innerHTML = "<div>" + kws + "</div>"}
}
// Check for sanity:
function sane(w) {
        return (typeof w !== 'undefined' && w.length > 0);
}
// Move the cursor back to the end of the line (it moves to the beginning).
// TODO: When someone manually moves the cursor to the middle of the text and 
// begins typing, the pointer is moved to the end of the text after every
// keystroke. 
function moveCursorToEnd(elm) {
        elm.focus();
        const range = document.createRange();
        range.selectNodeContents(elm); range.collapse(false);
        const selection = window.getSelection();
        selection.removeAllRanges(); selection.addRange(range);
}
function getTag(tag) {
        console.log(tag)
}
