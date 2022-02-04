<template>
  <div
    class="md:bg-black md:w-screen md:h-screen flex justify-center items-center"
  >
    <div
      class="
        bg-white
        sm:w-screen
        sm:h-screen
        md:w-1/3
        md:h-1/2
        md:rounded-lg
        md:shadow-md
        flex flex-col
        justify-center
        items-center
      "
    >
      <h1 class="text-lg font-bold">Not found</h1>
      <p>That page wasn't found, sorry</p>
      <router-link class="text-blue-700"
      :to="{ name: 'ProjectList' }"
        >Go home?</router-link
      >
    </div>
  </div>
</template>
