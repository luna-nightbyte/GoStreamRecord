import { createRouter,createWebHashHistory } from 'vue-router';

// Components
import AnimatedButton from '../components/AnimatedButton.vue';
import TabButton from '../components/TabButton.vue';
import ProgressBar from '../components/ProgressBar.vue';

// Tabs
import LivestreamTab from '@/tabs/livestream/LivestreamTab.vue';  
import RecorderTab from '@/tabs/recorder/RecorderTab.vue';  
import GalleryTab from '../tabs/gallery/GalleriTab.vue';
import DownloadTab from '../tabs/download/DownloadTab.vue'; 
import AboutTab from '@/tabs/about/AboutTab.vue'; 
import SettingsTab from '@/tabs/settings/SettingsTab.vue'; 
// import SettingsTab from '@/tabs/SettingsTab.vue'; 

// Cards
import LogsTab from '@/tabs/logs/LogsTab.vue';  
import AdminSettings from '@/tabs/settings/AdminSettings.vue'; 
import GeneralSettings from '@/tabs/settings/GeneralSettings.vue'; 
import ApiSectionSection from '@/tabs/settings/sections/ApiSection.vue';  
import UserSettingsSection from '@/tabs/settings/sections/UserSection.vue'; 


import SettingsView from '@/views/SettingsView.vue'
import ClientViews from '@/views/ClientViews.vue'
import ClientDetail from '@/views/ClientDetail.vue'

const routes = [
  { path: '/', redirect: '/gallery' },
  { path: '/gallery', redirect: '/galleryTab' },
  { path: '/getVideo', redirect: '/downloadTab' },

  // Tabs
  { path: '/settingsTab', name: 'settings', component: SettingsTab },
  { path: '/galleryTab', component: GalleryTab },
  { path: '/downloadTab', component: DownloadTab }, 
  { path: '/livestreamTab', name: 'livestream', component: LivestreamTab }, 
  { path: '/recorderTab', name: 'recorder', component: RecorderTab }, 
  { path: '/logsTab', name: 'logs', component: LogsTab }, 

  { path: '/aboutTab', name: 'about', component: AboutTab }, 

  { path: '/animatedbtn', name: 'animatedbtn', component: AnimatedButton },
  { path: '/tabbtn', name: 'tabbtn', component: TabButton }, 
  { path: '/progress', name: 'progress', component: ProgressBar },
 
  { path: '/apiSettingsSection', name: 'apisection', component: ApiSectionSection },
  { path: '/userSettingsSection', name: 'usersection', component: UserSettingsSection },


    { path: '/clientView', redirect: '/clients' },
    { path: '/settingsView', component: SettingsView },
    { path: '/clients', component: ClientViews },
    { path: '/client/:id', component: ClientDetail, props: true },

// Old / to be discontinued
  { path: '/adminSettingsTab', name: 'settingsadmin', component: AdminSettings },
  { path: '/generalSettingsTab', name: 'settingsgeneraltab', component: GeneralSettings },

];

const router = createRouter({
  history:  createWebHashHistory(),
  routes

});

export default router;
