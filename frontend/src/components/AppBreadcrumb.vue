<script setup lang="ts">
import { ref, watch } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

const items = ref<BreadcrumbItem[]>([]);

watch(router.currentRoute, () => {
  const { params } = router.currentRoute.value;
  const i: BreadcrumbItem[] = [];
  if (params.database) {
    i.push({
      label: "Database",
      icon: "lucide:database",
      to: `/database/${params.database}`,
    });
  }
  if (params.schema) {
    i.push({
      label: params.schema,
      icon: "lucide:table-of-contents",
      to: `/database/${params.database}/schema/${params.schema}`,
    });
  }
  if (params.table) {
    i.push({
      label: params.table,
      icon: "lucide:table",
      to: `/database/${params.database}/schema/${params.schema}/table/${params.table}`,
    });
  }
  items.value = i;
});
</script>

<template>
  <UBreadcrumb :items="items" class="p-4" />
</template>
