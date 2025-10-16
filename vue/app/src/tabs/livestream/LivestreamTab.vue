<script setup>
import { ref, computed } from 'vue';
import { notify } from "@/composables/useNotifications";
const { showResponse } = notify();
 
const EMBED_HOST = 'https://cbxyz.com/in/';
const TOUR       = 'SHBY';    
const CAMPAIGN   = 'rQ9kN';  

const streamerName   = ref('');
const selectedOption = ref('live');
const liveStreamSrc  = ref(null);
const iframeKey      = ref(0);
const message        = ref({ text: null, isError: false });

const isLocalhost = computed(() => {
  const h = window.location.hostname;
  return h === 'localhost' || h === '127.0.0.1' || h.endsWith('.local');
});

const currentHostForEmbed = computed(() => {
  if (typeof window === 'undefined' || !window.location) return '';
  return window.location.host;  
});

const streamOptions = [
  { value: 'chat',        label: 'Show Chat (Video + Chat)' },
  { value: 'live',        label: 'Live Only (Video only)' },
  { value: 'interactive', label: 'Interactive (Different Embed)' },
];

function buildEmbedUrl(mode, roomRaw) {
  const room = roomRaw.trim();
  if (!room) return null;

  const u = new URL(EMBED_HOST);
  const p = u.searchParams;

  p.set('tour', TOUR);
  p.set('campaign', CAMPAIGN);
  p.set('track', 'embed');
  p.set('disable_sound', '1');
  p.set('embed_domain', currentHostForEmbed.value);  
  p.set('room', room);
 
  if (mode === 'live') p.set('embed_video_only', '1');
  else p.set('embed_video_only', '0');
 
  return u.toString();
}

const handleLivestreamForm = () => {
  if (!streamerName.value) {
    showResponse("Please enter a streamer name.", true);
    return;
  }
  const src = buildEmbedUrl(selectedOption.value, streamerName.value);
  liveStreamSrc.value = src;
  iframeKey.value++; 
  showResponse(`Loaded ${streamerName.value} (${selectedOption.value}).`);
};
</script>

<template>
    <meta name="referrer" content="origin-when-cross-origin" />
  <div class="item-container flex flex-col lg:flex-row gap-6 p-4 max-w-6xl mx-auto">
    <div class="flex-1 card bg-white rounded-xl shadow-2xl overflow-hidden min-w-[320px]">
      <div class="p-4 bg-gray-50 border-b border-gray-200">
        <div class="card-header p-4 border-b">Live Stream Preview</div>
      </div>

      <div class="card_rec p-2 aspect-video bg-gray-900 flex items-center justify-center">
        <iframe
          v-if="liveStreamSrc"
          :key="iframeKey"
          :src="liveStreamSrc"
          class="rounded-lg shadow-xl w-full h-full min-h-[300px]"
          frameborder="0"
          scrolling="no"
          allow="autoplay; fullscreen; picture-in-picture; encrypted-media"
          allowfullscreen
          referrerpolicy="origin-when-cross-origin"  
          title="Live stream embed"
        ></iframe>

        <div v-else class="flex flex-col items-center gap-3 text-center px-6">
          <h3 class="text-gray-200">Enter a streamer name and click “Apply”.</h3>
          <p v-if="isLocalhost" class="text-amber-300 text-sm">
            You’re on localhost. Some providers block localhost origins.
            Test via a public tunnel (Cloudflare: <code>cloudflared tunnel --url http://localhost:5173</code>)
            or deploy to a real HTTPS domain.
          </p>
        </div>
      </div>
    </div>

    <div class="lg:w-80 w-full card bg-white rounded-xl shadow-2xl p-6">
      <h3 class="card-header">Settings</h3>
      <form class="card_rec" @submit.prevent="handleLivestreamForm">
        <div class="form-group mb-6 mt-4">
          <h4 for="streamerNameInput">Streamer Name:</h4>
          <input id="streamerNameInput" v-model.trim="streamerName" class="w-full" placeholder="e.g., JaneDoe123" required />
        </div>

        <div class="options-group mb-8">
          <h4>Stream View Options</h4>
          <div class="options-container">
            <div v-for="option in streamOptions" :key="option.value" class="flex items-center gap-2">
              <input class="form-radio" type="radio" :id="`option-${option.value}`" :value="option.value" v-model="selectedOption" />
              <label :for="`option-${option.value}`" class="cursor-pointer select-none" :class="{ 'selected': selectedOption === option.value }">
                {{ option.label }}
              </label>
            </div>
          </div>
        </div>

        <button type="submit" class="buttonclass w-full">Apply Stream Settings</button>

 <div v-if="isLocalhost" class="mt-4 text-xs text-gray-600">
   Tip: When testing locally, we add <code>&embed_domain={{ currentHostForEmbed }}</code> automatically.
   Open your app via a tunnel URL (e.g., trycloudflare.com or ngrok) so the provider sees a real origin.
 </div>
      </form>
    </div>

    <Transition
      enter-active-class="transition ease-out duration-300"
      enter-from-class="opacity-0 translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-2">
      <div v-if="message.text" :class="['fixed bottom-4 right-4 z-50 p-3 rounded-lg shadow-lg text-white font-medium', message.isError ? 'bg-red-500' : 'bg-green-500']">
        {{ message.text }}
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.item-container { min-height: 100vh; }
.aspect-video { aspect-ratio: 16 / 9; }
</style>
