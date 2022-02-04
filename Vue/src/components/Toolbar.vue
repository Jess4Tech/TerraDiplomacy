<template>
  <div
    v-if="props.on.includes(route.path)"
    class="flex flex-row-reverse absolute p-0 m-0 top-0 left-0 toggle-box"
    :class="{ 'slide-in': data.open, 'slide-out': data.close }"
  >
    <svg
      class="
        cursor-pointer
        rounded-lg
        fill-current
        sm:text-black
        md:text-white
        toggle-button
      "
      @click="toggle"
      alt="Open Toolbar"
      height="32px"
      id="Layer_1"
      style="enable-background: new 0 0 32 32"
      version="1.1"
      viewBox="0 0 32 32"
      width="32px"
      xml:space="preserve"
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
    >
      <path
        d="M4,10h24c1.104,0,2-0.896,2-2s-0.896-2-2-2H4C2.896,
        6,2,6.896,2,8S2.896,10,4,10z M28,14H4c-1.104,0-2,
        0.896-2,2  s0.896,2,2,2h24c1.104,0,2-0.896,2-2S29.104,
        14,28,14z M28,22H4c-1.104,0-2,0.896-2,2s0.896,2,2,
        2h24c1.104,0,2-0.896,2-2  S29.104,22,28,22z"
      />
    </svg>
    <div
      class="
        flex flex-col
        bg-white
        rounded-br-md rounded-tr-md
      "
    >
      <router-link
        class="border-2 shadow-md border-grey-700 rounded-md p-2"
        :to="{ name: 'ProjectList' }"
        >Projects</router-link
      >
      <router-link
        class="border-2 shadow-md border-grey-700 rounded-md p-2"
        :to="{ name: 'FactionViewer' }"
        >Faction</router-link
      >
      <router-link
        class="border-2 shadow-md border-grey-700 rounded-md p-2"
        :to="{ name: 'Leaderboard' }"
        >Leaderboard</router-link
      >
    </div>
  </div>
</template>

<script setup>
import {reactive} from 'vue';
import {useRoute} from 'vue-router';

// eslint-disable-next-line no-unused-vars
const route = useRoute();

// eslint-disable-next-line no-unused-vars
const props = defineProps({
  on: Array,
});

const data = reactive({
  open: false,
  close: true,
});

const toggle = function() {
  const temp = data.open;
  data.open = data.close;
  data.close = temp;
};
</script>

<style scoped>
.toggle-button {
  width: 10vw;
  height: 10vh;
}

.toggle-box {
  position: absolute;
  max-width: 35vw;
  max-height: 30vh;
  transform: translateX(50%);
  -webkit-transform: translateX(0%);
}

.slide-in {
  animation: slide-in 0.5s forwards;
  -webkit-animation: slide-in 0.5s forwards;
}

@keyframes slide-in {
  100% {
    transform: translateX(0%);
  }
}

@-webkit-keyframes slide-in {
  100% {
    -webkit-transform: translateX(0%);
  }
}

.slide-out {
  animation: slide-out 0.5s forwards;
  -webkit-animation: slide-out 0.5s forwards;
}

@keyframes slide-out {
  0% {
    transform: translateX(0%);
  }
  100% {
    transform: translateX(-60%);
  }
}

@-webkit-keyframes slide-out {
  0% {
    -webkit-transform: translateX(0%);
  }
  100% {
    -webkit-transform: translateX(-60%);
  }
}
</style>
