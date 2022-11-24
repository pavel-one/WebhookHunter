<template>
  <home v-if="!login"></home>
  <dashboard v-else></dashboard>
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
        let redirectUrl = url.join(".")
        redirectUrl = window.location.protocol + "//" + redirectUrl
        if (window.location.port) {
          redirectUrl = redirectUrl + ':' + window.location.port
        }
        window.location.href = redirectUrl + "/ui/"
        return
      }

      this.login = true
    }
  }
}
</script>