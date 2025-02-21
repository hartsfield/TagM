const myElement = document.getElementById('stream');

myElement.addEventListener('scroll', function() {
        // Get the current vertical scroll position
        const scrollTop = myElement.scrollTop;

        // Log the scroll position (for example)

        // Example: Change element style based on scroll position
        console.log('Scrolled to:', scrollTop);
        if (scrollTop > 50) {
                myElement.classList.add('scrolled');
                (function() {
                        window.scroll({
                                top: 100,
                                left: 0,
                                behavior: 'auto'
                        });

        console.log('wffewfew');
                        document.getElementById("stream").classList.add("stream-loaded");
                })();

        } else {
                myElement.classList.remove('scrolled');
        }
});

