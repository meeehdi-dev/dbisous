<script setup lang="ts">
const value = defineModel<string>();

const { isNew, initialValue, disabled } = defineProps<{
  isNew: boolean;
  initialValue: string | undefined;
  defaultValue: string | undefined;
  nullable: boolean;
  disabled: boolean;
}>();
</script>

<template>
  <UInput
    variant="ghost"
    v-model="value as string"
    :disabled="disabled"
    :highlight="isNew || value !== initialValue"
    :color="isNew ? 'primary' : value !== initialValue ? 'warning' : undefined"
    :ui="{
      base: 'pr-8 w-full overflow-ellipsis',
      trailing: 'pr-1',
    }"
  />
</template>
