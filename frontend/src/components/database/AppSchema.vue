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

async function navigateToTable(schemaId: string, tableId: string) {
  await router.push({
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
    GetSchemaTables(
      databaseId.value,
      {
        offset: (page - 1) * itemsPerPage,
        limit: itemsPerPage,
        filter: [],
        order: [],
      },
      schemaId.value,
    ),
  );
  fetchingData.value = false;
  if (result instanceof Error) {
    return;
  }
  columns.value = result.columns;
  data.value = {
    key: dataKey.value++,
    // eslint-disable-next-line @typescript-eslint/no-misused-spread
    ...result,
    columns: formatColumns(result.columns, undefined, undefined, true),
  };
}
await fetchData();
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
              (row.TABLE_SCHEMA || row.table_schema || row.schema) as string,
              (row.TABLE_NAME || row.table_name || row.name) as string,
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
