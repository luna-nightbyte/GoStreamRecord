<template>
    <header class="header">
      <h2  style="margin:0">Client Views</h2> 
      <RouterLink class="button_form ghost" to="/settingsView">Manage</RouterLink>
    </header>
                       
  <div class="item-container" style="gap:1rem;">
     
    <div  style="display:grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap:1rem;">
      <article class=" card" v-for="c in clients" :key="c.id" style="display:flex; flex-direction:column; gap:.5rem;">
        <header class="client_row" style="align-items:center; gap:.5rem;">
          <b >{{ c.name }}</b>
          <span class="pill">{{ c.id }}</span>
          <div class="spacer"></div>
        </header>
        <div class="client_url">
          {{ c.url }}
        </div>
        <div class="client_row" style="gap:.5rem; margin-top:.5rem;">
          <RouterLink class="client_button_open" :to="`/client/${c.id}`">Open</RouterLink>
          <a class="client_button_newTab ghost" :href="c.url" target="_blank" rel="noreferrer">New Tab</a>
        </div>
      </article>
  </div>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { portal } from '../stores/portal'

const clients = computed(()=> portal.clients.value)
</script>
<style>
.header {
  margin-left: 0%;
  padding: 20px;
  background: linear-gradient(45deg, #004776, #adb6be);
}
.content {
  display: flex;
  height: calc(100vh - 120px); 
}

.client_button_modify,
.client_button_remove,
.client_button_open,
.client_button_newTab,
.button_form {
  list-style-type:none; 
  margin: 5px;
  position: relative;
  display: inline-block;
  padding: 15px 30px;
  font-size: 18px;
  margin: 5px;
  font-weight: bold;
  color: rgba(255, 255, 255, 0.357);
  background: #4d6162;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  overflow: hidden;
  transition: transform 0.5s, box-shadow 0.4s;
  z-index: 1;
  box-shadow: 0 4px 25px rgb(44, 42, 51);
}

.client_button_newTab {
  list-style-type:none; 
  color: rgba(41, 60, 66, 0.615);
  background: linear-gradient(45deg, #2eab9d, #00f2ff);
}

.client_button_open {
  list-style-type:none; 
  color: rgba(41, 60, 66, 0.615);
  background: linear-gradient(-45deg, #2eab9d, #00f2ff);
}

.client_button_modify {
  list-style-type:none; 
  color: rgba(41, 60, 66, 0.615);
  background: linear-gradient(45deg, #60ab2e, #ccff00);
}

.client_button_remove {
  list-style-type:none; 
  color: rgba(137, 199, 218, 0.615);
  background: linear-gradient(-45deg, #2b0e36, rgb(114, 3, 34));
}


.client_url {
  list-style-type:none; 
  color: #253e46;
  font-weight: bold;
  font-size: 16px;
  background: linear-gradient(45deg, #2eab9d, #00f2ff);
  box-shadow: 0 0px 50px #0d69b5a2;
  border-radius: 5px;
  margin: 0 auto;
  padding: 10px;
  overflow: hidden;
  max-width: 100%;
}

.button_form:hover {
  color: rgb(255, 255, 255);
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.091);
  overflow: hidden;
  transition: transform 0.5s, box-shadow 0.4s;
  background: linear-gradient(45deg, #4c6b6d, #42344f);

} 
</style>