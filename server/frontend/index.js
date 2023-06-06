const container = document.getElementById('AgentWindow');

function handleOfferResponse(evt) {
    const offer = JSON.parse(evt.target.textContent);
    //createPeerConnection(offer);
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