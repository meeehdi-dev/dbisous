<script setup lang="ts">
const value = defineModel<number>();

const { disabled } = defineProps<{
  disabled: boolean;
}>();
</script>

<template>
  <UInputNumber variant="ghost" v-model="value" :disabled="disabled" />
</template>
