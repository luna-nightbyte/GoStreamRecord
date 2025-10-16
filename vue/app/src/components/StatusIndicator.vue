<template>
  <div class="status" :class="statusClass" aria-live="polite" :aria-label="`Status: ${statusText}`">
    <span class="status__dot" aria-hidden="true"></span>
    <span class="status__text">{{ statusText }}</span>
  </div>
</template>

<script>
import axios from 'axios'; 

export default {
  
  name: 'ProcessStatusIndicator',
   
  lastStatus: '', 
  
  data() {
    return { 
      statusResp: {
        status: {
          is_online: false,
          is_recording: false,
          is_downloading: false,
          is_fixing_codec: false,
        }
      },
      statusText: 'Loading...',  
    };
  },
   
  created() {
    this.fetchProgress();
     
    this.statusInterval = setInterval(this.fetchProgress, 5000);  
  },
    
  beforeUnmount() {
    clearInterval(this.statusInterval);
  },
   
  computed: {
    statusClass() { 
       const text = this.statusText.toLowerCase();
       return {  
        'status--down': text.includes('fixing') || text.includes('offline') || text.includes('error'),   
        'status--warn': text.includes('recording') || text.includes('downloading'),  
        'status--ok': text.includes('online') || text.includes('idle'),
      };
    }
  },
  
  watch: { 
    statusResp: { 
      handler(newResp) { 
        this.updateStatusText(newResp.status);
         
      },
      deep: true, 
    },
  },
  
  methods: { 
    async fetchProgress() {
      try {
        const response = await axios.get('/api/status');
         
        this.lastStatus = this.statusText;
         
        this.statusResp = response.data;
      } catch (error) {
        console.error('Error fetching progress:', error); 
        this.statusText = 'API Error';
        this.lastStatus = this.statusText;
      }
    },
     
    updateStatusText(status) { 
      if (status.is_fixing_codec) {
        this.statusText = 'Fixing Codec';
      } 
      else if (status.is_recording) {
        this.statusText = 'Recording...';
      } 
      else if (status.is_downloading) {
        this.statusText = 'Downloading...';
      } 
      else if (status.is_online) {
        this.statusText = 'Online & Idle';
      } 
      else {
        this.statusText = 'Offline';
      }
       
    }
  }
}
</script> 