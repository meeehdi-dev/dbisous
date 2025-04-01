<script setup lang="ts">
const value = defineModel<string>();

const { disabled } = defineProps<{
  disabled: boolean;
}>();
</script>

<template>
  <UInput
    v-model="value"
    spellcheck="false"
    variant="ghost"
    :disabled="disabled"
    :ui="{
      base: 'w-full overflow-ellipsis',
      trailing: 'pr-1',
    }"
  />
</template>
