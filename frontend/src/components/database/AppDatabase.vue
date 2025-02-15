<script setup lang="ts">
import { useRouter } from "vue-router";
import { ref } from "vue";
import { Effect } from "effect";
import {
  GetDatabaseInfo,
  GetDatabaseSchemas,
} from "../../../wailsjs/go/app/App";
import { useUrlParams } from "../../composables/useUrlParams";
import { useWails } from "../../wails";
import { formatColumns, FormattedQueryResult, RowAction } from "./table";

const wails = useWails();
const router = useRouter();
const { databaseId } = useUrlParams();

const tabs = [
  {
    label: "Schemas",
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
    wails(() => GetDatabaseSchemas(databaseId.value, page, itemsPerPage)).pipe(
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
    wails(() => GetDatabaseInfo(databaseId.value, page, itemsPerPage)).pipe(
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

function navigateToSchema(schemaId: string) {
  router.push({ name: "schema", params: { schemaId } });
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
            navigateToSchema(
              row.original.SCHEMA_NAME ||
                row.original.schema_name ||
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
