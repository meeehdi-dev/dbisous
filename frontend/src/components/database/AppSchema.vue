<script setup lang="ts">
import { useRouter } from "vue-router";
import { ref } from "vue";
import { useUrlParams } from "../../composables/useUrlParams";
import { useWails } from "../../wails";
import { GetTables } from "../../../wailsjs/go/app/App";
import { Effect } from "effect";
import { formatColumns, FormattedQueryResult, RowAction } from "./table";

const wails = useWails();
const router = useRouter();
const { databaseId, schemaId } = useUrlParams();

const tabs = [
  {
    label: "Tables",
    slot: "data",
    icon: "lucide:list-ordered",
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
await Effect.runPromise(
  wails(() => GetTables(databaseId.value, schemaId.value)).pipe(
    Effect.tap((result) => {
      data.value = {
        ...result.data,
        columns: formatColumns(result.data.columns),
      };
      info.value = {
        ...result.info,
        columns: formatColumns(result.info.columns, false),
      };
    }),
    Effect.catchTags({
      WailsError: Effect.succeed,
    }),
  ),
);

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
        :rows="data?.rows"
        :columns="data?.columns"
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
      <AppRows :rows="info?.rows" :columns="info?.columns" />
    </template>
    <template #script>
      <AppScript />
    </template>
  </UTabs>
</template>
