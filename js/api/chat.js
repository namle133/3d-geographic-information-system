import apiClient from "./client.js";
export default {
  unread: () => apiClient.get("/chat/unread"),
  people: () => apiClient.get("/chat/people"),
  start: (userID) =>
    apiClient.post("/chat/conversations", userID ? { user_id: userID } : {}),
  messages: (id, before) =>
    apiClient.get(
      `/chat/conversations/${encodeURIComponent(id)}/messages${before ? "?before=" + encodeURIComponent(before) : ""}`,
    ),
  send: (id, body, clientID) =>
    apiClient.post(`/chat/conversations/${encodeURIComponent(id)}/messages`, {
      body,
      client_id: clientID,
    }),
  read: (id, messageID) =>
    apiClient.put(`/chat/conversations/${encodeURIComponent(id)}/read`, {
      message_id: messageID,
    }),
  status: (id, status) =>
    apiClient.put(`/chat/conversations/${encodeURIComponent(id)}/status`, {
      status,
    }),
};
