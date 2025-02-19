<script setup lang="ts">
import { useUrlParams } from "../../composables/useUrlParams";
import { GetTableInfo, GetTableRows } from "../../../wailsjs/go/app/App";

const { databaseId, schemaId, tableId } = useUrlParams();
</script>

<template>
  <AppTabs :default-query="`SELECT * FROM ${tableId};`">
    <template #data>
      <AppRowsTab
        :fetch-fn="
          (page, itemsPerPage) =>
            GetTableRows(databaseId, page, itemsPerPage, schemaId, tableId)
        "
      />
    </template>
    <template #info>
      <AppRowsTab
        :fetch-fn="
          (page, itemsPerPage) =>
            GetTableInfo(databaseId, page, itemsPerPage, schemaId, tableId)
        "
      />
    </template>
  </AppTabs>
</template>
