<template> 
  <div class="notification-container" role="region" aria-live="polite" aria-label="Notifications">
    <div
      v-for="msg in messages"
      :key="msg.id"
      class="notification-box"
      :class="msg.isError ? 'error' : 'info'"
    >
      <span class="text">{{ msg.text }}</span>
      <button class="close-btn" @click="close(msg.id)" aria-label="Dismiss">×</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';

const messages = ref([]);

/* Simple unique id */
const uid = () => `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;

function close(id) {
  messages.value = messages.value.filter(m => m.id !== id);
}

function onNotify(e) {
  const d = e?.detail;
  if (!d || typeof d.text !== 'string') return; 
  const id = uid();
  messages.value.push({ id, text: d.text, isError: !!d.isError });
  if (messages.value.length > 6) messages.value.shift();  
  const ttl = Math.max(1000, Number(d.ttl ?? 5000));
  setTimeout(() => close(id), ttl);
}

onMounted(() => window.addEventListener('notify', onNotify));
onBeforeUnmount(() => window.removeEventListener('notify', onNotify));
</script>

<style scoped> 
.notification-container {
  position: fixed;
  bottom: 10%;
  right: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-height: 80vh;
  overflow-y: auto;
  pointer-events: none;    
  z-index: 2147483646;  
}

.notification-box {
  pointer-events: auto; 
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 16px;
  min-width: 260px;
  max-width: 360px;
  border-radius: 10px;
  color: #fff;
  background: #01bb8d;
  box-shadow: 0 8px 24px rgba(0,0,0,0.25);
}

.notification-box.error { background: #ef4444; }

.text { flex: 1; word-break: break-word; }

.close-btn {
  appearance: none;
  background: transparent;
  border: 0;
  color: #fff;
  font-size: 18px;
  line-height: 1;
  cursor: pointer;
  padding: 4px 6px;
  border-radius: 6px;
}
.close-btn:hover { background: rgba(255,255,255,0.15); }
</style>
