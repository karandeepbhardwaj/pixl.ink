function toggleTheme() {
    var html = document.documentElement;
    var current = html.getAttribute('data-theme');
    var next = current === 'dark' ? 'light' : 'dark';
    html.setAttribute('data-theme', next);
    localStorage.setItem('theme', next);
}

(function() {
    var saved = localStorage.getItem('theme') ||
        (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
    document.documentElement.setAttribute('data-theme', saved);
})();

function copyUrl() {
    var input = document.getElementById('shareUrl');
    if (input) {
        navigator.clipboard.writeText(input.value).then(function() {
            var btn = input.nextElementSibling;
            var original = btn.textContent;
            btn.textContent = 'Copied!';
            setTimeout(function() { btn.textContent = original; }, 2000);
        });
    }
}

var dropZone = document.getElementById('dropZone');
if (dropZone) {
    ['dragenter', 'dragover'].forEach(function(e) {
        dropZone.addEventListener(e, function(ev) {
            ev.preventDefault();
            dropZone.classList.add('dragover');
        });
    });
    ['dragleave', 'drop'].forEach(function(e) {
        dropZone.addEventListener(e, function() {
            dropZone.classList.remove('dragover');
        });
    });
}
