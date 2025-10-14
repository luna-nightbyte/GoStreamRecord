<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

import AnimatedButton from '../../components/AnimatedButton.vue';
import AdminSettings from './AdminSettings.vue';

// Reactive state variables
const statusText = ref('Not recording');
const statusColor = ref('#fab1a0');
const recorderProcesses = ref([]); 
const logTerminal = ref(null);

const activeProcesses = ref([]); 

let updateLogsInterval = null;
let updateStatusInterval = null;
let updateRecordersInterval = null; 
// Helper function for transient messages
const showResponse = (message, isError = false) => {
  const responseArea = document.getElementById("responseArea");
  const alertDiv = document.createElement("div");
  alertDiv.className = `alert ${isError ? "alert-danger" : "alert-info"}`;
  alertDiv.innerText = message;
  responseArea.appendChild(alertDiv);
  setTimeout(() => alertDiv.remove(), 5000);
};
const updateLogs = async () => {
  try {
    const res = await fetch("/api/logs");
    const data = await res.json();
    if (logTerminal.value) {
      logTerminal.value.innerText = data.join("\n");
      logTerminal.value.scrollTop = logTerminal.value.scrollHeight;
    }
  } catch (err) {
    console.error("Error fetching logs:", err);
  }
}; 

const updateStatus = async () => {
  try {
    const res = await fetch("/api/status");
    const data = await res.json();
    const status = data.status || "Not recording";
    statusText.value = status;
    statusColor.value = status === "Running" ? "#3dbb24af" : "#fab1a0";
    activeProcesses.value = Array.isArray(data.botStatus) ? data.botStatus : [];
  } catch (err) {
    console.error("Error fetching status:", err);
  }
};

const checkSourceAvailability = async (streamer, provider) => {
  try {
    const res = await fetch("/api/get-online-status", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ streamer, provider }),
    });
    const data = await res.json();
    return data.message === "true";
  } catch (err) {
    showResponse("Error starting process: " + err, true);
    return false;
  }
};

const addStreamer = async () => {
  const streamerInput = document.getElementById("streamerInput");
  const providerInput = document.getElementById("providerInput");
  const provider = providerInput.value.trim();
  const name = streamerInput.value.trim();
  if (!name) {
    showResponse("Streamer name cannot be empty.", true);
    return;
  }
  try {
    await fetch("/api/add-streamer?provider=" + provider, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ data: name }),
    });
    showResponse(`Added '${name}' to the list`);
    streamerInput.value = "";
    // Re-fetch recorders to update the UI
    updateRecorders();
  } catch (err) {
    console.error(err);
    showResponse("Error adding streamer.", true);
  }
};

const updateRecorders = async () => {
  try {
    const res = await fetch("/api/get-streamers");
    const streamers = await res.json();
    const activeMap = new Map(activeProcesses.value.map((p) => [p.website.username, p.is_recording]));

    const newRecorderProcesses = [];
    for (const streamer of streamers) {
      const isRecording = activeMap.get(streamer) || false;
      const isAvailable = await checkSourceAvailability(streamer, "chaturbate");
      newRecorderProcesses.push({
        name: streamer,
        isRecording,
        isAvailable,
      });
    }
    recorderProcesses.value = newRecorderProcesses;
  } catch (err) {
    console.error(err);
    showResponse("Error fetching streamers.", true);
  }
};


const uploadFile = async () => {
  const fileInput = document.getElementById("importFile");
  const file = fileInput.files[0];
  if (!file) {
    showResponse("Please select a file to upload.", true);
    return;
  }
  const formData = new FormData();
  formData.append("file", file);
  try {
    const res = await fetch("/api/import", {
      method: "POST",
      body: formData
    });
    console.log(res)
    // const data = await res.json();
    showResponse("Success!", false);
  } catch (err) {
    console.error(err);
    showResponse("Error uploading file.", true);
  }
};

const downloadFile = () => {
  window.location.href = "/api/export";
};

// Lifecycle hooks
onMounted(() => {
  updateRecorders();
  updateLogs();
  updateStatus();

  updateLogsInterval = setInterval(updateLogs, 3000);
  updateStatusInterval = setInterval(updateStatus, 5000);
  updateRecordersInterval = setInterval(updateRecorders, 5000);
});

onUnmounted(() => {
  clearInterval(updateLogsInterval);
  clearInterval(updateStatusInterval);
  clearInterval(updateRecordersInterval);
});AdminSettings

</script>


