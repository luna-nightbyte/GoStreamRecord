<script setup>
import { ref, onMounted } from 'vue';

// Refs
const statusText = ref('Not recording');
const statusColor = ref('#fab1a0');
const recorderProcesses = ref([]);
const videoFiles = ref({});
const logTerminal = ref(null);
// const selectedVideos = ref(new Set()); // Stores video URLs
const activeProcesses = ref([]);

// Helper: show messages
const showResponse = (message, isError = false) => {
  const responseArea = document.getElementById("responseArea");
  if (!responseArea) return;
  
  // Keep the response area from getting too cluttered
  while (responseArea.childNodes.length > 5) {
      responseArea.removeChild(responseArea.firstChild);
  }

  const alertDiv = document.createElement("div");
  alertDiv.className = `alert ${isError ? "alert-danger" : "alert-info"}`;
  alertDiv.innerText = message;
  responseArea.appendChild(alertDiv);
  setTimeout(() => alertDiv.remove(), 5000);
};

// API commands
const sendCommand = async (command, name = '') => {
  try {
    await fetch("/api/control", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ command, name }) });
    showResponse("Success!");
    await refreshAllData(); // Automatically refresh data after a command
  } catch (err) { showResponse(`Error: ${err}`, true); }
};

const fetchLogs = async () => {
  try {
    const res = await fetch("/api/logs");
    const data = await res.json();
    if (logTerminal.value) {
      logTerminal.value.innerText = data.join("\n");
      logTerminal.value.scrollTop = logTerminal.value.scrollHeight;
    }
  } catch (err) { console.error(err); }
};

const fetchStatus = async () => {
  try {
    const res = await fetch("/api/status");
    const data = await res.json();
    statusText.value = data.status || "Not recording";
    statusColor.value = statusText.value === "Running" ? "#3dbb24af" : "#fab1a0";
    activeProcesses.value = Array.isArray(data.botStatus) ? data.botStatus : [];
  } catch (err) { console.error(err); }
};

// NEW: Manual refresh function to replace intervals
const refreshAllData = async () => {
  showResponse("Refreshing data...");
  try {
    // Fetch status first, as updateRecorders depends on its data (activeProcesses)
    await fetchStatus();
    // Update recorders and logs in parallel for efficiency
    await Promise.all([updateRecorders(), fetchLogs()]);
    showResponse("Data refreshed!");
  } catch (err) {
    showResponse(`Error refreshing data: ${err}`, true);
    console.error(err);
  }
};


// Streamer API
const addStreamer = async () => {
  // const providerName = document.getElementById("providerName");
  const streamerName = document.getElementById("streamerName");
  if (!streamerName.value.trim()) { // || !providerName.value.trim()) {
    showResponse("Streamer and provider cannot be empty.", true);
    return;
  }
  try {
    await fetch(
      `/api/add-streamer?provider=chaturbate`, { // ${providerName.value}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ data: streamerName.value.trim() })
    });
    showResponse(`Added '${streamerName.value}'`);
    streamerName.value = '';
    // providerName.value = '';
    updateRecorders(); // Refresh the list after adding
  } catch (err) { showResponse(err, true); }
};

// Recorders
const checkAvailability = async (streamer, provider) => {
  try {
    const res = await fetch("/api/get-online-status", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ streamer, provider })
    });
    const data = await res.json();
    return data.message === "true";
  } catch { return false; }
};

const updateRecorders = async () => {
  const section = document.getElementById("videoFilesSection");
  if (section) section.hidden = true;
  const nowsection = document.getElementById("activeRecordersSection");
  if (nowsection) nowsection.hidden = false;
  
  try {
    const res = await fetch("/api/get-streamers");
    const streamers = await res.json();
    const activeMap = new Map(activeProcesses.value.map(p => [ p.website.username, p.is_recording]));
    const updated = [];
    for (const s of streamers) {  
      const   isRecording = activeMap.get(s) || false; 
      const isAvailable = await checkAvailability(s,"");
      updated.push({ name: s, isRecording, isAvailable });
    }
    recorderProcesses.value = updated;
  } catch (err) { showResponse(err, true); }
};

