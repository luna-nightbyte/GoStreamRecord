<script setup>
import {  onMounted } from 'vue';
 

onMounted(() => { 
});
</script>
 
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
  