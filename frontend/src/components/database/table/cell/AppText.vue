<script setup lang="ts">
const value = defineModel<string>();

const { initialValue, disabled } = defineProps<{
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
    :highlight="value !== initialValue"
    :color="value !== initialValue ? 'warning' : undefined"
    :ui="{
      root: 'w-full',
      base: 'pr-8 w-full overflow-ellipsis',
      trailing: 'pr-1',
    }"
  />
</template>
