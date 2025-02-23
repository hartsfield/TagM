function hideAddressBar() {
        if (navigator.userAgent.match(/Android/i)) {
                //document.body.scrollTo(0, 10);
                var pageHeight = document.documentElement.clientHeight;
                if (window.outerHeight > pageHeight) {
                        document.body.style.height = window.outerHeight / window.devicePixelRatio + 'px';
                }
                //document.body.scrollTo(0, -200);
        }
        //} else {
        //        window.addEventListener("load", function() {
        //                setTimeout(function() { window.scrollTo(0, 1); }, 100);
        //        }, false);
        //}
}


const myElement = document.getElementById('stream');
let scrollTop2 = myElement.scrollTop;
let scrollTop3 = document.body.scrollTop;
let isdid = false;

//window.addEventListener("load", function() {
//        setTimeout(function() { 
//                if (!myElement.classList.contains('scrolled')) {
//                        myElement.classList.add('scrolled');
//                        //document.getElementById("stream").classList.add("stream-loaded");
//                }
//
//                hideAddressBar();
//                //document.getElementById("nav-outer").style.height = "500%";
//                window.scrollTo(0, 0); 
//                hideAddressBar();
//                window.scrollTo(0, 300); 
//                hideAddressBar();
//                window.scrollTo(0, 0); 
//                window.scrollTo(0, 300); 
//                hideAddressBar();
//                window.scrollTo(0, 0); 
//                window.scrollTo(0, 300); 
//                window.scrollTo(0, 0); 
//                hideAddressBar();
//                window.scrollTo(0, 300); 
//                window.scrollTo(0, 0); 
//                hideAddressBar();
//                window.scrollTo(0, 300); 
//                window.scrollTo(0, 0); 
//                hideAddressBar();
//                window.scrollTo(0, 300); 
//                window.scrollTo(0, 0); 
//                hideAddressBar();
//                window.scrollTo(0, 300); 
//                hideAddressBar();
//                window.scrollTo(0, 0); 
//                hideAddressBar();
//                window.scrollTo(0, 300); 
//                hideAddressBar();
//                window.scrollTo(0, 0); 
//                hideAddressBar();
//                window.scrollTo(0, 300); 
//                hideAddressBar();
//                window.scrollTo(0, 0); 
//                hideAddressBar();
//                window.scrollTo(0, 300); 
//        }, 200);
//
//}, false);