// Videos
const populateVideos = async () => {
  const section = document.getElementById("activeRecordersSection");
  if(section) section.hidden = true;
  const nowsection = document.getElementById("videoFilesSection");
  if (nowsection) nowsection.hidden = false;
  
  try {
    const res = await fetch("/api/get-videos");
    const data = await res.json();
    const folders = {};
    data.forEach(video => {
      const urlParts = video.url.split("/");
      const folder = urlParts[urlParts.length - 2] || "Uncategorized";

      folders[folder] = folders[folder] || [];
      folders[folder].push(video);
    });
    videoFiles.value = folders;
  } catch (err) { showResponse(err, true); }
};

// const toggleFolder = (folderId) => {
//   const folder = document.getElementById(folderId);
//   if (!folder) return;
//   folder.style.display = folder.style.display === 'none' ? 'grid' : 'none';
// };

// const handleVideoSelection = (event, videoUrl) => {
//   if (event.target.checked) {
//     selectedVideos.value.add(videoUrl);
//   } else {
//     selectedVideos.value.delete(videoUrl);
//   }
// };

// const deleteSelectedVideos = async () => {
//   if (!selectedVideos.value.size) return showResponse("No videos selected!", true);
//   const videosToDelete = Array.from(selectedVideos.value);

//   if (!confirm(`Are you sure you want to delete ${videosToDelete.length} video(s)?`)) {
//     return;
//   }

//   try {
//     const res = await fetch("/api/delete-videos", {
//       method: "POST",
//       headers: { "Content-Type": "application/json" },
//       body: JSON.stringify({ videos: videosToDelete })
//     });
//     const result = await res.json();
//     if (result.data.success) {
//       showResponse("Videos deleted!");
//       populateVideos();
//       selectedVideos.value.clear();
//     } else {
//       showResponse("Failed to delete all videos.", true);
//     }
//   } catch { showResponse("Error deleting videos.", true); }
// };

// Lifecycle: onMounted now just calls the initial data load functions once.
onMounted(() => {
  refreshAllData();
  populateVideos();
  updateRecorders();
});
</script>


