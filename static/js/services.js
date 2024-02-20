const address = window.location.hostname + ":" + window.location.port;
const socket = new WebSocket("ws:://" + address + "/service-sockets");

/**
 * @param {Event} event
 **/
function onOpenSocketCallback(event){

}

socket.onopen() = onOpenSocketCallback

