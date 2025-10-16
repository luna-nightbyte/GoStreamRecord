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
  
  // Note: lastStatus is better as a data property, not a static option.
  // Move 'lastStatus' into data() if you intend for it to be reactive 
  // or specific to this component instance. For now, we'll keep it as a simple 
  // property but know this is less "Vue standard."
  lastStatus: '', 
  
  data() {
    return {
      // Initialize with a representative structure for statusResp
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
    
    // 💡 Setup a polling interval to fetch the status regularly
    this.statusInterval = setInterval(this.fetchProgress, 5000); // e.g., every 5 seconds
  },
   
  // 💡 Don't forget to clear intervals when the component is destroyed!
  beforeUnmount() {
    clearInterval(this.statusInterval);
  },
   
  computed: {
    statusClass() {
       // Class names based on the current statusText
       const text = this.statusText.toLowerCase();
       return { 
        // Red/Danger status
        'status--down': text.includes('fixing') || text.includes('offline') || text.includes('error'),  
        // Yellow/Warning status
        'status--warn': text.includes('recording') || text.includes('downloading'), 
        // Green/Success status
        'status--ok': text.includes('online') || text.includes('idle'),
      };
    }
  },
  
  watch: { 
    statusResp: {
      // Handler runs whenever statusResp changes (e.g., after fetchProgress updates it)
      handler(newResp) {
        // 1. Update the display text and set lastStatus
        this.updateStatusText(newResp.status);
        
        // 2. The component will automatically update the UI based on statusText and statusClass
        
        // 3. REMOVED THE RELOAD: No need to reload, Vue handles the update!
        /* if (this.lastStatus !== this.statusText) {
          window.location.reload(); 
        }
        */
      },
      deep: true, // Needed because statusResp is an object
    },
  },
  
  methods: { 
    async fetchProgress() {
      try {
        const response = await axios.get('/api/status');
        
        // ⚠️ Before updating statusResp, save the current text for the watcher logic
        this.lastStatus = this.statusText;
        
        // Updating statusResp triggers the 'watch' handler
        this.statusResp = response.data;
      } catch (error) {
        console.error('Error fetching progress:', error); 
        this.statusText = 'API Error';
        this.lastStatus = this.statusText;
      }
    },
     
    updateStatusText(status) {
      // Determine the most critical status and set the text accordingly
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