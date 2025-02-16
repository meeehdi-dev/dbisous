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
const fetchingData = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  await Effect.runPromise(
    wails(() => GetDatabaseSchemas(databaseId.value, page, itemsPerPage)).pipe(
      Effect.tap((result) => {
        data.value = {
          ...result,
          columns: formatColumns(result.columns),
        };
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
    wails(() => GetDatabaseInfo(databaseId.value, page, itemsPerPage)).pipe(
      Effect.tap((result) => {
        info.value = {
          ...result,
          columns: formatColumns(result.columns, false),
        };
        fetchingInfo.value = false;
      }),
      Effect.catchTags({
        WailsError: Effect.succeed,
      }),
    ),
  );
}
fetchInfo();

function navigateToSchema(schemaId: string) {
  router.push({ name: "schema", params: { schemaId } });
}
</script>

<template>
  <UTabs
    :items="tabs"
    variant="link"
    :ui="{ root: 'flex flex-auto', content: 'flex flex-auto flex-col gap-2' }"
  >
    <template #data>
      <AppRows
        :loading="fetchingData"
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
    <template #script>
      <AppScript />
    </template>
  </UTabs>
</template>
