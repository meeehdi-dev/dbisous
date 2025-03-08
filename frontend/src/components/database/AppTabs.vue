<script setup lang="ts">
const { defaultQuery } = defineProps<{ defaultQuery?: string }>();

defineSlots<{
  rows: unknown;
  columns: unknown;
}>();

const tabs = [
  {
    label: "Rows",
    slot: "rows",
    icon: "lucide:list",
  },
  {
    label: "Columns",
    slot: "columns",
    icon: "lucide:info",
  },
  {
    label: "Script",
    slot: "script",
    icon: "lucide:square-terminal",
  },
  {
    label: "Export",
    slot: "export",
    icon: "lucide:upload",
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
    <template #rows>
      <slot name="rows" />
    </template>
    <template #columns>
      <slot name="columns" />
    </template>
    <template #script>
      <AppScript :default-query="defaultQuery" />
    </template>
    <template #export>
      <AppExport />
    </template>
  </UTabs>
</template>
