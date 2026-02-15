<script setup lang="ts">
const props = defineProps({
  variant: {
    type: String,
    default: "default",
  },
});

const colorMode = useColorMode();

function toggleDarkMode() {
  colorMode.preference = colorMode.value === "dark" ? "light" : "dark";
}
</script>

<template>
  <ClientOnly>
    <Button
      text
      rounded
      @click="toggleDarkMode"
      :icon="colorMode.value === 'dark' ? 'pi pi-moon' : 'pi pi-sun'"
      :class="{
        'text-primary-contrast! hover:bg-white/10! transition-colors duration-150':
          props.variant === 'toolbar',
      }"
    />
  </ClientOnly>
</template>
