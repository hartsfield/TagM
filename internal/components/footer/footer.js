const myElement = document.getElementById('stream');

myElement.addEventListener('scroll', function() {
        // Get the current vertical scroll position
        const scrollTop = myElement.scrollTop;

        // Log the scroll position (for example)
        console.log('Scrolled to:', scrollTop);

        // Example: Change element style based on scroll position
        if (scrollTop > 10) {
                myElement.classList.add('scrolled');
                (function() {
                        window.scroll({
                                top: 500,
                                left: 0,
                                behavior: 'auto'
                        });
                        document.getElementById("stream").classList.add("stream-loaded");
                })();

        } else {
                myElement.classList.remove('scrolled');
        }
});

