console.log("hello")
const address = window.location.hostname + ":" + window.location.port;
const socket = new WebSocket("ws:://" + address + "/services-socket");

/**
 * @param {Event} event
 **/
function onOpenSocketCallback(event){

}

/**
 * @param {Event} event
 **/
function onMessageSocketCallback(event) {
    console.log("Received message from server:", event.data);
};

/**
 * @param {Event} event
 **/
function onCloseSocketCallback(event) {
    console.log("WebSocket connection closed.");
};



// Event handler for when the WebSocket connection is established
socket.onopen() = onOpenSocketCallback

// Event handler for when a message is received from the server
socket.onmessage =onMessageSocketCallback

// Event handler for when the WebSocket connection is closed
socket.onclose = onCloseSocketCallback
