<script setup lang="ts">
const value = defineModel<string>();

// TODO: handle all props
// eslint-disable-next-line @typescript-eslint/no-unused-vars
const { defaultValue, nullable, initialValue, disabled } = defineProps<{
  initialValue: string | undefined;
  defaultValue: string | undefined;
  nullable: boolean;
  disabled: boolean;
}>();
</script>

<template>
  <UInput variant="ghost" :value="value" :disabled="disabled" class="w-full" />
</template>
