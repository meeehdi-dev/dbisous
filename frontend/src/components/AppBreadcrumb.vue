<script setup lang="ts">
import { ref, watchEffect } from "vue";
import { useConnections } from "@/composables/useConnections";
import { BreadcrumbItem } from "@nuxt/ui/dist/module";
import { useApp } from "@/composables/useApp";

const { database, schema, table } = useApp();
const { connections } = useConnections();

const items = ref<Array<BreadcrumbItem>>([]);

watchEffect(() => {
  const i: BreadcrumbItem[] = [];
  if (database.value) {
    i.push({
      label:
        connections.value.find((c) => c.id === database.value)?.name ||
        "Database",
      icon: "lucide:database",
      to: "/database",
    });
  }
  if (schema.value) {
    i.push({
      label: schema.value,
      icon: "lucide:table-of-contents",
      to: "/schema",
    });
  }
  if (table.value) {
    i.push({
      label: table.value,
      icon: "lucide:table",
      to: "/table",
    });
  }
  items.value = i;
});
</script>

<template>
  <UBreadcrumb :items="items" class="flex flex-initial p-4" />
</template>
