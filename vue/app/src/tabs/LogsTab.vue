<script setup>
import { ref, onMounted, onUnmounted } from 'vue';


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
});

</script>

<template>
  <div class="item-container"  >
    <!-- Header with Title and Status -->
    <header >
      <h1 class="text-3xl font-bold text-gray-800">GoStreamRecord WebUI</h1>
      <div id="statusIndicator"
           class="px-4 py-2 rounded-full text-white font-semibold transition-colors duration-300"
           :style="{ backgroundColor: statusColor }">
        Status: {{ statusText }}
      </div>
    </header> 
    <!-- Response area for transient messages -->
    <div id="responseArea" class="fixed bottom-4 right-4 z-50 space-y-2"></div>
  </div>
</template>

<style scoped>
.nav-pills .nav-link {
  color: #4b5563;
  border-radius: 9999px;
  text-align: left;
}
.nav-pills .nav-link.active {
  background-color: #3b82f6;
  color: white;
}
.tab-content > .tab-pane {
  display: none;
}
.tab-content > .show.active {
  display: block;
}
.nav-tabs .nav-link {
  color: #4b5563;
}
.nav-tabs .nav-link.active {
  background-color: #f3f4f6;
  border-color: #e5e7eb #e5e7eb #f3f4f6;
}

</style>
