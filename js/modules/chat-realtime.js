import apiClient from "../api/client.js";

// One connection per page, shared by the chat view and navigation badges.
let socket,
  retry,
  attempts = 0,
  stopped = false;
let connected = false;
const announce = (type, detail) =>
  window.dispatchEvent(new CustomEvent(type, { detail }));
function connect() {
  clearTimeout(retry);
  const token = apiClient.getToken();
  if (stopped || !token || (socket && socket.readyState < 2)) return;
  const url = new URL(apiClient.baseURL + "/chat/ws");
  url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
  const current = new WebSocket(url);
  socket = current;
  const timeout = setTimeout(() => current.close(), 10000);
  current.onopen = () => current.send(JSON.stringify({ token }));
  current.onmessage = (event) => {
    if (socket !== current) return;
    let message;
    try {
      message = JSON.parse(event.data);
    } catch {
      return;
    }
    if (message.type === "ready") {
      clearTimeout(timeout);
      connected = true;
      attempts = 0;
      announce("chat:connection", { connected: true });
      // Fetch durable history after every reconnect to recover missed events.
      announce("chat:updated", {});
    } else if (message.type === "chat.updated")
      announce("chat:updated", message);
  };
  current.onerror = () => current.close();
  current.onclose = (event) => {
    clearTimeout(timeout);
    if (socket !== current) return;
    connected = false;
    socket = null;
    announce("chat:connection", { connected: false });
    if (event.code === 1008) return; // Expired or revoked session: do not retry it.
    if (!stopped)
      retry = setTimeout(
        connect,
        Math.min(30000, 1000 * 2 ** attempts++) + Math.random() * 500,
      );
  };
}
window.addEventListener("pagehide", () => {
  stopped = true;
  clearTimeout(retry);
  socket?.close();
});
window.addEventListener("pageshow", () => {
  stopped = false;
  connect();
});
window.addEventListener("online", connect);
document.addEventListener("visibilitychange", () => {
  if (!document.hidden) {
    connect();
    announce("chat:updated", {});
  }
});
export default {
  start: connect,
  get connected() {
    return connected;
  },
};
