<template>
  <div class="home-page">
    <div style="display: flex; flex-direction: column; align-items: center;">
      <logo></logo>
      <br>
      <br>
      <btn @click="click" class="center lg" icon="fa-solid fa-play">Получить</btn>
    </div>
  </div>
</template>

<script>
import Logo from "./components/Logo";
import Btn from "./components/Btn";
import axios from "axios";
export default {
  components: {Btn, Logo},
  methods: {
    click: async function () {
      let response = await axios.post('/api/v1/', {})

      if (response.data.uri) {
        let redirect = window.location.protocol + "//" + response.data.uri + "." + window.location.hostname
        if (window.location.port) {
          redirect = redirect + ':' + window.location.port
        }
        window.location.href = redirect
      }

      console.log()
    }
  }
}
</script>