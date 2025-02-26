<script setup lang="ts">
import { useRouter } from "vue-router";
import { useUrlParams } from "@/composables/useUrlParams";
import {
  formatColumns,
  FormattedQueryResult,
  RowAction,
} from "@/components/database/table/table";
import { useWails } from "@/composables/useWails";
import { GetSchemaTables } from "_/go/app/App";
import { ref } from "vue";
import { client } from "_/go/models";

const router = useRouter();
const { databaseId, schemaId } = useUrlParams();
const wails = useWails();

function navigateToTable(schemaId: string, tableId: string) {
  router.push({
    name: "table",
    params: { schemaId, tableId },
  });
}

const data = ref<FormattedQueryResult & { key: number }>();
const dataKey = ref(0);
const columns = ref<client.ColumnMetadata[]>();
const fetchingData = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  const result = await wails(() =>
    GetSchemaTables(databaseId.value, page, itemsPerPage, schemaId.value),
  );
  if (result instanceof Error) {
    // TODO: specific error handling
  } else {
    columns.value = result.columns;
    data.value = {
      key: dataKey.value++,
      ...result,
      columns: formatColumns(result.columns, undefined, undefined, true),
    };
  }
  fetchingData.value = false;
}
fetchData();
</script>

<template>
  <AppTabs>
    <template #data>
      <AppRows
        :loading="fetchingData"
        :data="data"
        :actions="[RowAction.View]"
        @view="
          (row) =>
            navigateToTable(
              // @ts-expect-error tkt
              row.TABLE_SCHEMA || row.table_schema || row.schema,
              // @ts-expect-error tkt
              row.TABLE_NAME || row.table_name || row.name,
            )
        "
        @pagination-change="fetchData"
      />
    </template>
    <template #info>
      <AppColumns :loading="fetchingData" :data="columns" />
    </template>
  </AppTabs>
</template>
