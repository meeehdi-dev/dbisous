<script setup lang="ts">
const value = defineModel<number>();

const { disabled } = defineProps<{
  disabled: boolean;
}>();
</script>

<template>
  <UInputNumber
    v-model="value"
    variant="ghost"
    orientation="vertical"
    :disabled="disabled"
  />
</template>
