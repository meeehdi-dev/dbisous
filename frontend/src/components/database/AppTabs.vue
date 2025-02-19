<script setup lang="ts">
const { defaultQuery } = defineProps<{ defaultQuery?: string }>();

defineSlots<{
  data: unknown;
  info: unknown;
}>();

const tabs = [
  {
    label: "Data",
    slot: "data",
    icon: "lucide:list",
  },
  {
    label: "Info",
    slot: "info",
    icon: "lucide:info",
  },
  {
    label: "Script",
    slot: "script",
    icon: "lucide:square-terminal",
  },
];
</script>

<template>
  <UTabs
    :items="tabs"
    variant="link"
    :ui="{
      root: 'flex flex-auto overflow-hidden',
      content: 'flex flex-auto flex-col gap-2 overflow-hidden',
    }"
  >
    <template #data>
      <slot name="data" />
    </template>
    <template #info>
      <slot name="info" />
    </template>
    <template #script>
      <AppScript :default-query="defaultQuery" />
    </template>
  </UTabs>
</template>
