// If you're on JS instead of TS, use this file name:
// File: src/composables/useNotifications.js
export function notify(text, isError = false, ttl = 5000) {
  window.dispatchEvent(new CustomEvent('notify', { detail: { text, isError, ttl } }));
  const showResponse = (t, e = false, time = 5000) => notify(t, e, time);
  return { showResponse };
}

export function useNotifications() {
  const showResponse = (text, isError = false, ttl = 5000) => notify(text, isError, ttl);
  return { notify, showResponse };
}

export default notify;
