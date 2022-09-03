<template>
  <div>
    <home v-if="!login"></home>
    <dashboard v-else></dashboard>
  </div>
</template>

<script>
import Home from "./Home";
import Dashboard from "./Dashboard";
import axios from "axios";
export default {
  components: {Dashboard, Home},
  data: function () {
    return {
      login: false
    }
  },
  async mounted() {
    let url = window.location.hostname.split(".")

    // if the url is subdomain
    if (url[0] !== "ww" && url[0] !== "hunter") {
      try {
        const resp = await axios.post("/api/v1/check/", {
          "uri": url[0]
        })
      } catch (e) {
        url.splice(0, 1)
        const redirectUrl = url.join(".")
        window.location.href = window.location.protocol + "//" + redirectUrl
        return
      }

      this.login = true
    }
  }
}
</script>