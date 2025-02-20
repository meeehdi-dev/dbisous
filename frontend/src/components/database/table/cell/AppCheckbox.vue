<script setup lang="ts">
const value = defineModel<boolean>();

// TODO: handle all props
const { disabled } = defineProps<{
  initialValue: boolean | undefined;
  defaultValue: boolean | undefined;
  nullable: boolean;
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
