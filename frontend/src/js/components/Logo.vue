<template>
  <div class="logo-wrapper">
    <div v-if="!small" class="logo">
      <div class="text">
        PAVEL
      </div>
      <div class="dot">.</div>
      <div class="text">
        ONE
      </div>
    </div>
    <div class="post-title">HOOK HUNTER</div>
    <slot></slot>
  </div>
</template>

<script>
export default {
  props: {
    small: Boolean
  }
}
</script>