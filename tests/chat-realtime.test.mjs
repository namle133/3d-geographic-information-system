import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import vm from "node:vm";
function client() {
  const sockets = [],
    events = [],
    listeners = {},
    timers = new Map();
  let seq = 0,
    token = "private-token";
  class Socket {
    constructor(url) {
      this.url = String(url);
      this.readyState = 0;
      sockets.push(this);
    }
    send(data) {
      this.sent = JSON.parse(data);
    }
    close() {
      this.readyState = 3;
      this.onclose?.({ code: 1000 });
    }
  }
  const scope = {
    WebSocket: Socket,
    URL,
    Math,
    apiClient: { baseURL: "https://api.example/api/v1", getToken: () => token },
    setTimeout: (fn, delay) => {
      timers.set(++seq, { fn, delay });
      return seq;
    },
    clearTimeout: (id) => timers.delete(id),
    window: {
      dispatchEvent: (event) => events.push(event),
      addEventListener: (name, fn) => (listeners[name] = fn),
    },
    document: {
      hidden: false,
      addEventListener: (name, fn) => (listeners[name] = fn),
    },
    CustomEvent: class {
      constructor(type, options) {
        this.type = type;
        this.detail = options.detail;
      }
    },
  };
  let source = fs
    .readFileSync(
      new URL("../js/modules/chat-realtime.js", import.meta.url),
      "utf8",
    )
    .replace(/^import .*;\s*$/gm, "")
    .replace("export default", "globalThis.client =");
  vm.runInNewContext(source, scope);
  return {
    api: scope.client,
    sockets,
    events,
    listeners,
    timers,
    token: (value) => (token = value),
  };
}
test("authenticates in a frame and waits for ready before claiming a live connection", () => {
  const c = client();
  c.api.start();
  c.api.start();
  assert.equal(c.sockets.length, 1);
  const s = c.sockets[0];
  assert.equal(s.url, "wss://api.example/api/v1/chat/ws");
  assert.doesNotMatch(s.url, /private-token/);
  s.onopen();
  assert.equal(s.sent.token, "private-token");
  assert.equal(c.api.connected, false);
  s.onmessage({ data: JSON.stringify({ type: "ready" }) });
  assert.equal(c.api.connected, true);
  assert.equal(c.timers.size, 0);
  s.onmessage({
    data: JSON.stringify({ type: "chat.updated", conversation_id: "thread" }),
  });
  assert.equal(c.events.at(-1).detail.conversation_id, "thread");
});
test("reconnects with the current token and requests catch-up; pagehide stops retries", () => {
  const c = client();
  c.api.start();
  const first = c.sockets[0];
  first.close();
  c.token("renewed");
  [...c.timers.values()][0].fn();
  assert.equal(c.sockets.length, 2);
  const second = c.sockets[1];
  second.onopen();
  assert.equal(second.sent.token, "renewed");
  second.onmessage({ data: '{"type":"ready"}' });
  assert.equal(c.events.at(-1).type, "chat:updated");
  c.listeners.pagehide();
  assert.equal(c.timers.size, 0);
  assert.equal(c.api.connected, false);
});
test("authentication failures do not enter a retry loop", () => {
  const c = client();
  c.api.start();
  c.sockets[0].onclose({ code: 1008 });
  assert.equal(c.timers.size, 0);
  assert.equal(c.api.connected, false);
});
