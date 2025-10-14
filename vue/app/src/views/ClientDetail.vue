<template>
  <div v-if="client" class="col" style="height:100%; gap:.75rem;">
    <header class="row" style="align-items:center; gap:.5rem;">
      <RouterLink class="btn ghost" to="/clients">← Back</RouterLink>
      <h2 style="margin:0">{{ client.name }}</h2>
      <span class="pill">{{ client.id }}</span>
      <div class="spacer"></div>
      <button class="btn" @click="enterGameMode">Enter Game Mode</button>
      <button class="btn ghost" @click="exitGameMode">Exit</button>
      <a class="btn ghost" :href="client.url" target="_blank" rel="noreferrer">Open in new tab</a>
    </header> 
    <div ref="rootEl" class="card" style="position:relative; flex:1; padding:0; overflow:hidden;">
      <iframe
      class="stream"
        ref="frameEl"
        :src="client.url"
        allow="fullscreen; pointer-lock; keyboard-lock"
        allowfullscreen
        style="border:0; width:100%; height:100%; background:#000;"
      ></iframe>

      <div v-if="hint" style="position:absolute; bottom:1rem; left:1rem; background:rgba(0,0,0,.65); padding:.5rem .75rem; border-radius:.5rem;">
        <small>{{ hint }}</small>
      </div>
    </div>
  </div>
  <div v-else class="col" style="gap:1rem;">
    <p>Client not found.</p>
    <RouterLink to="/clients" class="btn">Back</RouterLink>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref, computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { portal } from '../stores/portal'

const route = useRoute()
const id = String(route.params.id || '')
const client = computed(()=> portal.getClient(id))
const rootEl = ref(null)
const frameEl = ref(null)
const hint = ref('Click "Enter Game Mode" to enable fullscreen, pointer- and keyboard-lock. ESC exits.')

function log(){
  // eslint-disable-next-line prefer-rest-params
  const args = Array.from(arguments)
  console.debug('[portal]', ...args)
}

async function enterGameMode(){
  const root = rootEl.value
  const frame = frameEl.value
  if (!root || !frame) return
  try {
    if (!document.fullscreenElement && root.requestFullscreen) await root.requestFullscreen()
    frame.focus()
    if (root.requestPointerLock) root.requestPointerLock()
    if (navigator.keyboard && navigator.keyboard.lock) {
      await navigator.keyboard.lock([])
      hint.value = 'Game Mode enabled. Reserved OS keys (Windows/Meta, Alt+Tab, etc.) are not capturable.'
    } else {
      hint.value = 'Keyboard Lock not supported; most keys will work, but some browser shortcuts may still fire.'
    }
  } catch (e) {
    log('enterGameMode error', e)
  }
}

function exitGameMode(){
  if (navigator.keyboard && navigator.keyboard.unlock) navigator.keyboard.unlock()
  if (document.pointerLockElement && document.exitPointerLock) document.exitPointerLock()
  if (document.fullscreenElement && document.exitFullscreen) document.exitFullscreen()
  hint.value = 'Exited Game Mode.'
}

function onKeyDown(e){
  const key = String(e.key || '').toLowerCase()
  const ctrl = e.ctrlKey || e.metaKey
  const dangerous = (ctrl && ['w','s','r','p','t','n','l'].includes(key)) || ['f1','f5','f11','f12'].includes(key)
  if (dangerous) { e.preventDefault(); e.stopPropagation() }
}

function onFullscreenChange(){
  if (!document.fullscreenElement) {
    if (navigator.keyboard && navigator.keyboard.unlock) navigator.keyboard.unlock()
    if (document.pointerLockElement && document.exitPointerLock) document.exitPointerLock()
  }
}

onMounted(()=>{
  window.addEventListener('keydown', onKeyDown, { capture:true })
  document.addEventListener('fullscreenchange', onFullscreenChange)
})

onBeforeUnmount(()=>{
  window.removeEventListener('keydown', onKeyDown, { capture:true })
  document.removeEventListener('fullscreenchange', onFullscreenChange)
  exitGameMode()
})
</script>

<style scoped>
:host, .card { height: 100%; width: 100%;
  min-height: max-content; }

.stream {
  min-height:1280px;
}
</style>
