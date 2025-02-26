<script setup lang="ts">
const value = defineModel<string>();

const { disabled } = defineProps<{
  disabled: boolean;
}>();
</script>

<template>
  <UInput
    variant="ghost"
    v-model="value as string"
    :disabled="disabled"
    :ui="{
      base: 'pr-8 w-full overflow-ellipsis',
      trailing: 'pr-1',
    }"
  />
</template>
