<script setup>
const colorMode = useColorMode()

function toggleDarkMode() {
  colorMode.preference =
    colorMode.value === 'dark'
      ? 'light'
      : 'dark'
}
</script>

<template>
  <div class="fixed top-8 right-8">
    <Button
      type="button"
      @click="toggleDarkMode"
      rounded
      severity="secondary"
      :icon="colorMode.value === 'dark'
        ? 'pi pi-moon'
        : 'pi pi-sun'"
    />
  </div>
</template>
