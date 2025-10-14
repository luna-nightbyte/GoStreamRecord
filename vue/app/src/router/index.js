import { createRouter,createWebHashHistory } from 'vue-router';

import GalleryTab from '../tabs/gallery/GalleriTab.vue';
import DownloadTab from '../tabs/download/DownloadTab.vue'; 

import AnimatedButton from '../components/AnimatedButton.vue';
import TabButton from '../components/TabButton.vue';
import ProgressBar from '../components/ProgressBar.vue';

import GoStreamRecord from '@/tabs/LogsTab.vue';  
import LivestreamTab from '@/tabs/livestream/LivestreamTab.vue';  
import RecorderTab from '@/tabs/recorder/RecorderTab.vue';  

import SettingsTab from '@/tabs/SettingsTab.vue'; 
import AdminSettings from '@/tabs/settings/AdminSettings.vue'; 
import GeneralSettings from '@/tabs/settings/GeneralSettings.vue'; 
import ApiSectionSection from '@/tabs/settings/sections/ApiSection.vue'; 
import StreamerSectionSection from '@/tabs/settings/sections/StreamersSection.vue'; 
import UserSettingsSection from '@/tabs/settings/sections/UserSection.vue'; 


import SettingsView from '@/views/SettingsView.vue'
import ClientViews from '@/views/ClientViews.vue'
import ClientDetail from '@/views/ClientDetail.vue'

const routes = [
  { path: '/', redirect: '/gallery' },
  { path: '/gallery', redirect: '/galleryTab' },
  { path: '/getVideo', redirect: '/downloadTab' },
  { path: '/galleryTab', component: GalleryTab },
  { path: '/downloadTab', component: DownloadTab }, 
  { path: '/btn', name: 'btn', component: AnimatedButton },
  { path: '/tabbtn', name: 'tabbtn', component: TabButton }, 
  { path: '/progress', name: 'progress', component: ProgressBar },
  { path: '/livestreamTab', name: 'livestream', component: LivestreamTab }, 
  { path: '/recorderTab', name: 'recorder', component: RecorderTab }, 
  { path: '/logsTab', name: 'GoStreamRecord', component: GoStreamRecord }, 

  { path: '/settingsTab', name: 'settingsadmin', component: SettingsTab },
  { path: '/adminSettings', name: 'settingsadmin', component: AdminSettings },
  { path: '/generalSettings', name: 'settingsgeneraltab', component: GeneralSettings },
  { path: '/streamersSettingsSection', name: 'streamersection', component: StreamerSectionSection },
  { path: '/apiSettingsSection', name: 'apisection', component: ApiSectionSection },
  { path: '/userSettingsSection', name: 'usersection', component: UserSettingsSection },


    { path: '/clientView', redirect: '/clients' },
    { path: '/settingsView', component: SettingsView },
    { path: '/clients', component: ClientViews },
    { path: '/client/:id', component: ClientDetail, props: true },
];

const router = createRouter({
  history:  createWebHashHistory(),
  routes

});

export default router;
