<script setup lang="ts">
const value = defineModel<boolean>();

const { disabled } = defineProps<{
  disabled: boolean;
}>();
</script>

<template>
  <UCheckbox
    v-model="value"
    :disabled="disabled"
    :ui="{
      root: 'pl-2 items-center',
    }"
  />
</template>