//window.onscroll = function() {
//        let scrollTop = document.body.scrollTop;
//        console.log(scrollTop, scrollTop3);
//        //window.scroll({top: 1153.5, left: 0, behavior: 'auto'});
//        //window.scroll({top: 0, left: 0, behavior: 'auto'});
//        //window.scroll({top: 1153.5, left: 0, behavior: 'auto'});
//        //window.scroll({top: 0, left: 0, behavior: 'auto'});
//        //window.scroll({top: 1153.5, left: 0, behavior: 'auto'});
//        //window.scroll({top: 0, left: 0, behavior: 'auto'});
//        //window.scroll({top: 1153.5, left: 0, behavior: 'auto'});
//        //window.scroll({top: 0, left: 0, behavior: 'auto'});
//        //window.scroll({top: 1153.5, left: 0, behavior: 'auto'});
//        //
//        if (scrollTop > 0 && scrollTop > scrollTop3 && !isdid) {         // scrolling down
//                // if (!myElement.classList.contains('scrolled')) {
//                //        //myElement.scroll({top: 90, left: 0, behavior: 'auto'});
//                //        myElement.classList.add('scrolled');
//                //        //requestFullscreen("hide")
//                //        //document.getElementById("stream").classList.add("stream-loaded");
//                //
//                //        //requestFullscreen()
//                //        //hideAddressBar();
//                //        //hideAddressBar();
//                //        (function(){
//                //                //document.body.scroll({top: 900, left: 0, behavior: 'auto'});
//                //                //document.body.scroll({top: 900, left: 0, behavior: 'auto'});
//                //                //myElement.scroll({top: 10, left: 0, behavior: 'smooth'});
//                //        })(); 
//                //}
//
//                hideAddressBar();
//                //document.body.scroll({top: 0, left: 0, behavior: 'auto'});
//        } 
//        scrollTop3 = document.body.scrollTop;
//}
//document.body.scroll({top: -10, left: 0, behavior: 'auto'});
//document.body.scroll({top: -10, left: 0, behavior: 'auto'});
//document.body.scroll({top: -10, left: 0, behavior: 'auto'});
//document.body.scroll({top: -10, left: 0, behavior: 'auto'});
//document.body.scroll({top: -10, left: 0, behavior: 'auto'});
//document.body.scroll({top: -10, left: 0, behavior: 'auto'});
//document.body.scroll({top: -10, left: 0, behavior: 'auto'});
//document.body.scroll({top: -10, left: 0, behavior: 'auto'});
//document.body.scroll({top: -10, left: 0, behavior: 'auto'});
//document.body.scroll({top: -10, left: 0, behavior: 'auto'});
//document.body.scroll({top: -10, left: 0, behavior: 'auto'});
//document.body.scroll({top: -10, left: 0, behavior: 'auto'});
window.onload = function() {
        //window.scrollTo(0, 5); 
        //setTimeout(function(){
        //        document.getElementById("ni").innerHTML = window.innerHeight;
        //        //document.getElementById("ni").style.display = "none";
        //}, 5000);

        //setTimeout(function(){
        //        myElement.scrollTo(0, 1000);
        //        document.body.scrollTo(0, 200);
        //        document.body.style.height = "110%";
        //        window.scrollTo(0, 200); 
        //}, 4000);

        let isdid2 = false;
        let initwh = window.innerHeight;
        let started = false;
        window.addEventListener("scroll", async function() {
                if (window.innerHeight > initwh) {
                        document.getElementById("ni").style.display = "none";
                        window.removeEventListener("scroll");
                        nextwh = window.innerHeight;
                        isdid = true;
                        return;
                }       else{
                        //if (window.innerHeight == initwh && !isdid) {
                        if (!started) {
                                setTimeout(function(){
                                        started = true;
                                        if (document.body.scrollTop >= 1) {
                                                window.scrollTo(0, 0);
                                                window.scrollTo(0, 30);
                                        } else {
                                                window.scrollTo(0, 30);
                                                window.scrollTo(0, 0);
                                        }
                                        started = false;
                                }, 200);
                        }
                }
                //isdid = true;
                //if (document.body.scrollTop <= 1) {
                //        window.scrollTo(0, 2);
                //} else {
                //        window.scrollTo(0, 0);
                //}
                //isdid = true;
                //}
        });
        setTimeout(function(){
                //document.body.scrollTo(0, 1);
                //window.scrollTo(0, 5); 

                window.scrollTo(0, 30);
        }, 200);

};
//myElement.addEventListener('scroll', function() {
//        document.body.scrollTo(0, 200);
//        window.scrollTo(0, 200); 
//
//});
//myElement.addEventListener('scroll', function() {
//        //hideAddressBar();
//        let scrollTop = myElement.scrollTop;
//        console.log(scrollTop, scrollTop2);
//        //myElement.scroll({top: 0, left: 0, behavior: 'auto'});
//        //document.body.scroll({top: 153.5, left: 0, behavior: 'auto'});
//        //isdid =true;
//        if (scrollTop > 50 && scrollTop > scrollTop2 && !isdid) {     
//                window.scroll({top: 1153.5, left: 0, behavior: 'auto'});
//                hideAddressBar();
//                if (!myElement.classList.contains('scrolled')) {
//                        //myElement.scroll({top: 90, left: 0, behavior: 'auto'});
//                        myElement.classList.add('scrolled');
//                        //requestFullscreen("hide")
//                        //document.getElementById("stream").classList.add("stream-loaded");
//
//                        //requestFullscreen()
//                        //hideAddressBar();
//                        //hideAddressBar();
//                        (function(){
//                                //document.body.scroll({top: 900, left: 0, behavior: 'auto'});
//                                //document.body.scroll({top: 900, left: 0, behavior: 'auto'});
//                                //myElement.scroll({top: 10, left: 0, behavior: 'smooth'});
//                        })(); 
//                }
//
//        }
//        scrollTop2 = myElement.scrollTop;
//});
//myElement.addEventListener('scroll', function() {
//        let scrollTop = myElement.scrollTop;
//        console.log(scrollTop, scrollTop2);
//        if (scrollTop > 10 && scrollTop > scrollTop2 && !isdid) {         // scrolling down
//                //
//                if (!myElement.classList.contains('scrolled')) {
//                        myElement.classList.add('scrolled');
//                        //document.getElementById("stream").classList.add("stream-loaded");
//                }
//                //document.getElementById("nav-outer").dispatchEvent(new Event("touchstart", { 
//                //        touches: [{ clientX: 0, clientY: 10 }] // Update touch coordinates as needed
//                //}));
//
//                //document.getElementById("nav-outer").dispatchEvent(new Event("touchmove", { 
//                //        touches: [{ clientX: 0, clientY: 1 }] // Update touch coordinates as needed
//                //}));
//                //document.getElementById("nav-outer").dispatchEvent(new Event("touchend", { 
//                //        touches: [{ clientX: 0, clientY: 0 }] // Update touch coordinates as needed
//                //}));
//
//                //myElement.dispatchEvent(new TouchEvent("touchmove", { 
//                //        touches: [{ clientX: 50, clientY: 20 }] // Update touch coordinates as needed
//                //}));
//                //myElement.dispatchEvent(new TouchEvent("touchmove", { 
//                //        touches: [{ clientX: 50, clientY: 5 }] // Update touch coordinates as needed
//                //}));
//                //myElement.dispatchEvent(new TouchEvent("touchmove", { 
//                //        touches: [{ clientX: 50, clientY: 20 }] // Update touch coordinates as needed
//                //}));
//                (function(){
//                        window.scroll({top: 501, left: 0, behavior: 'auto'});
        //                        //myElement.scroll({top: 10, left: 0, behavior: 'smooth'});
