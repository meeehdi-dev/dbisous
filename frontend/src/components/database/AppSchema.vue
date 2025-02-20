<script setup lang="ts">
import { useRouter } from "vue-router";
import { useUrlParams } from "@/composables/useUrlParams";
import {
  FormattedQueryResult,
  RowAction,
} from "@/components/database/table/table";
import { useWails } from "@/composables/useWails";
import { GetSchemaTables } from "_/go/app/App";
import { Effect } from "effect";
import { ref } from "vue";
import { client } from "_/go/models";
import { formatQueryResult } from "@/effects/columns";

const router = useRouter();
const { databaseId, schemaId } = useUrlParams();
const wails = useWails();

function navigateToTable(schemaId: string, tableId: string) {
  router.push({
    name: "table",
    params: { schemaId, tableId },
  });
}

const data = ref<FormattedQueryResult>();
const columns = ref<client.ColumnMetadata[]>();
const fetchingData = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  await Effect.runPromise(
    wails(() =>
      GetSchemaTables(databaseId.value, page, itemsPerPage, schemaId.value),
    ).pipe(
      Effect.tap((result) => {
        columns.value = result.columns;
      }),
      Effect.andThen(formatQueryResult),
      Effect.tap((result) => {
        data.value = result;
        fetchingData.value = false;
      }),
      Effect.catchTags({
        WailsError: Effect.succeed,
      }),
    ),
  );
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
      <AppColumns :loading="fetchingData" :data="columns" />
    </template>
  </AppTabs>
</template>
