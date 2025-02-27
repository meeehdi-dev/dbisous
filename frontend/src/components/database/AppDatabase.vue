<script setup lang="ts">
import { useRouter } from "vue-router";
import { GetDatabaseSchemas } from "_/go/app/App";
import { useUrlParams } from "@/composables/useUrlParams";
import {
  formatColumns,
  FormattedQueryResult,
  RowAction,
} from "@/components/database/table/table";
import { useWails } from "@/composables/useWails";
import { ref } from "vue";
import { client } from "_/go/models";

const router = useRouter();
const { databaseId } = useUrlParams();

async function navigateToSchema(schemaId: string) {
  await router.push({ name: "schema", params: { schemaId } });
}

const wails = useWails();

const data = ref<FormattedQueryResult & { key: number }>();
const dataKey = ref(0);
const columns = ref<client.ColumnMetadata[]>();
const fetchingData = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  const result = await wails(() =>
    GetDatabaseSchemas(databaseId.value, page, itemsPerPage),
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
            navigateToSchema(
              (row.SCHEMA_NAME || row.schema_name || row.name) as string,
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
