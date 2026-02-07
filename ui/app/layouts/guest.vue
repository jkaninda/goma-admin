
<script setup lang="ts">
</script>

<template>
   <div class="min-h-screen flex flex-col justify-center items-center mx-auto bg-surface-50 dark:bg-surface-950 px-4">
   <Toast/>
   <div class="fixed top-3 right-4 sm:top-8 sm:right-8">
   <ThemeToggle/>
   </div>
      <main class="w-full max-w-lg">
        <slot />
      </main>
    </div>
</template>

<style lang="css" scoped>
</style>