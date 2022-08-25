<template>
  <div>
    <logo></logo>
  </div>
</template>

<script>
import Logo from "./components/Logo";
export default {
  components: {Logo},
  data: function () {
    return {
      message: 'Hi vuejs!',
    }
  }
}
</script>