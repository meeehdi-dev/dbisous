<script setup lang="ts">
import { useRouter } from "vue-router";
import { GetDatabaseSchemas } from "../../../wailsjs/go/app/App";
import { useUrlParams } from "../../composables/useUrlParams";
import { FormattedQueryResult, RowAction } from "./table/table";
import { useWails } from "../../wails";
import { formatQueryResult } from "../../effects/columns";
import { ref } from "vue";
import { Effect } from "effect";
import { client } from "../../../wailsjs/go/models";

const router = useRouter();
const { databaseId } = useUrlParams();

function navigateToSchema(schemaId: string) {
  router.push({ name: "schema", params: { schemaId } });
}

const wails = useWails();

const data = ref<FormattedQueryResult>();
const columns = ref<client.ColumnMetadata[]>();
const fetchingData = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  await Effect.runPromise(
    wails(() => GetDatabaseSchemas(databaseId.value, page, itemsPerPage)).pipe(
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
      <AppColumns :loading="fetchingData" :data="columns" />
    </template>
  </AppTabs>
</template>
