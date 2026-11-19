/** Override before loading API modules when using a different backend. */
if (!window.__API_BASE_URL) {
  const local = ['localhost', '127.0.0.1', '[::1]'].includes(window.location.hostname);
  window.__API_BASE_URL = local
    ? 'http://127.0.0.1:8080/api/v1'
    : 'https://3d-geographic-information-system-production.up.railway.app/api/v1';
}
