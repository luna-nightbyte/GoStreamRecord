import { reactive, computed } from 'vue'

const STORAGE_KEY = 'portal-settings-v1'

function load() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return { clients: [] }
    const data = JSON.parse(raw)
    if (!Array.isArray(data.clients)) return { clients: [] }
    return data
  } catch {
    return { clients: [] }
  }
}

function persist(state) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state))
}

function newId() {
  const d = new Date()
  const pad = (n)=> n.toString().padStart(2,'0')
  const stamp = `${d.getFullYear()}${pad(d.getMonth()+1)}${pad(d.getDate())}${pad(d.getHours())}${pad(d.getMinutes())}${pad(d.getSeconds())}`
  const rand = Math.random().toString(36).slice(2,6)
  return `${stamp}-${rand}`
}

const state = reactive(load())

export const portal = {
  state,
  clients: computed(()=> state.clients),
  getClient(id) {
    return state.clients.find(c=>c.id===id)
  },
  addClient(input = {}) {
    const name = (input.name || '').trim() || 'TNV'
    const url = (input.url || '').trim() || '/api/stream'
    const id = (input.id || '').trim() || newId()
    state.clients.push({ id, name, url })
    persist(state)
    return id
  },
  updateClient(id, patch = {}) {
    const c = state.clients.find(c=>c.id===id)
    if (!c) return false
    if (typeof patch.name === 'string') c.name = patch.name
    if (typeof patch.url === 'string') c.url = patch.url
    if (typeof patch.id === 'string' && patch.id && patch.id!==id) {
      if (!state.clients.some(x=>x.id===patch.id)) c.id = patch.id
    }
    persist(state)
    return true
  },
  removeClient(id) {
    const i = state.clients.findIndex(c=>c.id===id)
    if (i>=0) {
      state.clients.splice(i,1)
      if (state.lastSelectedClientId===id) state.lastSelectedClientId = undefined
      persist(state)
      return true
    }
    return false
  },
  selectClient(id) {
    state.lastSelectedClientId = id
    persist(state)
  },
  exportJson() {
    return JSON.stringify(state, null, 2)
  },
  importJson(raw) {
    const data = JSON.parse(raw)
    if (!data || !Array.isArray(data.clients)) throw new Error('Invalid settings file')
    state.clients.splice(0, state.clients.length, ...data.clients)
    state.lastSelectedClientId = data.lastSelectedClientId
    persist(state)
  }
}
