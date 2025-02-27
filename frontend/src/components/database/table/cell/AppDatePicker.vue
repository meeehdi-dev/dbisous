<script setup lang="ts">
const value = defineModel<string>();

const { disabled } = defineProps<{
  disabled: boolean;
}>();
</script>

<template>
  <UInput v-model="value" variant="ghost" :disabled="disabled" />
</template>
