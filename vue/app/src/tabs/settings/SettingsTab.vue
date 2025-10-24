<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import AnimatedButton from '../../components/AnimatedButton.vue';

// ---------- Refs ---------- 
const activeProcesses = ref([]);
const recorderProcesses = ref([]);

const users = ref([]);
const selectedUser = ref('');

const apiKeys = ref([]);

const streamerName = ref('');
const importFileRef = ref(null);

let refreshTimer = null;

// ---------- Helpers ----------
// const badgeColor = (status) => (status === 'Running' ? '#3dbb24af' : '#fab1a0');

const safeJson = async (res) => {
  try { return await res.json(); } catch { return null; }
};


const fetchApiKeys = async () => {
  try {
    console.log("Fetching API keys...");
    const res = await fetch(`/api/keys`, {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    });
    const data = await res.json();
    apiKeys.value = data.map(key => ({ name: key.name, key: '*************' }));
  } catch (error) {
    console.error("Error fetching API key:", error);
  }
};

const generateNewApiKey = async () => {
  console.log("Generating new API key...");
  const apiKeyName = prompt("Enter a name for the new API key:");
  if (!apiKeyName) return;

  try {
    const response = await fetch("/api/generate-api-key?name=" + apiKeyName, { method: "POST" });
    const data = await response.json();
    if (data.status) {
      alert(`New API Key: ${data.key}\nThis will only be shown once!`);
      fetchApiKeys();
    } else {
      alert("Error generating API key.");
    }
  } catch (error) {
    console.error("Request failed", error);
    alert("An unexpected error occurred.");
  }
};

const deleteApiKey = async (name) => {
  console.log("Deleting API key:", name);
  if (confirm(`Are you sure you want to delete '${name}'?`)) {
    try {
      await fetch("/api/delete-api-key", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ data: { new: name } }),
      });
      fetchApiKeys();
    } catch (err) {
      console.error("Error deleting api key.", err);
    }
  }
};
// ---------- API: Control ----------
const sendCommand = async (command, name = '') => {
  console.log(`Sending command: ${command} for ${name || 'all'}`);
  try {
    await fetch('/api/control', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ command, name })
    });
    console.log('Success!');
    await refreshAllData();
  } catch (err) {
    console.log(`Error: ${err}`, true);
  }
};

// ---------- API: Users ----------
const fetchUsers = async () => {
  console.log('Fetching users...');
  try {
    const res = await fetch('/api/users');
    const data = await safeJson(res);
    users.value = Array.isArray(data) ? data : [];
  } catch (err) {
    console.error('fetchUsers', err);
  }
};

const addUser = async (evt) => {
  console.log('Adding user...');
  // Why: prevent reading stale DOM when browser autofill kicks in
  const form = evt?.target?.closest('form') || document;
  const name = form.querySelector('#addUserName')?.value?.trim();
  const pwd = form.querySelector('#addUserPassword')?.value || '';
  if (!name || !pwd) return console.log('Username and password required.', true);
  try {
    await fetch('/api/users', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, password: pwd })
    });
    console.log(`User '${name}' added`);
    form.querySelector('#addUserName').value = '';
    form.querySelector('#addUserPassword').value = '';
    await fetchUsers();
  } catch (err) {
    console.log(`Add user failed: ${err}`, true);
  }
};

const updateUser = async (evt) => {
  console.log('Updating user...');
  const form = evt?.target?.closest('form') || document;
  const name = form.querySelector('#updateUserName')?.value?.trim();
  const pwd = form.querySelector('#updateUserPassword')?.value || '';
  if (!name) return console.log('Select a user to update.', true);
  try {
    await fetch(`/api/users/${encodeURIComponent(name)}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password: pwd || undefined })
    });
    console.log(`User '${name}' updated`);
    form.querySelector('#updateUserPassword').value = '';
    await fetchUsers();
  } catch (err) {
    console.log(`Update failed: ${err}`, true);
  }
};

const deleteUser = async (name) => {
  console.log('Deleting user...');
  if (!name) return;
  try {
    await fetch(`/api/users/${encodeURIComponent(name)}`, { method: 'DELETE' });
    console.log(`User '${name}' deleted`);
    if (selectedUser.value === name) selectedUser.value = '';
    await fetchUsers();
  } catch (err) {
    console.log(`Delete failed: ${err}`, true);
  }
};



// ---------- API: Streamers ----------
const addStreamer = async () => {
  console.log('Adding streamer...');
  const name = streamerName.value.trim();
  if (!name) return;
  try {
    await fetch(`/api/add-streamer?provider=chaturbate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ data: name })
    });
    console.log(`Added '${name}'`);
    streamerName.value = '';
    await updateRecorders();
  } catch (err) {
    console.log(`Add streamer failed: ${err}`, true);
  }
};

