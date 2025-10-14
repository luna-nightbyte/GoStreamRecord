<template>
  <div id="app">
    <header class="header_class"> 
        <h1 class="header-text">Stream CTRL</h1> 
    </header>
    <div class="main-layout-container">
      <nav class="tabs-nav"> 
          <!-- <div class="tab-btn"><router-link to="/logsTab">
              <TabButton text="Logs" hovertext="Logs" />
            </router-link></div>  -->
          <div class="tab-btn"><router-link to="/galleryTab">
              <TabButton text="Gallery" hovertext="View local files" />
            </router-link></div> 
          <div class="tab-btn"><router-link to="/downloadTab">
              <TabButton text="Download" hovertext="Download files" />
            </router-link></div> 
          <!-- <div class="tab-btn"><router-link to="/settingsTab">
              <TabButton text="Settings" hovertext="Edit settings" />
            </router-link></div>   -->
          <div class="tab-btn"><router-link to="/livestreamTab">
              <TabButton text="LiveStream" hovertext="Watch live stream" />
            </router-link></div>  
          <div class="tab-btn"><router-link to="/recorderTab">
              <TabButton text="Recorder" hovertext="Record live stream videos" />
            </router-link></div>  
      </nav>
       <router-view></router-view>
      
    </div>

    <header class="bottom-header">
      <div class="header-content">
        <p>Thank you for using our tool. Happy downloading!</p>
      </div>
    </header>
  </div>
</template>

<script>
import TabButton from './components/TabButton.vue';

export default {
  name: 'App',
  components: {
    TabButton,
  },
};
</script>
 