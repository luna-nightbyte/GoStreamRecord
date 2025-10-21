<template>
  <div id="app">

    <header class="header_class">
      <h1 class="header-text">Stream CTRL
        <StatusIndicator/></h1>
    </header>


    <ResponseArea />
    <div class="main-layout-container">
      <nav class="tabs-nav"> 
        <div class="tab-btn" v-for="tab in tabs" :key="tab.id">
          <router-link :to="tab.routePath"> 
            <TabButton :text="tab.displayText" :hovertext="tab.description" />
          </router-link>
        </div> 
      </nav>

      <router-view />
    </div>

    <header class="bottom-header">
      <div class="header-content">
        <p>Thank you for using our tool. Happy downloading!</p>
      </div>
    </header>
  </div>
</template>

<script setup>
import StatusIndicator from './components/StatusIndicator.vue';
import ResponseArea from './components/ResponseArea.vue';
import TabButton from './components/TabButton.vue';

import { onMounted, ref } from 'vue';
 
const tabs = ref([]);
 
const getDisplayText = (name) => {
    const cleanName = name.replace('_tab', '');
    return cleanName.charAt(0).toUpperCase() + cleanName.slice(1);
};
 
const getRoutePath = (name) => { 
    const baseName = name.replace(/_(\w)/g, (match, p1) => p1.toUpperCase());
    return '/' + baseName.charAt(0).toUpperCase() + baseName.slice(1);
};


const updateStatus = async () => {
  try {
    const res = await fetch("/api/user_info");
    const resp = await res.json();
    const availableTabsObject =resp.tabs 
    console.log(availableTabsObject);
 
    const processedTabs = Object.values(availableTabsObject).map(tab => ({
        ...tab,
        displayText: getDisplayText(tab.name),
        routePath: getRoutePath(tab.name)
    }));
 
    tabs.value = processedTabs;
    console.log(processedTabs)
  } catch (err) {
    console.error("Error fetching status:", err); 
    tabs.value = [];
  }
}

onMounted(() => {
    updateStatus();
})
</script>
 
