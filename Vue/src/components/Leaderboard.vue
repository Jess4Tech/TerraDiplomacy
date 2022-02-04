<template>
    <div
    class="
      flex
      justify-center
      items-center
      flex-col
      md:bg-black
      w-screen
      h-screen
    "
  >
    <ul
      class="
        w-2/3
        h-2/3
        bg-white
        sm:w-screen
        sm:h-screen
        md:w-2/3
        md:h-3/4
        md:rounded-lg
        md:shadow-md
        flex
        justify-center
        items-center
      "
    >
      <li v-if="data.entries.length == 0">
        <p class="font-bold">No entries were found</p>
      </li>
      <li
        v-bind:key="entry.id"
        v-for="entry in data.entries"
        class="border-2 shadow-md w-2/3 m-2 p-2 relative"
      >
        <div class="flex justify-between">
          <p class="font-medium">{{ entry.id }}</p>
        </div>
        <div class="overflow-auto">{{ entry.tension }}</div>
      </li>
    </ul>
  </div>
</template>

<script setup>
import {reactive} from '@vue/reactivity';
import {onMounted} from '@vue/runtime-core';
import {useRouter} from 'vue-router';
import {getTension, ifAuthorized, redirectIfUnauthorized} from '../utility/api';

const router = useRouter();

const data = reactive({
  entries: [],
});

/**
 * Internal function to fetch data, saves boilerplate
 */
async function fetchData() {
  redirectIfUnauthorized(router);
  ifAuthorized(
      async () => {
        data.entries = await getTension();
      },
  );
}

onMounted(async () => {
  await fetchData();
});

</script>
