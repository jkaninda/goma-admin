<template>
  <div class="guest-layout">
    <router-view />
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.guest-layout {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: var(--bg-secondary);
  padding: 20px;
}
</style>
