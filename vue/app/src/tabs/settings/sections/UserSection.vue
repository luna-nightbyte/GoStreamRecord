<script setup>
import { ref, onMounted } from 'vue';

const users = ref([]); 
const selectedUser = ref(null);

const fetchUsers = async () => {
  try {
    const res = await fetch("/api/get-users");
    users.value = await res.json();
  } catch (err) {
    console.error("Error fetching users:", err);
  }
};

const addUser = async () => {
  const username = document.getElementById("addUserName").value.trim();
  const password = document.getElementById("addUserPassword").value;
  if (!username || !password) {
    alert("Username and password are required.");
    return;
  }
  try {
    await fetch("/api/add-user", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });
    document.getElementById("addUserName").value = "";
    document.getElementById("addUserPassword").value = "";
    fetchUsers();
  } catch (err) {
    console.error("Error adding user:", err);
  }
};

const updateUser = async () => {
  if (!selectedUser.value) {
    alert("No user selected.");
    return;
  }
  const newUsername = document.getElementById("updateUserName").value.trim();
  const newPassword = document.getElementById("updateUserPassword").value;
  if (!newUsername && !newPassword) {
    alert("Please provide new username and/or new password.");
    return;
  }
  try {
    await fetch("/api/update-user", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ oldUsername: selectedUser.value, newUsername, newPassword }),
    });
    selectedUser.value = null;
    document.getElementById("updateUserName").value = "";
    document.getElementById("updateUserPassword").value = "";
    fetchUsers();
  } catch (err) {
    console.error("Error updating user:", err);
  }
};

const deleteUser = async (name) => {
  if (confirm(`Are you sure you want to delete '${name}'?`)) {
    try {
      await fetch("/api/delete-user", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name }),
      });
      fetchUsers();
    } catch (err) {
      console.error("Error deleting user:", err);
    }
  }
};

onMounted(() => {
  fetchUsers(); 
});
</script>

<template>
  <div class="Form p-4 bg-gray-50 rounded-lg shadow-sm">
    <div class="card-body">
      <div class="h2 card-header font-semibold p-2">Add user</div>
      <article class="mb-4 p-2">
        <form @submit.prevent="addUser" class="Form">
          <div class="mb-3">
            <input type="text" id="addUserName" class="form-control w-full p-2 border rounded-md" placeholder="Username" required />
          </div>
          <div class="mb-3">
            <input type="password" id="addUserPassword" class="form-control w-full p-2 border rounded-md" placeholder="Password" required />
          </div>
          <button type="submit" class="btn bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-full transition-colors duration-300">Add User</button>
        </form>
      </article>
      <article class="mb-4 p-2">
        <div class="h2 card-header font-semibold">Update User</div>
        <form @submit.prevent="updateUser" class="Form">
          <div class="mb-3">
            <input type="text" id="updateUserName" class="form-control w-full p-2 border rounded-md" placeholder="Select a user from the list to modify." :value="selectedUser || ''" readonly/>
          </div>
          <div class="mb-3">
            <input type="password" id="updateUserPassword" class="form-control w-full p-2 border rounded-md" placeholder="New Password" />
          </div>
          <button type="submit" class="btn bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-full transition-colors duration-300">Update User</button>
        </form>
      </article>
      <article class="p-2">
        <div class="h2 card-header font-semibold">Current users</div>
        <ul class="Form list-group space-y-2">
          <li v-for="user in users" :key="user.name" class="list-group-item d-flex justify-content-between items-center p-3 border rounded-md bg-white">
            <span>{{ user.name }}</span>
            <div class="space-x-2">
              <button @click="selectedUser = user.name" class="btn bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-1 px-3 rounded-full transition-colors duration-300">Edit</button>
              <button @click="deleteUser(user.name)" class="btn bg-red-500 hover:bg-red-600 text-white font-bold py-1 px-3 rounded-full transition-colors duration-300">Delete</button>
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
  