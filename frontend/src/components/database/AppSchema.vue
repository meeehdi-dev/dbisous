<script setup lang="ts">
import { useRouter } from "vue-router";
import { ref } from "vue";
import { useUrlParams } from "../../composables/useUrlParams";
import { useWails } from "../../wails";
import { Effect } from "effect";
import { FormattedQueryResult, RowAction } from "./table";
import { GetSchemaInfo, GetSchemaTables } from "../../../wailsjs/go/app/App";
import { formatQueryResult } from "../../effects/columns";

const wails = useWails();
const router = useRouter();
const { databaseId, schemaId } = useUrlParams();

const data = ref<FormattedQueryResult>();
const fetchingData = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  await Effect.runPromise(
    wails(() =>
      GetSchemaTables(databaseId.value, page, itemsPerPage, schemaId.value),
    ).pipe(
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

const info = ref<FormattedQueryResult>();
const fetchingInfo = ref(false);
async function fetchInfo(page = 1, itemsPerPage = 10) {
  fetchingInfo.value = true;
  await Effect.runPromise(
    wails(() =>
      GetSchemaInfo(databaseId.value, page, itemsPerPage, schemaId.value),
    ).pipe(
      Effect.andThen(formatQueryResult),
      Effect.tap((result) => {
        info.value = result;
        fetchingInfo.value = false;
      }),
      Effect.catchTags({
        WailsError: Effect.succeed,
      }),
    ),
  );
}
fetchInfo();

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
        @pagination-change="fetchData"
      />
    </template>
    <template #info>
      <AppRows
        :loading="fetchingInfo"
        :data="info"
        @pagination-change="fetchInfo"
      />
    </template>
  </AppTabs>
</template>
