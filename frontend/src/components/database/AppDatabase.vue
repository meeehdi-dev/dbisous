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
import { FormattedQueryResult, RowAction } from "./table";
import { formatQueryResult } from "../../effects/columns";

const wails = useWails();
const router = useRouter();
const { databaseId } = useUrlParams();

const data = ref<FormattedQueryResult>();
const fetchingData = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  await Effect.runPromise(
    wails(() => GetDatabaseSchemas(databaseId.value, page, itemsPerPage)).pipe(
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
    wails(() => GetDatabaseInfo(databaseId.value, page, itemsPerPage)).pipe(
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

function navigateToSchema(schemaId: string) {
  router.push({ name: "schema", params: { schemaId } });
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
  </AppTabs>
</template>
