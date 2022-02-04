<template>
  <div class="relative">
    <Toolbar v-bind:on="['/leaderboard', '/projects']"/>
    <router-view></router-view>
  </div>
</template>

<script setup>
import './assets/index.css';
import Toolbar from './components/Toolbar.vue';
</script>
