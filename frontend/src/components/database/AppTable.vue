<script setup lang="ts">
import { ref } from "vue";
import { useUrlParams } from "../../composables/useUrlParams";
import { Effect } from "effect";
import { FormattedQueryResult } from "./table";
import { useWails } from "../../wails";
import { GetTableInfo, GetTableRows } from "../../../wailsjs/go/app/App";
import { formatQueryResult } from "../../effects/columns";

const wails = useWails();
const { databaseId, schemaId, tableId } = useUrlParams();

const data = ref<FormattedQueryResult>();
const fetchingData = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  await Effect.runPromise(
    wails(() =>
      GetTableRows(
        databaseId.value,
        page,
        itemsPerPage,
        schemaId.value,
        tableId.value,
      ),
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
      GetTableInfo(
        databaseId.value,
        page,
        itemsPerPage,
        schemaId.value,
        tableId.value,
      ),
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
</script>

<template>
  <AppTabs :default-query="`SELECT * FROM ${tableId};`">
    <template #data>
      <AppRows
        :loading="fetchingData"
        :data="data"
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
