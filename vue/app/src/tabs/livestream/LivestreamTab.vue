<template>
  <div class="item-container flex flex-col lg:flex-row gap-6 p-4 max-w-6xl mx-auto">

    
    <div class="flex-1 card bg-white rounded-xl shadow-2xl overflow-hidden min-w-[320px]">
      <div class="p-4 bg-gray-50 border-b border-gray-200">
      <div class="card-header p-4 border-b">Live Stream Preview</div> 
      </div>
      <div class="card_rec p-2 aspect-video bg-gray-900 flex items-center justify-center">
        <iframe v-if="liveStreamSrc" :src="liveStreamSrc" width="100%" height="100%" frameborder="0" scrolling="no"
          class="rounded-lg shadow-xl w-full h-full min-h-[300px]"></iframe>
        <h3 v-else class="text-gray-400">Enter a streamer name and click Apply to load the live stream.</h3>
      </div>
    </div>

    
    <div class="lg:w-80 w-full card bg-white rounded-xl shadow-2xl p-6">
      <h3 class="card-header">Settings</h3> 

      <form class="card_rec" @submit.prevent="handleLivestreamForm">
        <div class="form-group mb-6 mt-4 ">
          <h4 for="streamerNameInput" >Streamer Name:</h4>  
          <input type="text" id="streamerNameInput" v-model.trim="streamerName"
            class="w-full"
            placeholder="e.g., JaneDoe123" required />
        </div>

        
        <div class="options-group mb-8">
          <h4>Stream View Options</h4> 
          <div class="options-container">
            <div v-for="option in streamOptions" :key="option.value" class="flex items-center">
              
              <input class="form-radio" type="radio"
                :id="`option-${option.value}`" :value="option.value" v-model="selectedOption">
              
              <label 
                :class="{ 'selected': selectedOption === option.value }"
                class="cursor-pointer select-none"
                :for="`option-${option.value}`">
                {{ option.label }}
              </label>
            </div>
          </div>
        </div>

        
        <button type="submit"
          class="buttonclass w-full">
          Apply Stream Settings
        </button>
      </form>
 

    </div>

    
    <Transition enter-active-class="transition ease-out duration-300" enter-from-class="opacity-0 translate-y-2"
      enter-to-class="opacity-100 translate-y-0" leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100 translate-y-0" leave-to-class="opacity-0 translate-y-2">
      <div v-if="message.text"
        :class="['fixed bottom-4 right-4 z-50 p-3 rounded-lg shadow-lg text-white font-medium', message.isError ? 'bg-red-500' : 'bg-green-500']">
        {{ message.text }}
      </div>
    </Transition>

  </div>
</template>

<script setup>
import { ref } from 'vue';

// --- Reactive State ---
const streamerName = ref('');
const selectedOption = ref('live'); // Default to 'live'
const liveStreamSrc = ref(null);
const message = ref({ text: null, isError: false }); // Reactive message state

// Define the available options for the radio buttons
const streamOptions = [
  { value: 'chat', label: 'Show Chat (Video + Chat)' },
  { value: 'live', label: 'Live Only (Video only)' },
  { value: 'interactive', label: 'Interactive (Different Embed)' },
];

// --- Utility Functions ---

/**
 * Shows a transient notification message using Vue's reactivity.
 * @param {string} text - The message text.
 * @param {boolean} isError - Whether it's an error message.
 */
const showResponse = (text, isError = false) => {
  // Clear any existing timeout
  if (window.messageTimeout) {
    clearTimeout(window.messageTimeout);
  }

  // Set the new message
  message.value = { text, isError };

  // Set a timeout to clear the message after 5 seconds
  window.messageTimeout = setTimeout(() => {
    message.value = { text: null, isError: false };
    delete window.messageTimeout;
  }, 5000);
};

// --- Main Form Handler ---

const handleLivestreamForm = () => {
  if (!streamerName.value) {
    showResponse("Please enter a streamer name.", true);
    return;
  }

  const streamer = streamerName.value;
  let newSrc = null;

  // Centralized URL construction based on the selected option
  switch (selectedOption.value) {
    case "chat":
      newSrc = `https://cbxyz.com/in/?tour=9oGW&campaign=Ln9WI&track=embed&room=${streamer}&disable_sound=1&embed_video_only=0&target=_parent&mobileRedirect=auto&`;
      break;
    case "live":
      newSrc = `https://cbxyz.com/in/?tour=9oGW&campaign=Ln9WI&track=embed&room=${streamer}&disable_sound=1&mobileRedirect=auto&embed_video_only=1`;
      break;
    case "interactive":
      newSrc = `https://cbxyz.com/in/?tour=Limj&campaign=Ln9WI&track=embed&signup_notice=1&b=${streamer}&disable_sound=1&mobileRedirect=never`;
      break;
    default:
      showResponse("Invalid stream option selected.", true);
      return;
  }

  liveStreamSrc.value = newSrc;
  showResponse(`Successfully set stream for: ${streamer} (${selectedOption.value} mode)`);
};

</script>

<style scoped>
/* Scoped styles ensure styles only apply to this component */
.item-container {
  min-height: 100vh;
}
/* This is important for the iframe to fill its responsive container */
.aspect-video {
  /* This is a simple Tailwind class, but defining it ensures flexibility */
  aspect-ratio: 16 / 9;
}
</style>
