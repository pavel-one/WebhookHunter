<template>
  <div>
    <div style="display: flex; flex-direction: column; align-items: center;">
      <logo></logo>
      <br>
      <br>
      <btn class="center" icon="fa-solid fa-play">Получить</btn>
    </div>
  </div>
</template>

<script>
import Logo from "./components/Logo";
import Btn from "./components/Btn";
export default {
  components: {Btn, Logo},
  data: function () {
    return {
      message: 'Hi vuejs!',
    }
  }
}
</script>