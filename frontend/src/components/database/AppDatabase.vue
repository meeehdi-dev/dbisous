<script setup lang="ts">
import { useRouter } from "vue-router";
import type { TableColumn, TableData } from "@nuxt/ui/dist/module";
import { ref } from "vue";
import { Effect } from "effect";
import { GetSchemas } from "../../../wailsjs/go/app/App";
import { useUrlParams } from "../../composables/useUrlParams";
import { useWails } from "../../wails";
import { client } from "../../../wailsjs/go/models";
import { formatColumns, RowAction } from "./table";

const wails = useWails();
const router = useRouter();
const { databaseId } = useUrlParams();

const tabs = [
  {
    label: "Schemas",
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

const data = ref<
  Omit<client.QueryResult, "convertValues" | "columns"> & {
    columns: Array<TableColumn<TableData>>;
  }
>({
  columns: [],
  rows: [],
  sql_duration: "",
  total_duration: "",
});
const info = ref<
  Omit<client.QueryResult, "convertValues" | "columns"> & {
    columns: Array<TableColumn<TableData>>;
  }
>({
  columns: [],
  rows: [],
  sql_duration: "",
  total_duration: "",
});
await Effect.runPromise(
  wails(() => GetSchemas(databaseId.value)).pipe(
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
        :rows="data.rows"
        :columns="data.columns"
        :actions="[RowAction.View]"
        @view="
          (row) =>
            navigateToSchema(
              row.original.SCHEMA_NAME ||
                row.original.schema_name ||
                row.original.name,
            )
        "
      />
    </template>
    <template #info>
      <AppRows :rows="info.rows" :columns="info.columns" />
    </template>
    <template #script>
      <AppScript />
    </template>
  </UTabs>
</template>
