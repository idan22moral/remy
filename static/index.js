/**
 * @type {WebSocket | null | undefined}
 */
let socket;
let isMouseOn = false;

startMovementTracker();

function startMovementTracker() {
    const gyroscope = new Gyroscope({ frequency: 60 });

    gyroscope.addEventListener("reading", () => {
        const { x, z } = gyroscope;
        sendMovement(z, x);
    });

    gyroscope.start();
}

async function sendMovement(x, y) {
    await sendMessage({ type: 'movement', x, y });
}

async function sendRightClick() {
    await sendMessage({ type: 'click', button: 'right' });
}

async function sendLeftClick() {
    await sendMessage({ type: 'click', button: 'left' });
}

function switchMouseOnOff() {
    isMouseOn = !isMouseOn;

    if (isMouseOn) {
        socket = initWebSocket();
    } else {
        socket.close();
        socket = null;
    }
}

async function sendMessage(message) {
    if (socket?.readyState !== WebSocket.OPEN) {
        return;
    }

    const serializedMessage = JSON.stringify(message);
    socket.send(serializedMessage);
}

function initWebSocket() {
    socket = new WebSocket('/ws');

    socket.onopen = function () {
        console.log("WebSocket connection opened");
    };

    socket.onmessage = function (event) {
        const data = JSON.parse(event.data);
        console.log("Received message from server:", data);
    };

    socket.onerror = function (error) {
        console.error("WebSocket error:", error);
    };

    socket.onclose = function () {
        console.log("WebSocket connection closed");
    };

    return socket;
}
