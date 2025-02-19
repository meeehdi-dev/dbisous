<script setup lang="ts">
import { useRouter } from "vue-router";
import {
  GetDatabaseInfo,
  GetDatabaseSchemas,
} from "../../../wailsjs/go/app/App";
import { useUrlParams } from "../../composables/useUrlParams";
import { RowAction } from "./table";

const router = useRouter();
const { databaseId } = useUrlParams();

function navigateToSchema(schemaId: string) {
  router.push({ name: "schema", params: { schemaId } });
}
</script>

<template>
  <AppTabs>
    <template #data>
      <AppRowsTab
        :fetch-fn="
          (page, itemsPerPage) =>
            GetDatabaseSchemas(databaseId, page, itemsPerPage)
        "
        :actions="[RowAction.View]"
        @view="
          (row) =>
            navigateToSchema(
              row.original.SCHEMA_NAME ||
                row.original.schema_name ||
                row.original.name,
            )
        "
      />
    </template>
    <template #info>
      <AppRowsTab
        :fetch-fn="
          (page, itemsPerPage) =>
            GetDatabaseInfo(databaseId, page, itemsPerPage)
        "
      />
    </template>
  </AppTabs>
</template>
