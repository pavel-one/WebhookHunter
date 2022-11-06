<template>
  <div class="btn">
    <div class="btn-wrapper">
      <div v-if="icon" class="btn-icon">
        <font-awesome-icon :icon="icon" />
      </div>
      <div class="btn-text">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    icon: String
  }
}
</script>