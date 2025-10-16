// File: src/App.vue
<template>
  <div id="app">

    <header class="header_class">
      <h1 class="header-text">Stream CTRL
      <StatusIndicator/></h1>
    </header>
    

    <ResponseArea />
    <div class="main-layout-container">
      <nav class="tabs-nav">
        <div class="tab-btn">
          <router-link to="/galleryTab">
            <TabButton text="Gallery" hovertext="View local files" />
          </router-link>
        </div>
        <div class="tab-btn">
          <router-link to="/downloadTab">
            <TabButton text="Download" hovertext="Download files" />
          </router-link>
        </div>
        <div class="tab-btn">
          <router-link to="/livestreamTab">
            <TabButton text="LiveStream" hovertext="Watch live stream" />
          </router-link>
        </div>
        <div class="tab-btn">
          <router-link to="/recorderTab">
            <TabButton text="Recorder" hovertext="Record live stream videos" />
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
</script>
