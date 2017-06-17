hljs.initHighlightingOnLoad();
var xhr = new XMLHttpRequest();
var xhrPublish = new XMLHttpRequest();
var path = window.location.pathname.split("/");
var blogId = path[path.length-1];

xhr.onload = function() {
    document.getElementById("view").innerHTML = this.responseText;

    if (this.responseText.includes("![]()")) {
        modifiedText = this.responseText.replace("![]()", "<form action='/file/post' class='dropzone' id='fileUpload'></form>");
        document.getElementById("view").innerHTML = modifiedText;
        var imageUpload = new Dropzone(".dropzone", { url: "/file/post/"+blogId});
        imageUpload.on("success", function(file, responseText) { 
            document.getElementById("blogText").value = document.getElementById("blogText").value.replace("![]()", "![]("+JSON.parse(responseText).fileName+")")
            sendPost();
        });
    }

    nodeList = document.getElementsByTagName("pre");
    nodes = Array.prototype.slice.call(nodeList, 0); 
    nodes.forEach(function(item, index) {
        codeList = item.getElementsByTagName("code")
        codes = Array.prototype.slice.call(codeList, 0);
        codes.forEach(function(code, index) {
            hljs.highlightBlock(code);
        });
    });
};

var sendPost = debounce(function (e) {
        send(blogId, 
            document.getElementById("blogTitle").value, 
            document.getElementById("blogText").value);
    }, 1200, false
);
document.getElementById("blogText").addEventListener('keyup', sendPost);

function debounce(func, threshold, execAsap) {
    var timeout;
    return function debounced () {
        var obj = this, args = arguments;
        function delayed () {
            if (!execAsap)
                func.apply(obj, args);
            timeout = null; 
        };

        if (timeout)
            clearTimeout(timeout);
        else if (execAsap)
            func.apply(obj, args);

        timeout = setTimeout(delayed, threshold || 100); 
    };
}

function send(blogId, blogTitle, blogText) {
    xhr.open('POST', '/push/post/' + blogId, true);
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.send('blogId=' + blogId+ '&blogTitle=' + blogTitle + '&blogText=' + blogText);
}

document.getElementById("publish").addEventListener('click', publish);

function publish() {
    blogTitle = document.getElementById("blogTitle").value;
    blogText = document.getElementById("blogText").value;
    xhrPublish.open('POST', '/push/publish/' + blogId, true);
    xhrPublish.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhrPublish.send('blogId=' + blogId+ '&blogTitle=' + blogTitle + '&blogText=' + blogText);
}

xhrPublish.onload = function() {
    document.getElementById("published").innerHTML = this.responseText;
};

window.addEventListener("load", function(event) {
    console.log("aaa", document.getElementById("blogText").value);
    if (document.getElementById("blogText").value != null) {
        console.log("bbb", document.getElementById("blogText").value != null);
        send(blogId, 
            document.getElementById("blogTitle").value, 
            document.getElementById("blogText").value);
    }
});