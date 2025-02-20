<script setup lang="ts">
import { ref, watchEffect } from "vue";
import { useUrlParams } from "@/composables/useUrlParams";
import { useConnections } from "@/composables/useConnections";

const { databaseId, schemaId, tableId } = useUrlParams();
const { connections } = useConnections();

const items = ref<BreadcrumbItem[]>([]);

watchEffect(() => {
  const i: BreadcrumbItem[] = [];
  if (databaseId.value) {
    i.push({
      label:
        connections.value.find((c) => c.id === databaseId.value)?.name ||
        "Database",
      icon: "lucide:database",
      to: `/database/${databaseId.value}`,
    });
  }
  if (schemaId.value) {
    i.push({
      label: schemaId.value,
      icon: "lucide:table-of-contents",
      to: `/database/${databaseId.value}/schema/${schemaId.value}`,
    });
  }
  if (tableId.value) {
    i.push({
      label: tableId.value,
      icon: "lucide:table",
      to: `/database/${databaseId.value}/schema/${schemaId.value}/table/${tableId.value}`,
    });
  }
  items.value = i;
});
</script>

<template>
  <UBreadcrumb :items="items" class="flex flex-initial p-4" />
</template>
