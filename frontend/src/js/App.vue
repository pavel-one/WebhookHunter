<template>
  <div>
    <home v-if="!login"></home>
    <dashboard v-else></dashboard>
  </div>
</template>

<script>
import Home from "./Home";
import Dashboard from "./Dashboard";
export default {
  components: {Dashboard, Home},
  data: function () {
    return {
      login: false
    }
  },
  mounted() {
    console.log(this.login)
  }
}
</script>