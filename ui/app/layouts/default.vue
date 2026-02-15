<script setup lang="ts">
import CustomBreadcumb from "~/components/CustomBreadcumb.vue";
import RightPanel from "~/components/RightPanel.vue";
</script>

<template>
  <div class="min-h-screen mx-auto bg-surface-50 dark:bg-surface-950 px-4">
    <Toast />
    <Header />
    <RightPanel />

    <MobileSidebar />

    <main
      class="pt-10 transition-all duration-300"
      :class="{
        'lg:ml-56': sidebarOpen,
        'lg:ml-16': !sidebarOpen,
      }"
    >
      <CustomBreadcumb />
      <div class="p-4">
        <slot />
      </div>
    </main>
    <Footer />
  </div>
</template>

<style lang="css" scoped></style>
