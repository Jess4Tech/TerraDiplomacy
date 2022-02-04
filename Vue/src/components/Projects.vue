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
      <li v-if="data.projects.length == 0">
        <p class="font-bold">No projects were found</p>
      </li>
      <li
        v-bind:key="proj.name"
        v-for="proj in data.projects"
        class="border-2 shadow-md w-2/3 m-2 p-2 relative"
      >
        <div class="flex justify-between">
          <p class="font-medium">{{ proj.name }}</p>
          <p class="font-light">({{ proj.weight }})</p>
        </div>
        <div class="overflow-auto">{{ proj.description }}</div>
        <p
          v-if="data.tier >= 3"
          class="font-light absolute bottom-2 right-2 cursor-pointer"
          @click="deleteProjectReload(proj.name)"
        >
          x
        </p>
      </li>
    </ul>
  </div>
</template>

<script setup>
import {getProjects, deleteProject, ifAuthorized, redirectIfUnauthorized}
  from '../utility/api';
import {useRouter} from 'vue-router';
import {reactive, onMounted} from 'vue';

const router = useRouter();

const data = reactive({
  projects: [],
  tier: 0,
});

/**
 * Internal function to fetch data, saves boilerplate
 */
async function fetchData() {
  await redirectIfUnauthorized(router);
  ifAuthorized(
      async (res) => {
        data.projects = await getProjects();
        data.tier = res.tier;
      },
  );
}

const deleteProjectReload = async function(name) {
  await deleteProject(name);
  await fetchData();
};

onMounted(async function() {
  await fetchData();
});
</script>
