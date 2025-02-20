<script setup lang="ts">
const value = defineModel<number>();

// TODO: handle all props
// eslint-disable-next-line @typescript-eslint/no-unused-vars
const { defaultValue, nullable, initialValue, disabled } = defineProps<{
  initialValue: number | undefined;
  defaultValue: number | undefined;
  nullable: boolean;
  disabled: boolean;
}>();
</script>

<template>
  <UInputNumber
    variant="ghost"
    :value="value"
    :disabled="disabled"
    class="w-full"
  />
</template>
