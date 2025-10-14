<script setup>
import { ref, onMounted } from 'vue';

const apiKeys = ref([]); 
 
const fetchApiKeys = async () => {
  try {
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

onMounted(() => { 
  fetchApiKeys();
});
</script>

<template> 
  <div class="Form p-4 bg-gray-50 rounded-lg shadow-sm mt-4">
    <div class="h2 card-header font-semibold p-2" id="apiHeader">API keys</div>
    <div class="card-body p-2">
      <form @submit.prevent="generateNewApiKey" class="Form mb-3">
        <label for="generateNewAPIButton" class="form-label font-semibold">Generate new API secret key</label>
        <div class="input-group">
          <button type="submit" class="btn bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-full transition-colors duration-300" id="generateNewAPIButton">Add</button>
        </div>
      </form>
      <article class="mb-4">
        <div class="card-header font-semibold">API Keys</div>
        <ul class="list-group space-y-2">
          <li v-for="key in apiKeys" :key="key.name" class="list-group-item d-flex justify-content-between items-center p-3 border rounded-md bg-white">
            <span>{{ key.name }}</span>
            <div class="space-x-2">
              <span>{{ key.key }}</span>
              <button @click="deleteApiKey(key.name)" class="btn bg-red-500 hover:bg-red-600 text-white font-bold py-1 px-3 rounded-full transition-colors duration-300">Delete</button>
            </div>
          </li>
        </ul>
      </article>
    </div>
  </div>
</template>

<style scoped>
:root {
  --primary-color: #ff007f;
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

.header {
  text-align: center;
  margin-bottom: 30px;
  padding: 20px;
  background: linear-gradient(45deg, #adb6be, #004776);
  border-radius: 15px;
}

.h1 {
  font-size: 48px;
  margin: 0;
  color: #06153e;
}
 
.h2 {
  font-size: 24px;
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
  background: linear-gradient(45deg, #adb6be, #004776);
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
  