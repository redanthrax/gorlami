const container = document.getElementById('AgentWindow');

//we websockets up in the biotch

function handleOfferResponse(evt) {
    const offer = JSON.parse(evt.target.textContent);
    //createPeerConnection(offer);
}

const socket = new WebSocket('ws://localhost:3000/ws')

socket.onopen = function(event) {
    console.log("websocket connection open")
}

socket.onmessage = function(event) {
    // Handle received message
    const msg = JSON.parse(event.data);
    console.log(msg);
    switch(msg.type) {
        case "agentlist":
            console.log(msg.data);
            break;
        default:
            console.log("Unknown Type");
    }
};

socket.onclose = function(event) {
    // Connection closed
    console.log(event)
};

function sendMessage(message) {
    const payload = {
        type: 'message',
        content: message
    };

    socket.send(JSON.stringify(payload));
}

/*
document.addEventListener('htmx:afterSettle', handleOfferResponse);

function createPeerConnection(offer) {
    const peerConnection = new RTCPeerConnection();
    peerConnection.setRemoteDescription(offer)
    .then(() => {
        return peerConnection.createAnswer();
    })
    .then((answer) => {
        return peerConnection.setLocalDescription(answer);
    })
    .then(() => {
        //send our answer back and get going
        
        console.log(peerConnection.localDescription);
    })
    .catch((error) => {
        console.error(error);
    });
}
*/