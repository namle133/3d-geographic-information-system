import apiClient from "./client.js";
const uuid = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
const cache = new Map();
const tr = (text) => window.SiteI18n?.translate(text) || text;
export function validateImage(file) {
  if (
    !["image/jpeg", "image/png"].includes(file.type) ||
    file.size > 2 * 1024 * 1024
  )
    throw new Error(tr("Upload a JPEG or PNG image up to 2 MB"));
}
export function imageGallery(ids = []) {
  return (
    '<div class="stored-image-gallery">' +
    ids
      .filter((id) => uuid.test(id))
      .map(
        (id) =>
          `<a target="_blank" rel="noopener" class="stored-image-link"><img data-image-id="${id}" alt="${tr("Attached image")}" loading="lazy" /></a>`,
      )
      .join("") +
    "</div>"
  );
}
export async function hydrateImages(root) {
  await Promise.allSettled(
    [...root.querySelectorAll("img[data-image-id]")].map(async (img) => {
      const id = img.dataset.imageId;
      if (!uuid.test(id)) return;
      if (!cache.has(id))
        cache.set(
          id,
          apiClient
            .get("/images/" + id)
            .then((r) => r.blob())
            .then((blob) => URL.createObjectURL(blob))
            .catch((error) => {
              cache.delete(id);
              throw error;
            }),
        );
      try {
        const url = await cache.get(id);
        if (!img.isConnected) return;
        img.src = url;
        img.closest("a").href = url;
      } catch {
        img.alt = tr("Unable to load image");
      }
    }),
  );
}
window.addEventListener("pagehide", () => {
  for (const promise of cache.values())
    promise.then((url) => URL.revokeObjectURL(url)).catch(() => {});
  cache.clear();
});
export default {
  async upload(file, purpose, conversationID = "") {
    validateImage(file);
    const form = new FormData();
    form.append("file", file);
    form.append("purpose", purpose);
    form.append("conversation_id", conversationID);
    return apiClient.fetch("/images", {
      method: "POST",
      body: form,
      headers: { Authorization: "Bearer " + apiClient.getToken() },
    });
  },
  remove: (id) => apiClient.delete("/images/" + encodeURIComponent(id)),
};
