<script setup lang="ts">
import type { TableColumn, TableData } from "@nuxt/ui/dist/runtime/types";
import { useRouter } from "vue-router";
import { ref } from "vue";
import { client } from "../../../wailsjs/go/models";
import { useUrlParams } from "../../composables/useUrlParams";
import { useWails } from "../../wails";
import { GetSchemaInfo, GetTables } from "../../../wailsjs/go/app/App";
import { Effect } from "effect";
import { cell } from "./cell";
import { RowAction } from "./row";

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
await Effect.runPromise(
  wails(() => GetTables(databaseId.value, schemaId.value)).pipe(
    Effect.andThen((result) => ({
      ...result,
      columns: result.columns
        .map((column) => ({
          accessorKey: column.name,
          header: column.name,
          cell: cell(column.type),
        }))
        .concat([
          // @ts-expect-error tkt
          {
            accessorKey: "action",
            header: "Actions",
          },
        ]),
    })),
    Effect.tap((result) => {
      data.value = result;
    }),
  ),
);

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
  wails(() => GetSchemaInfo(databaseId.value, schemaId.value)).pipe(
    Effect.andThen((result) => ({
      ...result,
      columns: result.columns.map((column) => ({
        accessorKey: column.name,
        header: column.name,
        cell: cell(""),
      })),
    })),
    Effect.tap((result) => {
      info.value = result;
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
  <UTabs :items="tabs" variant="link" :ui="{ content: 'flex flex-col gap-2' }">
    <template #data>
      <AppRows
        :rows="data.rows"
        :columns="data.columns"
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
      <AppRows :rows="info.rows" :columns="info.columns" />
    </template>
    <template #script>
      <AppScript />
    </template>
  </UTabs>
</template>