<template> 
        <!-- Settings Section --> 
            <div class="container card bg-white rounded-lg shadow-sm">
              <div class="Form card-body p-4">
                <div class="Header card-header p-4 border-b"><h1>Settings</h1></div>
                
                    <ul class="nav nav-tabs flex border-b mb-4" id="settingsTabs" role="tablist">
                        <li class="nav-item flex-grow" role="presentation">
                            <AnimatedButton class=" nav-link active w-full py-3 px-4 text-center"
                                id="users-management-tab" data-bs-toggle="tab" data-bs-target="#adminSection"
                                type="button" role="tab" text="Users">Users</AnimatedButton>
                        </li>
                        <li class="nav-item flex-grow" role="presentation">
                            <AnimatedButton class="button nav-link w-full py-3 px-4 text-center" id="recorder-management-tab"
                                data-bs-toggle="tab" data-bs-target="#recorderManagementSection" type="button"
                                role="tab" text="Recorder">Recorder</AnimatedButton>
                        </li>
                        <li class="nav-item flex-grow" role="presentation">
                            <AnimatedButton class="button nav-link w-full py-3 px-4 text-center" id="api-management-tab"
                                data-bs-toggle="tab" data-bs-target="#apiManagementSection" type="button" role="tab"  text="API">API
                            </AnimatedButton>
                        </li>
                    </ul>
                    <div class="tab-content">
                        <div class="tab-pane fade show active" id="adminSection" role="tabpanel">
                            <AdminSettings />
                        </div>
                        <div class="tab-pane fade" id="apiManagementSection" role="tabpanel">
                            <AdminSettings />
                        </div>

                        <div class="tab-pane fade" id="recorderManagementSection" role="tabpanel">
                            <div class="card-header p-4 border-b">Streamers</div>
                            <div class="card-body p-4">
                                <form @submit.prevent="addStreamer" class="mb-4">
                                    <label for="streamerInput" class="form-label font-semibold">Add a New Streamer:</label>
                                    <div class="flex gap-2">
                                        <input type="text" id="streamerInput"
                                            class="form-control flex-1 p-2 border rounded-md"
                                            placeholder="Streamer name" required />
                                        <input type="text" id="providerInput"
                                            class="form-control flex-1 p-2 border rounded-md" placeholder="Hosting site"
                                            required />
                                        <AnimatedButton type="submit"
                                            class="btn bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-full transition-colors duration-300">
                                            Add</AnimatedButton>
                                    </div>
                                </form>
                                <article class="mb-4">
                                    <div class="card-header font-semibold">Current Streamers</div>
                                    <ul class="list-group">
                                        <!-- You'll need to populate this list dynamically -->
                                    </ul>
                                </article>
                            </div>

                            <div class="card bg-white rounded-lg shadow-sm">
                                <div class="card-header p-4 border-b">Import/Export</div>
                                <div class="card-body p-4">
                                    <div class="mb-4">
                                        <label for="importFile" class="form-label font-semibold">Import File:</label>
                                        <input type="file" id="importFile"
                                            class="form-control mb-2 p-2 border rounded-md" />
                                        <AnimatedButton @click="uploadFile"
                                            class="btn bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded-full transition-colors duration-300">
                                            Upload</AnimatedButton>
                                    </div>
                                    <div>
                                        <label class="form-label font-semibold">Export File:</label><br />
                                        <AnimatedButton @click="downloadFile"
                                            class="btn bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded-full transition-colors duration-300">
                                            Download</AnimatedButton>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>  <div id="responseArea" class="fixed bottom-4 right-4 z-50 space-y-2"></div>

            </div>
    </template>

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
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: #00477640;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  border-radius: 15px;
}

header {
  text-align: center;
  margin-bottom: 30px;
  padding: 20px;
  background: linear-gradient(45deg, #adb6be, #004776);
  border-radius: 15px;
}

h1 {
  font-size: 48px;
  margin: 0;
  color: #06153e;
}
 

.Form {
  text-align: center;
  margin-bottom: 30px;
  padding: 20px;
  background: linear-gradient(45deg, #004776, #adb6be);
  border-radius: 15px;
}
.Header {
  text-align: center;
  margin-bottom: 30px;
  padding: 20px;
  background: linear-gradient(45deg, LivestreamTab);
  border-radius: 15px;
}

.form-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.form-option {
  display: flex;
  flex-direction: row;
  align-items: center;
}

.form-option input[type="text"],
.form-option input[type="radio"],
.form-option input[type="checkbox"] {
  margin-top: 5px;
  padding: 10px;
  border: none;
  border-radius: 5px;
  width: 100%;
}

.form-option label {

  font-size: 16px;
  color: #333;
}


.progress-container {
    width: 100%;
    background-color: #f3f3f3;
    border: 1px solid #ccc;
    border-radius: 5px;
    overflow: hidden;
  }
  
  .progress-bar {
    height: 30px;
    background-color: #4caf50;
    transition: width 0.5s;
  }
  
  .progress-text {
    margin-top: 10px;
    font-size: 16px;
  }
  </style>
  