<template>
    <div>
      <div class="progress-container">
        <div class="progress-bar" :style="{ width: progress + '%' }"></div>
      </div>
      <div class="progress-text">{{ progressText }}</div>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  
  export default {
    data() {
      return {
        progress: 0,
        progressText: '',
        formData2: {
          running: '',
          total: '',
          progress: 0   ,
          progressText: ''
      }
      }
    },
    created() {
      this.fetchProgress();
    },
    methods: {  
      async fetchProgress() {
        try {
          const response = await axios.get('/api/progress');
          this.progress = response.data.progress;
          this.progressText = response.data.text;
  
          if (this.progress < 100) {
            setTimeout(this.fetchProgress, 1000); // Poll every second
          }
        } catch (error) {
          console.error('Error fetching progress:', error);
        }
      }
    }
  }
  </script>
  
  <style>
  
  </style>
  