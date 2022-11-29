<template>
  <div class="log-item">
    <div class="buttons">
      <div class="button" @click="this.preview = !this.preview" :class="{active: preview}">
        <svg v-if="!preview" xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-eye"
             viewBox="0 0 24 24"
             stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
          <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
          <circle cx="12" cy="12" r="2"></circle>
          <path d="M22 12c-2.667 4.667 -6 7 -10 7s-7.333 -2.333 -10 -7c2.667 -4.667 6 -7 10 -7s7.333 2.333 10 7"></path>
        </svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-eye-off" viewBox="0 0 24 24"
             stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
          <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
          <line x1="3" y1="3" x2="21" y2="21"></line>
          <path d="M10.584 10.587a2 2 0 0 0 2.828 2.83"></path>
          <path
              d="M9.363 5.365a9.466 9.466 0 0 1 2.637 -.365c4 0 7.333 2.333 10 7c-.778 1.361 -1.612 2.524 -2.503 3.488m-2.14 1.861c-1.631 1.1 -3.415 1.651 -5.357 1.651c-4 0 -7.333 -2.333 -10 -7c1.369 -2.395 2.913 -4.175 4.632 -5.341"></path>
        </svg>
      </div>
    </div>

    <div class="date">
      {{ splitData[0] }} <br>
      <span class="small">
        {{ splitData[1] }}
      </span>
    </div>
    <div class="data">
      <div>
        <h2>Request</h2>
        <json-viewer theme="my-awesome-json-theme"
                     :show-double-quotes="preview"
                     :preview-mode="preview"
                     :copyable="true"
                     :value="request"></json-viewer>
      </div>
      <div>
        <h2>Headers</h2>
        <json-viewer theme="my-awesome-json-theme monospace"
                     :show-double-quotes="preview"
                     :preview-mode="preview"
                     :copyable="true"
                     :value="headers"></json-viewer>
      </div>
      <div>
        <h2>Query</h2>
        <json-viewer theme="my-awesome-json-theme"
                     :show-double-quotes="preview"
                     :preview-mode="preview"
                     :copyable="true"
                     :value="query"></json-viewer>
      </div>
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
  },
  data: function () {
    return {
      splitData: Array,
      preview: false,
      expand: 1
    }
  },
  mounted() {
    this.splitData = this.date.split(" ")
  },
  watch: {
    preview(n) {
      if (n === true) {
        this.expand = 3
      } else {
        this.expand = 1
      }
    }
  },
}
</script>

<style lang="scss">
.log-item {
  position: relative;

  .buttons {
    position: absolute;
    left: 110px;
    top: 13px;
    z-index: 9999;

    .button {
      width: 26px;
      cursor: pointer;
      color: #49b3ff;
      transition: .25s;

      &:hover, &.active {
        color: #3ce09d;
      }
    }
  }

}

.data {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: flex-start;

  & > div {
    width: 31%;
    margin: 10px;
  }
}

.date .small {
  font-size: .8em;
  color: #525252;
  margin-top: 5px;
  text-align: right;
  font-weight: bold;
}
</style>