<script setup lang="ts">
import { ref, watchEffect } from "vue";
import { useConnections } from "@/composables/shared/useConnections";
import type { BreadcrumbItem } from "@nuxt/ui";
import { useApp } from "@/composables/shared/useApp";

const { connection, database, schema, table } = useApp();
const { connections } = useConnections();

interface AppBreadcrumbItem extends BreadcrumbItem {
  onClick?: () => void;
}

const items = ref<Array<AppBreadcrumbItem>>([]);

watchEffect(() => {
  const i: AppBreadcrumbItem[] = [];
  if (connection.value) {
    const c = connections.value.find((c) => c.id === connection.value);
    i.push({
      label: c?.name || c?.connection_string || "Connection",
      icon: "lucide:house-plug",
      to: "/connection",
      onClick: () => {
        database.value = "";
        schema.value = "";
        table.value = "";
      },
    });
  }
  if (database.value) {
    i.push({
      label: database.value,
      icon: "lucide:database",
      to: "/database",
      onClick: () => {
        schema.value = "";
        table.value = "";
      },
    });
  }
  if (schema.value) {
    i.push({
      label: schema.value,
      icon: "lucide:table-of-contents",
      to: "/schema",
      onClick: () => {
        table.value = "";
      },
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
  <UBreadcrumb :items="items" class="flex flex-initial p-2">
    <template #item="{ item }">
      <ULink
        :to="item.to"
        class="flex items-center gap-1"
        @click="item.onClick"
      >
        <UIcon v-if="item.icon" class="size-5" :name="item.icon" />
        <span>
          {{ item.label }}
        </span>
      </ULink>
    </template>
  </UBreadcrumb>
</template>
