<template>
  <div class="col" style="gap:1rem; max-width:980px; margin:0 auto;">
    <header class="row" style="align-items:center;">
      <h2 class="txt40pt" style="margin:0">Settings</h2>
      <div class="spacer"></div>
      <button class="btn button_form" @click="addNew">Add new client</button>
      <button class="btn button_form ghost" @click="exportFile">Export JSON</button>
      <label class="btn button_form ghost" style="cursor:pointer; display:inline-flex; align-items:center; gap:.4rem;">
        Import JSON
        <input type="file" accept="application/json" @change="importFile" style="display:none" />
      </label>
    </header>

    <section class="card" >
      <table>
        <thead>
          <tr><th  style="width:14rem;">Name</th><th style="width:16rem;">ID</th><th>URL</th><th style="width:10rem;"></th></tr>
        </thead>
        <tbody>
          <tr v-for="c in clients" :key="c.id">
            <td>{{ c.name }}</td>
            <td><code>{{ c.id }}</code></td>
            <td><a :href="c.url" target="_blank" rel="noreferrer">{{ c.url }}</a></td>
            <td class="row" style="gap:.5rem;">
              <button class="client_button_modify ghost" @click="edit(c.id)">Modify</button>
              <button class="client_button_remove ghost danger" @click="removeClient(c.id)">Remove</button>
            </td>
          </tr>
        </tbody>
      </table>
    </section>

    <section v-if="editing" class="card">
      <h3 style="margin-top:0;">Modify Client</h3>
      <ClientForm v-model="form" @save="save" @cancel="cancel">
        <template #danger>
          <button class="btn ghost danger" type="button" @click="removeClient(form.id)">Delete</button>
        </template>
      </ClientForm>
    </section>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { portal } from '../stores/portal'
import ClientForm from '../components/ClientForm.vue'

const clients = computed(()=> portal.clients.value)
const editing = ref(false)
const form = reactive({ id:'', name:'', url:'' })

function addNew(){
  const id = portal.addClient({ name: 'TNV', url: '/api/stream' })
  edit(id)
}

function edit(id){
  const c = portal.getClient(id)
  if (!c) return
  Object.assign(form, JSON.parse(JSON.stringify(c)))
  editing.value = true
}

function save(){
  portal.updateClient(form.id, { name: form.name, url: form.url })
  editing.value = false
}

function cancel(){ editing.value = false }

function removeClient(id){
  if (confirm('Remove this client?')) portal.removeClient(id)
}

function exportFile(){
  const blob = new Blob([portal.exportJson()], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'portal-settings.json'
  a.click()
  URL.revokeObjectURL(url)
}

function importFile(e){
  const input = e.target
  const file = input.files && input.files[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = () => {
    try {
      portal.importJson(String(reader.result))
      editing.value = false
    } catch(err){
      alert('Invalid JSON: ' + (err && err.message ? err.message : err))
    }
  }
  reader.readAsText(file)
}
</script>