//                })(); 
//                (function(){
        //                        window.scroll({top: 21, left: 0, behavior: 'auto'});
//                        //myElement.scroll({top: 10, left: 0, behavior: 'smooth'});
//                        isdid = true;
        //                })();
//                //setTimeout(function(){
//                //        window.scroll({top: 0, left: 0, behavior: 'smooth'});
//                //        //myElement.scroll({top: 0, left: 0, behavior: 'smooth'});
//                //}, 30); 
//
        //                //setTimeout(function(){
        //                //        window.scroll({top: 12, left: 0, behavior: 'smooth'});
//                //        //myElement.scroll({top: 10, left: 0, behavior: 'smooth'});
        //                //}, 60); 
//                //setTimeout(function(){
//                //        window.scroll({top: 13, left: 0, behavior: 'smooth'});
//                //        //myElement.scroll({top: 10, left: 0, behavior: 'smooth'});
//                //}, 60);  setTimeout(function(){
        //                //        window.scroll({top: 14, left: 0, behavior: 'smooth'});
//                //        //myElement.scroll({top: 10, left: 0, behavior: 'smooth'});
//                //}, 60);  setTimeout(function(){
//                //        window.scroll({top: 15, left: 0, behavior: 'smooth'});
        //                //        //myElement.scroll({top: 10, left: 0, behavior: 'smooth'});
//                //}, 60);  setTimeout(function(){
//                //        window.scroll({top: 16, left: 0, behavior: 'smooth'});
        //                //        //myElement.scroll({top: 10, left: 0, behavior: 'smooth'});
//                //}, 60); 
        //                //setTimeout(function(){
        //                //        window.scroll({top: 200, left: 0, behavior: 'smooth'});
//                //        isdid = true;
//                //}, 1500); 
        //                //myElement.scroll({top: 10, left: 0, behavior: 'smooth'});
        //
//        } else {                                  // scrolling up
//                if (myElement.classList.contains('scrolled') && isdid) {
//                        myElement.classList.remove('scrolled');
        //                }
//        }
        //        scrollTop2 = myElement.scrollTop;
        //});

//const element = document.getElementById('bodyy');
        //
//// Create a dragstart event
//const dragStartEvent = new DragEvent('dragstart', {
        //  bubbles: true,
//  cancelable: true,
        //  dataTransfer: new DataTransfer()
//});
//
        //// Set data to be transferred
//dragStartEvent.dataTransfer.setData('text/plain', 'Some data');
        //
//// Dispatch the dragstart event
//element.dispatchEvent(dragStartEvent);
//
//// Handle the drag event
//element.addEventListener('dragstart', (event) => {
//  console.log('Drag started', event.dataTransfer.getData('text/plain'));
