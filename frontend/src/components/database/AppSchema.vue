<script setup lang="ts">
import { useRouter } from "vue-router";
import { useUrlParams } from "../../composables/useUrlParams";
import { RowAction } from "./table";
import { GetSchemaInfo, GetSchemaTables } from "../../../wailsjs/go/app/App";

const router = useRouter();
const { databaseId, schemaId } = useUrlParams();

function navigateToTable(schemaId: string, tableId: string) {
  router.push({
    name: "table",
    params: { schemaId, tableId },
  });
}
</script>

<template>
  <AppTabs>
    <template #data>
      <AppRowsTab
        :fetch-fn="
          (page, itemsPerPage) =>
            GetSchemaTables(databaseId, page, itemsPerPage, schemaId)
        "
        :actions="[RowAction.View]"
        @view="
          (row) =>
            navigateToTable(
              row.original.TABLE_SCHEMA ||
                row.original.table_schema ||
                row.original.schema,
              row.original.TABLE_NAME ||
                row.original.table_name ||
                row.original.name,
            )
        "
      />
    </template>
    <template #info>
      <!-- <AppRowsTab
        :fetch-fn="
          (page, itemsPerPage) =>
            GetSchemaInfo(databaseId, page, itemsPerPage, schemaId)
        "
      /> -->
    </template>
  </AppTabs>
</template>
