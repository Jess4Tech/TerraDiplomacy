import {createWebHistory, createRouter} from 'vue-router';

import NotFound from '../components/NotFound.vue';
import Login from '../components/Login.vue';
import Projects from '../components/Projects.vue';
import Faction from '../components/Faction.vue';
import Leaderboard from '../components/Leaderboard.vue';

const routes = [
  {
    path: '/',
    name: 'LoginPage',
    component: Login,
  },
  {
    path: '/projects',
    name: 'ProjectList',
    component: Projects,
  },
  {
    path: '/faction',
    name: 'FactionViewer',
    component: Faction,
  },
  {
    path: '/leaderboard',
    name: 'Leaderboard',
    component: Leaderboard,
  },
  {
    path: '/:catchAll(.*)',
    component: NotFound,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
