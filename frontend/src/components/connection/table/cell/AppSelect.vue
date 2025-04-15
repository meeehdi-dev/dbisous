<script setup lang="ts">
const value = defineModel<string>();

const { disabled, items } = defineProps<{
  disabled?: boolean;
  items: Array<string>;
}>();
</script>

<template>
  <USelect
    v-model="value"
    variant="ghost"
    :items="items"
    :disabled="disabled"
    :ui="{ base: 'w-full' }"
  />
</template>
