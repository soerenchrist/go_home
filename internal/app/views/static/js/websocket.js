var ws;
var callbacks = [];
window.addEventListener("load", function () {
  // get current host
  var host = window.location.host;
  console.log(host);
  ws = new WebSocket("ws://" + host + "/ws");
  ws.onopen = function () {
    console.log("OPEN");
  };
  ws.onclose = function () {
    console.log("CLOSE");
  };

  ws.onmessage = function (event) {
    var data = JSON.parse(event.data);
    for (let callback of callbacks) {
      callback(data);
    }
  };
});

window.addEventListener("unload", function () {});

function registerToWebsocket(callback) {
  callbacks.push(callback);
}
