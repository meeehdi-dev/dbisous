<script setup lang="ts">
import { useRouter } from "vue-router";
import { ref } from "vue";
import { useUrlParams } from "../../composables/useUrlParams";
import { useWails } from "../../wails";
import { Effect } from "effect";
import { formatColumns, FormattedQueryResult, RowAction } from "./table";
import { GetSchemaInfo, GetSchemaTables } from "../../../wailsjs/go/app/App";

const wails = useWails();
const router = useRouter();
const { databaseId, schemaId } = useUrlParams();

const tabs = [
  {
    label: "Tables",
    slot: "data",
    icon: "lucide:list",
  },
  {
    label: "Info",
    slot: "info",
    icon: "lucide:info",
  },
  {
    label: "Script",
    slot: "script",
    icon: "lucide:square-terminal",
  },
];

const data = ref<FormattedQueryResult>();
const info = ref<FormattedQueryResult>();
async function getData(page = 1, itemsPerPage = 10) {
  await Effect.runPromise(
    wails(() =>
      GetSchemaTables(databaseId.value, page, itemsPerPage, schemaId.value),
    ).pipe(
      Effect.tap((result) => {
        data.value = {
          ...result,
          columns: formatColumns(result.columns),
        };
      }),
      Effect.catchTags({
        WailsError: Effect.succeed,
      }),
    ),
  );
}
async function getInfo(page = 1, itemsPerPage = 10) {
  await Effect.runPromise(
    wails(() =>
      GetSchemaInfo(databaseId.value, page, itemsPerPage, schemaId.value),
    ).pipe(
      Effect.tap((result) => {
        info.value = {
          ...result,
          columns: formatColumns(result.columns, false),
        };
      }),
      Effect.catchTags({
        WailsError: Effect.succeed,
      }),
    ),
  );
}
getData();
getInfo();

function navigateToTable(schemaId: string, tableId: string) {
  router.push({
    name: "table",
    params: { schemaId, tableId },
  });
}
</script>

<template>
  <UTabs
    :items="tabs"
    variant="link"
    :ui="{ root: 'h-full', content: 'flex flex-1 flex-col gap-2' }"
  >
    <template #data>
      <AppRows
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
        @pagination-change="getData"
      />
    </template>
    <template #info>
      <AppRows :data="info" @pagination-change="getInfo" />
    </template>
    <template #script>
      <AppScript />
    </template>
  </UTabs>
</template>
