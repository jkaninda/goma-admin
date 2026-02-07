
<script setup lang="ts">
import { onMounted } from 'vue';
onMounted(() => {});
</script>

<template>
   <div>
    <p>Some default layout content shared across all pages</p>
    <slot />
  </div>
</template>

<style lang="css" scoped>
</style>