const fetchStreamers = async () => {
  console.log('Fetching streamers...');
  try {
    const res = await fetch('/api/get-streamers');
    const data = await safeJson(res);
    return Array.isArray(data) ? data : [];
  } catch (err) {
    console.error('fetchStreamers', err);
    return [];
  }
};

const checkAvailability = async (streamer, provider = 'chaturbate') => {
  console.log(`Checking availability for ${streamer}...`);
  try {
    const res = await fetch('/api/get-online-status', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ streamer, provider })
    });
    const data = await safeJson(res);
    return data?.message === 'true';
  } catch {
    return false;
  }
};

const updateRecorders = async () => {
  console.log('Updating recorder statuses...');
  try {
    const streamers = await fetchStreamers();
    const activeMap = new Map(
      activeProcesses.value.map((p) => [p?.website?.username, p?.is_recording])
    );
    const updated = [];
    // Sequential availability checks to avoid hammering the API
    for (const s of streamers) {
      const isRecording = !!activeMap.get(s);
      const isAvailable = await checkAvailability(s, 'chaturbate');
      updated.push({ name: s, isRecording, isAvailable });
    }
    recorderProcesses.value = updated;
  } catch (err) {
    console.log(`Failed to update recorders: ${err}`, true);
  }
};

// ---------- Import / Export ----------
const uploadFile = async () => {
  console.log('Uploading import file...');
  const file = importFileRef.value?.files?.[0];
  if (!file) return console.log('Please select a file to upload.', true);
  const fd = new FormData();
  fd.append('file', file);
  try {
    const res = await fetch('/api/import', { method: 'POST', body: fd });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    console.log('Import success!');
    importFileRef.value.value = '';
    await refreshAllData();
  } catch (err) {
    console.log(`Error uploading file: ${err}`, true);
  }
};

const downloadFile = () => {
  console.log('Downloading export file...');
  // Why: simple and reliable to trigger file download
  window.location.href = '/api/export';
};

// ---------- Coordinator ----------
const refreshAllData = async () => {
  console.log('Refreshing data...');
  try {
    // await fetchStatus();
    await Promise.all([fetchUsers(), fetchApiKeys()]);
    await updateRecorders();
    console.log('Data refreshed!');
  } catch (err) {
    console.log(`Refresh failed: ${err}`, true);
  }
};

// ---------- Lifecycle ----------
onMounted(() => {
  fetchApiKeys();
  refreshAllData();
  refreshTimer = setInterval(async () => {
    try {
      // await fetchStatus();
      await Promise.all([updateRecorders()]);
    } catch (_) { /* noisy refresh not needed */ }
  }, 5000);
});

onUnmounted(() => {
  if (refreshTimer) clearInterval(refreshTimer);
});
</script>

