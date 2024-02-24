console.log("hello");
const address = window.location.hostname + ":" + window.location.port;
const socket = new WebSocket("ws://" + address + "/services-socket");

/**
 * @param {Event} event
 **/
function onOpenSocketCallback(event) {
  console.log("Connection opened");
}

/**
 * @param {MessageEvent} event
 **/
function onMessageSocketCallback(event) {
  console.log("Received message from server:");
  var newService = JSON.parse(event.data);
  console.log(newService);
}

/**
 * @param {CloseEvent} event
 **/
function onCloseSocketCallback(event) {
  console.log("WebSocket connection closed.");
}

// Event handler for when the WebSocket connection is established
socket.onopen = onOpenSocketCallback;

// Event handler for when a message is received from the server
socket.onmessage = onMessageSocketCallback;

// Event handler for when the WebSocket connection is closed
socket.onclose = onCloseSocketCallback;
