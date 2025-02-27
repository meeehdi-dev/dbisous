<script setup lang="ts">
const value = defineModel<string>();

const { disabled } = defineProps<{
  disabled: boolean;
}>();
</script>

<template>
  <UInput variant="ghost" v-model="value" :disabled="disabled" />
</template>
