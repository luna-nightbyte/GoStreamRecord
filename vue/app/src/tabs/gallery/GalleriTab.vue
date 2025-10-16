<template>
  <div class="videoD-container">

        <button @click="sendCommand('repair', '')"
          class="buttonclass bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-6 rounded-full transition-colors duration-300">
          Run video codec repair </button>
      <div class="card-header p-4 border-b">Gallery</div>

    <section class="media-list card_rec" id="localFilesID">
      <div v-if="localMedia.length === 0">
        No videos available
      </div>
      <div v-else class="row">
        <figure v-for="url in localMedia" :key="url" class="column media-container" @click="openFullScreen(url.url)">
          <div class="notransform" v-if="url.error.length !== 0">
            <p>{{ url.error }}</p>
          </div>
          <div v-else >
          
          <img v-if="isImage(url.url)" :src="url.url" alt="Image" class="image" />
          <video v-if="isVideo(url.url)" controls class="video">
            <source :src="url.url" type="video/mp4" />
          </video>
          <figure class="media-caption ">{{ url.name }}</figure>
        </div>
        </figure>

      </div>
      <div class="pagination">

        <button @click="prevPage" :disabled="page === 1">Previous</button>
        <span>Page {{ page }}</span>
        <button @click="nextPage" :disabled="localMedia.length < pageSize">Next</button>
      </div>
    </section>
  </div>
</template>

<script>
import axios from 'axios';
import { notify } from "@/composables/useNotifications";

const { showResponse } = notify();
export default {
  data() {
    return {
      localMedia: [],
      page: 1,
      pageSize: 10,
      searchTerm: '',
      fullScreenMedia: null,
    };
  },
  methods: {
    async sendCommand(command, name = '') {
      try {
        await fetch("/api/control", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ command, name }),
        });
        showResponse("Success!");
      } catch (err) {
        showResponse(`Error: ${err}`, true);
      }
    },

    showResponse(message, isError = false) {
      const responseArea = document.getElementById("responseArea");
      if (!responseArea) return;

      while (responseArea.childNodes.length > 5) {
        responseArea.removeChild(responseArea.firstChild);
      }

      const alertDiv = document.createElement("div");
      alertDiv.className = `alert ${isError ? "alert-danger" : "alert-info"}`;
      alertDiv.innerText = message;
      responseArea.appendChild(alertDiv);
      setTimeout(() => alertDiv.remove(), 5000);
    },

    fetchVideos() {
      axios.get(`/api/videos?page=${this.page}`)
        .then(response => {
          this.localMedia = response.data;
        })
        .catch(error => {
          console.error("There was an error fetching the videos!", error);
        });
    },

    nextPage() {
      this.page++;
      this.fetchVideos();
    },

    prevPage() {
      if (this.page > 1) {
        this.page--;
        this.fetchVideos();
      }
    },

    isImage(url) {
      return /\.(jpg|jpeg|png|gif)$/i.test(url);
    },

    isVideo(url) {
      return /\.(mp4)$/i.test(url);
    },

    searchMedia() { 
    }
  },
  mounted() {
    this.fetchVideos();
  },
};
</script>


<style scoped>
:root {
  --primary-color: #004776;
  --secondary-color: #adb6be;
  --background-color: #0a0a0a;
  --text-color: #e0e0e0;
  --font-family: 'Roboto', sans-serif;
}

body {
  margin: 0;
  font-family: var(--font-family);
  background-color: var(--background-color);
  color: var(--text-color);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.media-list {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}

.media-list {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}

.row {
  display: flex;
  flex-wrap: wrap;
  padding: 0 4px;
  width: 100%;
}

.column {
  flex: 1 1 calc(25% - 1px);
  /* 4 columns */
  padding: 0 8px;
  margin-bottom: 16px;
  box-sizing: border-box;
}

@media (max-width: 1200px) {
  .column {
    flex: 1 1 calc(25% - 16px); 
  }
}

@media (max-width: 768px) {
  .column {
    flex: 1 1 calc(50% - 16px); 
  }
}

@media (max-width: 480px) {
  .column {
    flex: 1 1 100%; 
  }
}

.media-container {
  position: relative;
  overflow: hidden;
  border-radius: 15px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.5);
  background: linear-gradient(45deg, var(--primary-color), var(--secondary-color));
  transition: transform 0.4s ease, box-shadow 0.4s ease;
  cursor: pointer;
}

.media-container:hover {
  transform: scale(1.05);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.7);
}

.image,
.video {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.media-caption {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  padding: 10px;
  background-color: rgba(0, 0, 0, 0.7);
  color: #fff;
  text-align: center;
  font-size: 16px;
  opacity: 0;
  transition: opacity 0.75s ease;
}

.media-container:hover .media-caption {
  opacity: 1;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 20px;
}

.pagination button {
  background-color: var(--primary-color);
  color: #fff;
  border: none;
  padding: 10px 20px;
  margin: 0 10px;
  cursor: pointer;
  border-radius: 5px;
  transition: background-color 0.3s ease;
}

.pagination button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.pagination span {
  color: var(--text-color);
  font-size: 18px;
}
</style>
