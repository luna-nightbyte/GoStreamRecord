<script setup>
import { ref, onMounted } from 'vue';
import { notify } from "@/composables/useNotifications";

const { showResponse } = notify();

// Refs
const statusText = ref('Not recording');
const statusColor = ref('#fab1a0');
const recorderProcesses = ref([]);
// const videoFiles = ref({});
const logTerminal = ref(null); 
const activeProcesses = ref([]);
 

// API commands
const sendCommand = async (command, name = '') => {
  try {
    await fetch("/api/control", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ command, name }) });
    showResponse("Success!");
    await refreshAllData();  
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
 
const refreshAllData = async () => {
  showResponse("Refreshing data...");
  try { 
    await fetchStatus(); 
    await Promise.all([updateRecorders(), fetchLogs()]);
    showResponse("Data refreshed!");
  } catch (err) {
    showResponse(`Error refreshing data: ${err}`, true);
    console.error(err);
  }
};


// Streamer API
const addStreamer = async () => { 
  const streamerName = document.getElementById("streamerName");
  if (!streamerName.value.trim()) {   
    showResponse("Streamer and provider cannot be empty.", true);
    return;
  }
  try {
    await fetch(
      `/api/add-streamer?provider=chaturbate`, { 
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ data: streamerName.value.trim() })
    });
    showResponse(`Added '${streamerName.value}'`);
    streamerName.value = ''; 
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

// // Videos
// const populateVideos = async () => {
//   const section = document.getElementById("activeRecordersSection");
//   if(section) section.hidden = true;
//   const nowsection = document.getElementById("videoFilesSection");
//   if (nowsection) nowsection.hidden = false;
  
//   try {
//     const res = await fetch("/api/get-videos");
//     const data = await res.json();
//     const folders = {};
//     data.forEach(video => {
//       const urlParts = video.url.split("/");
//       const folder = urlParts[urlParts.length - 2] || "Uncategorized";

//       folders[folder] = folders[folder] || [];
//       folders[folder].push(video);
//     });
//     videoFiles.value = folders;
//   } catch (err) { showResponse(err, true); }
// };
  
onMounted(() => {
  refreshAllData(); 
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
            <button type="submit"
              class="buttonclass bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-full transition-colors duration-300">
              Add</button>
          </div>
        </form>
      </div>
    </div> 
    <div class="card_rec bg-white rounded-lg shadow-sm">
      <div class="card-header p-4 border-b">Details</div>
      <div class="card-body p-4">
        <div class="flex space-x-2 mb-4"> 
        </div>
        <div class="tab-content">  
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
  </div>
</template>