<template> 

    <!-- Import / Export (compact) -->
    <div class="card_rec mb-6 bg-white rounded-lg shadow-sm">
      <div class="card-header p-4 border-b">Import & Export</div>
      <div class="card-body p-4 flex flex-wrap justify-center gap-4">
        <div class="mb-4">
          <input ref="importFileRef" id="importFile" type="file" class="form-control mb-2 p-2 border rounded-md" />
          <AnimatedButton @click="uploadFile" class="btn-neutral">Upload</AnimatedButton>
        </div>
        <div>
          <label class="form-label font-semibold">Export File:</label><br />
          <AnimatedButton @click="downloadFile" class="btn-neutral">Download</AnimatedButton>
        </div>
      </div>
    </div>
    <!-- Users & Groups (wider) -->
    <div class="card_rec mb-6 bg-white rounded-lg shadow-sm">
      <div class="card-header p-4 border-b font-semibold">Users & Groups</div>
      <div class="card-body p-4 flex flex-wrap justify-center gap-4">
        <!-- Add user -->
        <article class="p-2">
          <div class="card-header font-semibold p-2">Add user</div>
          <form @submit.prevent="addUser" class="Form">
            <div class="mb-3">
              <input type="text" id="addUserName" class="form-control w-full p-2 border rounded-md"
                placeholder="Username" required />
            </div>
            <div class="mb-3">
              <input type="password" id="addUserPassword" class="form-control w-full p-2 border rounded-md"
                placeholder="Password" required />
            </div>
            <AnimatedButton type="submit">Add User</AnimatedButton>
          </form>
        </article>

        <!-- Update user -->
        <article class="p-2">
          <div class="card-header font-semibold p-2">Update User</div>
          <form @submit.prevent="updateUser" class="Form">
            <div class="mb-3">
              <input type="text" id="updateUserName" class="form-control w-full p-2 border rounded-md"
                placeholder="Select a user from the list to modify." :value="selectedUser || ''" readonly />
            </div>
            <div class="mb-3">
              <input type="password" id="updateUserPassword" class="form-control w-full p-2 border rounded-md"
                placeholder="New Password" />
            </div>
            <AnimatedButton type="submit">Update User</AnimatedButton>
          </form>
        </article>
      </div>

      <!-- Current users -->
      <article class="p-2">
        <div class="card_rec mb-6 bg-white rounded-lg shadow-sm">
          <div class="card-header p-4 border-b font-semibold">Current users</div>
          <div class="card-body p-4 flex flex-wrap justify-center gap-4">
            <ul class="Form list-group space-y-2">
              <li v-for="user in users" :key="user.name"
                class="list-group-item d-flex justify-between items-center p-3 border rounded-md bg-white flex justify-between">
                <span>{{ user.name }}</span>
                <div class="space-x-2">
                  <button @click="selectedUser = user.name" class="chip">Edit</button>
                  <button @click="deleteUser(user.name)" class="chip-danger">Delete</button>
                </div>
              </li>
              <li v-if="!users.length" class="list-group-item p-3 text-center bg-white rounded-md">
                No users yet.
              </li>
            </ul>
          </div>
        </div>
      </article>
    </div>

    <!-- Streamers -->
    <div class="card_rec mb-6 bg-white rounded-lg shadow-sm">
      <div class="card-header p-4 border-b">streamers</div>
      <div class="card-body p-4 flex flex-wrap justify-center gap-4">
        <form @submit.prevent="addStreamer" class="space-y-3">
          <label class="form-label font-semibold">Add a New Chaturbate Streamer:</label>
          <div class="flex gap-2">
            <input v-model="streamerName" type="text" class="form-control flex-1 p-2 border rounded-md"
              placeholder="Streamer name" required />
            <AnimatedButton type="submit" class="btn-blue">Add</AnimatedButton>
          </div>
        </form>

        <ul class="list-group divide-y">
          <li v-for="proc in recorderProcesses" :key="proc.name" class="p-3 flex items-center justify-between">
            <div class="flex flex-col">
              <span class="font-medium">{{ proc.name }}</span>
              <small class="opacity-70">
                Recording: {{ proc.isRecording ? 'Yes' : 'No' }} · Online: {{ proc.isAvailable ? 'Yes' : 'No' }}
              </small>
            </div>
            <div class="flex items-center gap-2">
              <span class="status-tag" :class="proc.isAvailable ? 'ok' : 'off'">
                {{ proc.isAvailable ? 'Available' : 'Offline' }}
              </span>
              <button @click="sendCommand('start', proc.name)" :disabled="proc.isRecording || !proc.isAvailable"
                class="btn-primary disabled:opacity-50">Start</button>
              <button @click="sendCommand('stop', proc.name)" :disabled="!proc.isRecording"
                class="btn-danger disabled:opacity-50">Stop</button>
              <button @click="sendCommand('restart', proc.name)" :disabled="!proc.isRecording"
                class="btn-warn disabled:opacity-50">Restart</button>
            </div>
          </li>
          <li v-if="!recorderProcesses.length" class="p-3 text-center">
            No streamers yet.
          </li>
        </ul>
      </div>
    </div>


    <!-- API Keys -->
    <div class="card_rec mb-6 bg-white rounded-lg shadow-sm">
      <div class="card-header p-4 border-b">API keys</div>
      <div class="card-body p-4 flex flex-wrap justify-center gap-4">
        <form @submit.prevent="generateNewApiKey" class="Form mb-3">
          <label for="generateNewAPIButton" class="form-label font-semibold">Generate new API secret key</label>
          <div class="input-group">
            <button type="submit" id="generateNewAPIButton" class="btn-blue">Add</button>
          </div>
        </form>

        <ul class="list-group space-y-2">
          <li v-for="key in apiKeys" :key="key.name"
            class="list-group-item d-flex justify-between items-center p-3 border rounded-md bg-white flex justify-between">
            <span class="font-medium">{{ key.name }}</span>
            <div class="space-x-2 flex items-center">
              <span class="mono text-xs break-all">{{ key.key }}</span>
              <button @click="deleteApiKey(key.name)" class="btn-danger">Delete</button>
            </div>
          </li>
          <li v-if="!apiKeys.length" class="p-3 text-center bg-white rounded-md">
            No API keys yet.
          </li>
        </ul>
      </div> 
  </div>
</template>
