<script setup lang="ts">
const { label, value, active } = defineProps<{
  label: string;
  value: string;
  active: boolean;
}>();
</script>

<template>
  <UButton
    :icon="`simple-icons:${value}`"
    :label="label"
    size="xl"
    :ui="{ label: 'w-full text-center' }"
    :trailing-icon="`lucide:square${active ? '-check-big' : ''}`"
  />
</template>
