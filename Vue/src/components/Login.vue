<template>
  <div
    class="md:bg-black md:w-screen md:h-screen flex justify-center items-center"
  >
    <div
      @click="auth.error = ''"
      class="
        top-0
        left-0
        w-full
        h-auto
        absolute
        bg-red-600
        rounded-bl-md rounded-br-md
        text-center text-white
        cursor-pointer
      "
      v-if="data.error != ''"
    >
      <span class="p-4 m-4">{{ data.error }}</span>
    </div>
    <div
      class="
        bg-white
        sm:w-screen
        sm:h-screen
        md:w-1/3
        md:h-3/4
        md:rounded-lg
        md:shadow-md
        flex flex-col
        justify-center
        items-center
      "
    >
      <img
        src="../assets/logo.png"
        alt="Terra Logo"
        class="m-4 object-scale-down w-1/2 h-1/2"
      />
      <form class="flex flex-col justify-center items-center w-full w-full">
        <input
          type="text"
          placeholder="Username"
          v-model="data.user"
          class="border border-gray-500 rounded-md m-2 w-1/2"
        />
        <input
          type="password"
          placeholder="Password"
          v-model="data.otac"
          class="border border-gray-500 rounded-md m-2 w-1/2"
        />

        <input
          type="submit"
          value="Submit"
          class="
            bg-gray-300
            rounded
            p-2
            m-2
            hover:bg-gray-400
            cursor-pointer
            sm:w-1/2
            md:w-1/4
          "
          @click="login"
        />
      </form>
    </div>
  </div>
</template>

<script setup>
import {useRouter} from 'vue-router';
import {onMounted, reactive} from 'vue';
import {API_URL, ifAuthorized} from '../utility/api';

const router = useRouter();

onMounted(async () => {
  ifAuthorized(async () => {
    console.log('Not logged in, redirecting');
    await router.push({name: 'LoginPage'});
  });
});

const data = reactive({
  user: '',
  otac: '',
  error: '',
});

const login = async function(event) {
  event.preventDefault();
  if (data.user.trim() == '' || data.otac.trim() == '') {
    auth.error = 'Make sure all fields are filled in correctly';
    return;
  }
  const res = await fetch(`${API_URL}/auth/login`, {
    method: 'POST',
    mode: 'cors',
    credentials: 'include',
    body: JSON.stringify({user: data.user, otac: data.otac}),
  });

  if (res.ok) {
    await router.push({name: 'ProjectList'});
  } else if (res.status == 401) {
    data.error = 'Invalid credentials';
  } else {
    data.error = 'An unknown error has occurred';
    console.log(await res.text());
  }
};
</script>
