<template>
  <div class="log-item">
    <div class="date">
      {{ date }}
    </div>
    <div class="data">
      <json-viewer theme="my-awesome-json-theme" :value="request"></json-viewer>
      <json-viewer theme="my-awesome-json-theme" :value="headers"></json-viewer>
      <json-viewer theme="my-awesome-json-theme" :value="query"></json-viewer>
    </div>
  </div>
</template>

<script>
import JsonViewer from 'vue-json-viewer'
import 'vue-json-viewer/style.css'
export default {
  components: {JsonViewer},
  props: {
    date: String,
    request: JSON,
    headers: JSON,
    query: JSON,
    path: String,
  }
}
</script>