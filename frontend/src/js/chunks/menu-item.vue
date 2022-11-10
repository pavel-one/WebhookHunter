<template>
  <div class="menu-item">
    <div class="name">{{ name }}</div>
    <div class="count">Сообщений ({{ count }})</div>
    <div class="date">{{ date }}</div>
    <div class="buttons" @click="this.delete">
      <font-awesome-icon icon="fa-solid fa-trash-can-arrow-up" />
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
    date: String
  },
  methods: {
    delete: function () {
      axios.delete('/api/v1/channels/'+this.id)
    }
  }
}
</script>