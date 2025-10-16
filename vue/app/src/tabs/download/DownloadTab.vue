<template>
  <div class="videoD-container">
    <div class="card_rec">
      <div class="card-header">Download Video</div>
      <div class="card-body">
        <form @submit.prevent="submitForm" class="space-y-4">
          <fieldset class="form-option-group">
            <legend>Select a Source</legend>
            <div class="options-container">
              <label v-for="site in sites" :key="site.value" :class="{ 'selected': formData.option === site.value }">
                <input type="radio" class="select-download" v-model="formData.option" :value="site.value" />
                {{ site.name }}
              </label>
            </div>
          </fieldset>

          <!-- Bulk option -->
          <div class="checkbox-container" title="Enable to download multiple videos at once">
            <input type="checkbox" id="bulk" v-model="formData.bulk" />
            <label for="bulk">Bulk</label>
          </div>
          <!-- Search -->
          <div class="form-option"
            title="Optional for bulk downloading. Just search here instead of heading to the site.">
            <label for="search">Search:</label>
            <br>
            <input type="text" id="search" v-model="formData.search" placeholder="Enter search keywords" />
          </div>

          <!-- URL (if no search or bulk) -->
          <div v-if="!hasSearch && !hasBulk" class="form-option"
            title="Paste the URL of the video you wish to download">
            <label for="url">URL:</label>
            <br>
            <input type="text" id="url" v-model="formData.url" placeholder="Enter video URL" />
          </div>

          <!-- Save name -->
          <div v-if="!hasSearch && !hasBulk" class="form-option" title="Specify a name to save the video">
            <label for="name">Save name:</label>
            <br>
            <input type="text" id="name" v-model="formData.save" placeholder="Enter save name" />
          </div>

          <!-- Submit Button -->
          <div class="form-option">
            <AnimatedButton @click="submitForm" text="Download" hovertext="Download" />
          </div>

        </form>

        <!-- Progress -->
        <div v-if="downloading" class="progress-container mt-4">
          <div class="progress-bar" :style="{ width: progress + '%' }"></div>
        </div>
        <div class="progress-text mt-2">{{ progressText }}</div>
        <div class="progress-text">{{ queueText }}</div>

      </div>
    </div>
  </div>
</template>

<script>
import AnimatedButton from '../../components/AnimatedButton.vue';
import axios from 'axios';

import { notify } from "@/composables/useNotifications";

const { showResponse } = notify();
export default {
  components: { AnimatedButton },
  data() {
    return {
      formData: {
        option: 'Pornhub', // Default value
        bulk: false,
        search: '',
        url: '',
        save: 'Default_name'
      },
      sites: [
        { name: 'Pornhub', value: 'Pornhub' },
        { name: 'Xnxx', value: 'Xnxx' },
        { name: 'Xvideos', value: 'Xvideos' },
        //{ name: 'PornOne', value: 'Pornone' },
       // { name: 'Spankbang', value: 'Spankbang' },
      ],
      downloading: false,
      progress: 0,
      progressText: '',
      queueText: '',
      formData2: { running: false, total: 0, current: 0, progress: 0, progressText: '', queueText: '' },
    };
  },
  computed: {
    hasSearch() { return this.formData.search && this.formData.search.trim().length > 0; },
    hasBulk() { return this.formData.bulk; },
  },
  methods: {
    async fetchProgress() {
      try {
        const res = await axios.get('/api/progress');
        this.formData2 = res.data;
        this.progress = (this.formData2.current / this.formData2.total) * 100 || 0;
        this.progressText = this.formData2.progressText;
        this.downloading = this.formData2.running;
        this.queueText = this.formData2.queueText;

        if (this.downloading) setTimeout(this.fetchProgress, 250);
      } catch (err) { console.error('Error fetching progress:', err); }
    },
    async submitForm() {
      try {
        await axios.post('/api/download', this.formData);
        this.fetchProgress(); // Start polling for progress after submitting
        showResponse("Sucess! Starting video download")
      }
      catch (err) { 

        showResponse("Error submitting download:", err) }
    },
  },
  created() { this.fetchProgress(); }
};
</script>

<style>
/* No scoped styles: uses global theme */
.form-option {
  margin-bottom: 1rem;
}

.space-y-4>*+* {
  margin-top: 1rem;
}
</style>
