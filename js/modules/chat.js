import imagesAPI, { imageGallery, hydrateImages } from "../api/images.js";
import { initLayout } from "./layout.js";
import chatRealtime from "./chat-realtime.js";
import authManager from "./auth-manager.js";
import chatAPI from "../api/chat.js";
import { formatDate } from "../utils/helpers.js";

if (!(await initLayout())) throw new Error("Authentication required");
const admin = authManager.getUserRole() === "admin";
const $ = (id) => document.getElementById(id);
const escape = (value) =>
  String(value ?? "").replace(
    /[&<>"']/g,
    (c) =>
      ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" })[
        c
      ],
  );
const tr = (text) => window.SiteI18n?.translate(text) || text;
let people = [],
  selected = null,
  thread = null,
  messages = [],
  filter = "all",
  selectionVersion = 0,
  busy = false,
  loading = false,
  olderAvailable = false,
  pendingID = null;
const drafts = new Map();
let attachments = [];
function error(message = "") {
  $("chatError").textContent = message;
  $("chatError").hidden = !message;
  $("chatRetry").hidden = !message;
}
function drawPeople() {
  const search = $("chatSearch").value.toLowerCase().trim();
  const list = people.filter(
    (p) =>
      (!search || `${p.username} ${p.email}`.toLowerCase().includes(search)) &&
      (filter === "all" ||
        (filter === "unread" ? p.unread > 0 : p.status === filter)),
  );
  $("chatTotal").textContent = people.length;
  $("chatPeople").innerHTML = list.length
    ? list
        .map(
          (p) =>
            `<button type="button" class="chat-person" data-user="${escape(p.id)}" aria-current="${selected?.id === p.id}"><span class="chat-avatar" aria-hidden="true">${escape(p.username?.slice(0, 1).toUpperCase())}</span><span class="chat-person-body"><strong class="chat-person-name" data-no-translate>${escape(p.username)}</strong><span class="chat-person-role">${escape(tr({ resident: "Cư dân", manager: "Ban quản lý", staff: "Nhân viên" }[p.role] || p.role))} · ${escape(p.account_status === "active" ? "Đang hoạt động" : "Không hoạt động")}</span><span class="chat-person-preview" ${p.preview ? "data-no-translate" : ""}>${escape(p.preview || "Chưa có tin nhắn")}</span><span class="chat-person-status">${escape(p.status === "resolved" ? "Đã giải quyết" : p.status === "open" ? "Đang mở" : "")}</span></span>${p.unread ? `<span class="chat-unread" aria-label="${escape(tr("Chưa đọc"))}">${p.unread}</span>` : ""}</button>`,
        )
        .join("")
    : '<p class="chat-empty">Không tìm thấy cuộc trò chuyện</p>';
  $("chatPeople")
    .querySelectorAll("[data-user]")
    .forEach((button) =>
      button.addEventListener("click", () =>
        selectPerson(people.find((p) => p.id === button.dataset.user)),
      ),
    );
}
async function refreshPeople() {
  const res = await chatAPI.people();
  people = res.data || [];
  drawPeople();
}
function drawMessages(stick = false, older = false) {
  const history = $("chatHistory"),
    oldHeight = history.scrollHeight,
    oldTop = history.scrollTop;
  $("chatMessages").innerHTML = messages.length
    ? messages
        .map(
          (m) =>
            `<div class="chat-message-row ${m.from_admin === admin ? "mine" : ""}"><div class="chat-message"><div class="chat-message-body" data-no-translate>${escape(m.body)}</div>${imageGallery(m.image_ids)}<time class="chat-message-time" datetime="${escape(m.created_at)}">${escape(formatDate(m.created_at, "dd/mm/yyyy hh:mm"))}</time></div></div>`,
        )
        .join("")
    : '<p class="chat-empty">Bắt đầu cuộc trò chuyện. Chúng tôi luôn sẵn lòng lắng nghe.</p>';
  hydrateImages($("chatMessages"));
  $("chatOlder").hidden = !olderAvailable;
  if (older) history.scrollTop = history.scrollHeight - oldHeight + oldTop;
  else if (stick) history.scrollTop = history.scrollHeight;
  $("chatStatus").textContent =
    thread?.status === "resolved"
      ? "Mở lại cuộc trò chuyện"
      : "Đánh dấu đã giải quyết";
}
async function refreshMessages(initial = false) {
  if (!thread || loading) return;
  loading = true;
  const id = thread.id,
    version = selectionVersion;
  try {
    const history = $("chatHistory");
    const stick =
      initial ||
      history.scrollHeight - history.scrollTop - history.clientHeight < 80;
    const res = await chatAPI.messages(id);
    if (version !== selectionVersion) return;
    const known = new Map(messages.map((m) => [m.id, m]));
    (res.data || []).forEach((m) => known.set(m.id, m));
    const previousCount = messages.length;
    messages = [...known.values()].sort(
      (a, b) =>
        a.created_at.localeCompare(b.created_at) || a.id.localeCompare(b.id),
    );
    thread = res.conversation;
    if (initial || messages.length <= 100) olderAvailable = res.has_more;
    drawMessages(stick);
    if (previousCount && messages.length > previousCount)
      $("chatAnnouncement").textContent = tr("Có tin nhắn mới");
    if (stick && messages.length && !document.hidden)
      await chatAPI.read(id, messages.at(-1).id);
  } finally {
    loading = false;
  }
}
async function selectPerson(person) {
  if (busy) return;
  if (thread)
    drafts.set(thread.id, {
      body: $("chatDraft").value,
      pendingID,
      attachments: [...attachments],
    });
  const version = ++selectionVersion;
  selected = person;
  thread = null;
  messages = [];
  pendingID = null;
  attachments = [];
  drawAttachments();
  loading = false;
  error();
  $("chatWelcome").hidden = true;
  $("chatActive").hidden = false;
  $("chatMessages").innerHTML = '<p class="chat-empty">Đang tải...</p>';
  $("chatName").textContent = person?.username || tr("Ban quản trị Vinhomes");
  $("chatAvatar").textContent =
    person?.username?.slice(0, 1).toUpperCase() || "V";
  $("chatPersonInfo").textContent = person
    ? `${person.email} · ${tr(person.account_status === "active" ? "Đang hoạt động" : "Không hoạt động")}`
    : tr("Hỗ trợ cư dân");
  $("chatStatus").hidden = !admin;
  $("chatDraft").value = "";
  $("chatSend").disabled = true;
  if (admin) drawPeople();
  try {
    const opened = person?.conversation_id
      ? { id: person.conversation_id }
      : (await chatAPI.start(person?.id)).data;
    if (version !== selectionVersion) return;
    thread = opened;
    $("chatDraft").value = drafts.get(thread.id)?.body || "";
    pendingID = drafts.get(thread.id)?.pendingID || null;
    attachments = [...(drafts.get(thread.id)?.attachments || [])];
    drawAttachments();
    $("chatDraftCount").textContent = `${$("chatDraft").value.length} / 2000`;
    await refreshMessages(true);
    if (admin) await refreshPeople();
  } catch (e) {
    error(e.message);
  } finally {
    if (version === selectionVersion)
      $("chatSend").disabled =
        !thread || Boolean(person && person.account_status !== "active");
  }
}
function drawAttachments() {
  $("chatDraft").required = !attachments.length;
  $("chatAttachments").innerHTML = attachments
    .map(
      (id) =>
        `<div>${imageGallery([id])}<button type="button" data-remove-image="${id}" aria-label="Remove image">×</button></div>`,
    )
    .join("");
  hydrateImages($("chatAttachments"));
  $("chatAttachments")
    .querySelectorAll("[data-remove-image]")
    .forEach((button) =>
      button.addEventListener("click", async () => {
        if (busy) return;
        try {
          await imagesAPI.remove(button.dataset.removeImage);
          attachments = attachments.filter(
            (id) => id !== button.dataset.removeImage,
          );
          pendingID = null;
          drawAttachments();
        } catch (e) {
          error(e.message);
        }
      }),
    );
}
$("chatAttach").addEventListener("click", () => {
  if (!busy && thread) $("chatImageInput").click();
});
$("chatImageInput").addEventListener("change", async (event) => {
  if (busy || !thread) return;
  const files = [...event.target.files];
  event.target.value = "";
  if (files.length + attachments.length > 5) {
    error(tr("Choose up to 5 images"));
    return;
  }
  busy = true;
  $("chatSend").disabled = true;
  $("chatAttach").disabled = true;
  error();
  try {
    for (const file of files) {
      const res = await imagesAPI.upload(file, "chat", thread.id);
      attachments.push(res.data.id);
      pendingID = null;
      drawAttachments();
    }
  } catch (e) {
    error(e.message);
  } finally {
    busy = false;
    $("chatSend").disabled = false;
    $("chatAttach").disabled = false;
  }
});
$("chatSearch").addEventListener("input", drawPeople);
document.querySelectorAll("[data-filter]").forEach((button) =>
  button.addEventListener("click", () => {
    filter = button.dataset.filter;
    document
      .querySelectorAll("[data-filter]")
      .forEach((b) => b.setAttribute("aria-pressed", String(b === button)));
    drawPeople();
  }),
);
$("chatDraft").addEventListener("input", () => {
  $("chatDraftCount").textContent = `${$("chatDraft").value.length} / 2000`;
  pendingID = null;
});
$("chatForm").addEventListener("submit", async (event) => {
  event.preventDefault();
  const body = $("chatDraft").value.trim();
  if (!thread || busy || (!body && !attachments.length)) return;
  busy = true;
  error();
  $("chatSend").disabled = true;
  pendingID ||= crypto.randomUUID();
  try {
    const res = await chatAPI.send(thread.id, body, pendingID, attachments);
    if (!messages.some((m) => m.id === res.data.id)) messages.push(res.data);
    $("chatDraft").value = "";
    drafts.delete(thread.id);
    attachments = [];
    drawAttachments();
    pendingID = null;
    $("chatDraftCount").textContent = "0 / 2000";
    thread.status = "open";
    drawMessages(true);
    $("chatDraft").focus();
    if (admin) await refreshPeople();
  } catch (e) {
    error(e.message);
  } finally {
    busy = false;
    $("chatSend").disabled = false;
  }
});
$("chatOlder").addEventListener("click", async () => {
  if (!thread || !messages.length || loading) return;
  loading = true;
  const version = selectionVersion;
  $("chatOlder").disabled = true;
  try {
    const res = await chatAPI.messages(thread.id, messages[0].id);
    if (version !== selectionVersion) return;
    messages = [...(res.data || []), ...messages];
    olderAvailable = res.has_more;
    drawMessages(false, true);
  } catch (e) {
    error(e.message);
  } finally {
    loading = false;
    $("chatOlder").disabled = false;
  }
});
$("chatStatus").addEventListener("click", async () => {
  if (!thread) return;
  const id = thread.id,
    version = selectionVersion,
    next = thread.status === "resolved" ? "open" : "resolved";
  $("chatStatus").disabled = true;
  try {
    await chatAPI.status(id, next);
    if (version === selectionVersion) {
      thread.status = next;
      drawMessages();
    }
    await refreshPeople();
  } catch (e) {
    error(e.message);
  } finally {
    $("chatStatus").disabled = false;
  }
});
$("chatRetry").addEventListener("click", async () => {
  error();
  try {
    if (!thread && (!admin || selected)) await selectPerson(selected);
    else {
      await refreshMessages();
      if (admin) await refreshPeople();
    }
  } catch (e) {
    error(e.message);
  }
});
window.addEventListener("site:languagechange", () => {
  if (!admin) $("chatName").textContent = tr("Ban quản trị Vinhomes");
  if (admin) drawPeople();
  if (thread) drawMessages();
});
if (admin) {
  $("chatHelp").textContent =
    "Tin nhắn chỉ hiển thị cho người dùng này và quản trị viên.";
  $("chatDirectory").hidden = false;
  $("chatWorkspace").classList.add("admin-chat");
  $("chatPageTitle").textContent = "Trung tâm trò chuyện";
  $("chatSubtitle").textContent =
    "Kết nối với cư dân, theo dõi tin chưa đọc và quản lý mọi cuộc trò chuyện.";
  try {
    await refreshPeople();
  } catch (e) {
    error(e.message);
  }
} else await selectPerson(null);
let refreshing = false,
  refreshQueued = false;
async function refreshChat() {
  if (document.hidden || refreshing || busy) return;
  refreshing = true;
  try {
    await refreshMessages();
    if (admin) await refreshPeople();
  } catch (e) {
    error(e.message);
  } finally {
    refreshing = false;
  }
}
window.addEventListener("chat:updated", () => {
  if (refreshQueued) return;
  refreshQueued = true;
  const run = async () => {
    if (busy || loading || refreshing) {
      setTimeout(run, 100);
      return;
    }
    refreshQueued = false;
    await refreshChat();
  };
  run();
});
setInterval(() => {
  if (!chatRealtime.connected) refreshChat();
}, 5000);
