import { computed } from "vue";
import { useRoute } from "vue-router";
import { createSharedComposable } from "@vueuse/core";

export const useUrlParams = createSharedComposable(() => {
  const route = useRoute();

  const databaseId = computed(() => {
    return Array.isArray(route.params.databaseId)
      ? route.params.databaseId[0]
      : route.params.databaseId || "";
  });

  const schemaId = computed(() => {
    return Array.isArray(route.params.schemaId)
      ? route.params.schemaId[0]
      : route.params.schemaId || "";
  });

  const tableId = computed(() => {
    return Array.isArray(route.params.tableId)
      ? route.params.tableId[0]
      : route.params.tableId || "";
  });

  return {
    databaseId,
    schemaId,
    tableId,
  };
});
