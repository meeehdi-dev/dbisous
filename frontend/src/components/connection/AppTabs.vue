<script setup lang="ts">
import AppRows from "@/components/connection/table/AppRows.vue";
import AppColumns from "@/components/connection/table/column/AppColumns.vue";
import { Tab } from "@/utils/tabs";

const { defaultQuery } = defineProps<{ defaultQuery?: string }>();
const active = defineModel<Tab>({ default: Tab.Rows });

defineSlots<{
  rows: typeof AppRows;
  columns: typeof AppColumns;
}>();

const tabs = [
  {
    label: "Rows",
    value: Tab.Rows,
    slot: "rows",
    icon: "lucide:list",
  },
  {
    label: "Columns",
    value: Tab.Columns,
    slot: "columns",
    icon: "lucide:info",
  },
  {
    label: "Query",
    value: Tab.Query,
    slot: "query",
    icon: "lucide:square-terminal",
  },
  {
    label: "Export",
    value: Tab.Export,
    slot: "export",
    icon: "lucide:upload",
  },
  {
    label: "Import",
    value: Tab.Import,
    slot: "import",
    icon: "lucide:download",
  },
];
</script>

<template>
  <UTabs
    v-model="active"
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
    <template #query>
      <AppQuery :default-query="defaultQuery" />
    </template>
    <template #export>
      <AppExport />
    </template>
    <template #import>
      <AppImport />
    </template>
  </UTabs>
</template>
