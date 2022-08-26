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
export default {
  components: {Btn, Logo},
  methods: {
    click: function () {
      console.log('test')
    }
  }
}
</script>