<script setup lang="ts">
const { locale, locales, t } = useI18n();

const menu = ref();

const languageItems = computed(() => {
  return locales.value.map((l: any) => ({
    label: l.name,
    code: l.code,
    command: () => {
      locale.value = l.code;
    }
  }));
});

function toggleMenu(event: Event) {
  menu.value.toggle(event);
}
</script>

<template>
  <div class="flex items-center">
    <Button
      text
      rounded
      size="small"
      aria-label="common.language_menu"
      aria-haspopup="true"
      aria-controls="language_menu"
      @click="toggleMenu"
      class="text-primary-contrast! hover:bg-white/10! min-w-10 h-8 flex items-center justify-center gap-1"
    >
      <span class="text-[11px] font-bold uppercase tracking-tighter">
        {{ locale }}
      </span>
     
    </Button>

    <Menu 
      ref="menu" 
      id="language_menu" 
      :model="languageItems" 
      :popup="true"
    >
      <template #item="{ item, props }">
        <a v-bind="props.action" class="flex items-center gap-2 px-1.5 py-1.5">
          <span class="text-[10px] font-mono bg-surface-100 dark:bg-surface-800 px-1 rounded uppercase">
            {{ item.code }}
          </span>
          <span class="text-sm">{{ item.label }}</span>
        </a>
      </template>
    </Menu>
  </div>
</template>

<style scoped>

</style>