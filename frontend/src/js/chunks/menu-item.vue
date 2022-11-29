<template>
  <div class="menu-item" :class="{active: active}" @click="$emit('changeChannel', name)">
    <div class="name">{{ name }}</div>
    <div class="count">Сообщений ({{ count }})</div>
    <div class="date">{{ date }}</div>
    <div class="buttons" @click="this.delete">
      <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-trash-x" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
        <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
        <path d="M4 7h16"></path>
        <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12"></path>
        <path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3"></path>
        <path d="M10 12l4 4m0 -4l-4 4"></path>
      </svg>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  props: {
    id: Number,
    name: String,
    count: Number,
    date: String,
    active: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    delete: function () {
      axios.delete('/api/v1/channels/'+this.id)
    }
  }
}
</script>