<template>
  <div class="recorder-container">
    <div class="card_rec mb-6 bg-white rounded-lg shadow-sm">
      <div class="card-header p-4 border-b">Control Panel</div>
      <div class="card-body p-4 flex flex-wrap justify-center gap-4">
        <button @click="sendCommand('start', '')"
          class="buttonclass bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-6 rounded-full transition-colors duration-300">
          Start all recordings</button>
        <button @click="sendCommand('stop', '')"
          class="buttonclass bg-gray-800 hover:bg-gray-900 text-white font-bold py-2 px-6 rounded-full transition-colors duration-300">
          Stop all recordings</button>
        <button @click="sendCommand('restart', '')"
          class="buttonclass bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-6 rounded-full transition-colors duration-300">
          Restart all recordings</button>
        <!-- NEW REFRESH BUTTON -->
        <button @click="refreshAllData"
            class="buttonclass bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-6 rounded-full transition-colors duration-300">
            Refresh Data
        </button>
      </div>
    </div>
    <div class="card_rec mb-4 bg-white rounded-lg shadow-sm">
      <div class="card-header p-4 border-b">Add a New Chaturbate Streamer</div>
      <div class="card-body p-4">
        <form @submit.prevent="addStreamer" class="space-y-4">
          <div class="form-option grid grid-cols-1 md:grid-cols-3 gap-4 items-center">
            <input type="text" id="streamerName" class="form-control col-span-1 md:col-span-1 p-2 border rounded-md"
              placeholder="Streamer name" required />
            <!-- <input type="text" id="providerName" class="form-control col-span-1 md:col-span-1 p-2 border rounded-md"
              placeholder="Hosting site" required /> -->
            <button type="submit"
              class="buttonclass bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-full transition-colors duration-300">
              Add</button>
          </div>
        </form>
      </div>
    </div>
    <!-- Monitor Card with Nested Tabs -->
    <div class="card_rec bg-white rounded-lg shadow-sm">
      <div class="card-header p-4 border-b">Details</div>
      <div class="card-body p-4">
        <div class="flex space-x-2 mb-4">
          <!-- <button class="buttonclass nav-link flex-1 py-3 px-4 text-center" id="video-tab" data-bs-toggle="tab"
            data-bs-target="#videoFilesSection" type="button" role="tab" @click="populateVideos">Video Files
          </button> -->
          <!-- <button class="buttonclass nav-link flex-1 py-3 px-4 text-center" id="recorderStatus-tab" data-bs-toggle="tab"
            data-bs-target="#activeRecordersSection" type="button" role="tab" @click="updateRecorders">Individual bots
          </button> -->
        </div>
        <div class="tab-content">
          <!-- Video Files Section -->
          <div class="tab-pane fade" id="videoFilesSection" role="tabpanel">
            <div id="videoFilesContainer">
              <!-- <button @click="deleteSelectedVideos"
                class="buttonclass bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded-full mb-4">
                Delete Selected Videos ({{ selectedVideos.size }})</button> -->

              <!-- <div v-for="(videos, folder) in videoFiles" :key="folder"
                class="card_rec mb-4 shadow-sm rounded-lg overflow-hidden">
                <div class="card-header p-3 bg-gray-100 flex justify-between items-center cursor-pointer"
                  @click="toggleFolder(`folder-${folder}`)">
                  <strong class="text-gray-700">{{ folder }}</strong>
                  <span class="badge bg-blue-500 text-white px-2 py-1 rounded-full"> ({{ videos.length
                    }} files)</span>
                </div>

                <div :id="`folder-${folder}`" class="card-body p-4" style="display: none;">
                  <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                    <div v-for="video in videos" :key="video.url"
                      class="video-item card p-2 border rounded-lg shadow-sm relative bg-white">
                      <div class="form-check absolute top-2 right-2">
                        <input class="form-check-input video-checkbox w-5 h-5" type="checkbox" :value="video.url"
                          :checked="selectedVideos.has(video.url)" @change="handleVideoSelection($event, video.url)">
                      </div>
                      <h6 class="mt-4 mb-2 text-sm font-semibold truncate" :title="video.name">{{
                        video.name }}</h6>
                      <video controls preload="metadata" width="100%" class="rounded-lg object-cover h-32">
                        <source :src="video.url" type="video/mp4">
                        Your browser does not support the video tag.
                      </video>
                      <a :href="video.url" target="_blank" class="text-xs text-blue-500 hover:underline mt-1 block">Open
                        File</a>
                    </div>
                  </div>
                </div>
              </div> -->
            </div>
          </div>

          <!-- Active Recorders Section -->
          <div class="tab-pane fade" id="activeRecordersSection" role="tabpanel">
            <div id="recorderProcessesContainer" class="space-y-4">
              <div v-for="process in recorderProcesses" :key="process.name"
                class="recorder-process p-4 bg-gray-50 rounded-lg shadow-sm flex flex-col md:flex-row justify-between items-center"
                :data-name="process.name">
                <div class="flex-grow flex items-center gap-4 mb-2 md:mb-0">
                  <span class="font-bold text-lg"
                    :class="{ 'text-green-500': process.isRecording, 'text-red-500': !process.isRecording }">
                    {{ process.isRecording ? 'Recording ' : 'Not Recording ' }}
                  </span>
                  <span class="text-gray-600">
                    {{ process.name }} - <span
                      :class="{ 'text-green-500': process.isAvailable, 'text-red-500': !process.isAvailable }">
                      {{ process.isAvailable ? 'Online' : 'Not Online' }}
                    </span>
                  </span>
                </div>
                <div class="flex gap-2">
                  <button @click="sendCommand('start', process.name)"
                    :disabled="process.isRecording || !process.isAvailable"
                    class="buttonclass bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded-full transition-colors duration-300 disabled:opacity-50">
                    Start</button>
                  <button @click="sendCommand('stop', process.name)" :disabled="!process.isRecording"
                    class="buttonclass bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded-full transition-colors duration-300 disabled:opacity-50">
                    Stop</button>
                  <button @click="sendCommand('restart', process.name)" :disabled="!process.isRecording"
                    class="buttonclass bg-yellow-500 hover:bg-yellow-600 text-white font-bold py-2 px-4 rounded-full transition-colors duration-300 disabled:opacity-50">
                    Restart</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div id="responseArea" class="fixed bottom-4 right-4 z-50 space-y-2"></div>
  </div>
</template>
