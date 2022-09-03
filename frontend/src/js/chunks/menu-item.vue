<template>
  <div class="menu-item">
    <div class="name">{{ name }}</div>
    <div class="count">Сообщений ({{ count }})</div>
    <div class="date">{{ date }}</div>
  </div>
</template>

<script>
export default {
  props: {
    name: String,
    count: Number,
    date: String
  }
}
</script>