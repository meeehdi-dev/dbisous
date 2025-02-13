<script setup lang="ts">
import { TableColumn, TableData } from "@nuxt/ui/dist/runtime/types";
import { useRouter } from "vue-router";
import { ref } from "vue";
import { client } from "../../../wailsjs/go/models";
import { useUrlParams } from "../../composables/useUrlParams";
import { useWails } from "../../wails";
import { GetSchemaInfo, GetTables } from "../../../wailsjs/go/app/App";
import { Effect } from "effect";
import { cell } from "./cell";

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

const dataKey = ref(0);
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
      dataKey.value++;
    }),
  ),
);

const infoKey = ref(0);
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
      infoKey.value++;
    }),
  ),
);

function redirectToTable(schemaId: string, tableId: string) {
  router.push({
    name: "table",
    params: { schemaId, tableId },
  });
}

const columnPinning = ref({ right: ["action"] });
</script>

<template>
  <UTabs :items="tabs" variant="link" :ui="{ content: 'flex flex-col gap-2' }">
    <template #data>
      <UTable
        :data="data.rows"
        :columns="data.columns"
        v-model:column-pinning="columnPinning"
        :key="dataKey"
      >
        <template #action-cell="{ row }">
          <UButton
            icon="lucide:eye"
            variant="ghost"
            @click="
              redirectToTable(
                row.original.table_schema || row.original.schema,
                row.original.table_name || row.original.name,
              )
            "
          />
        </template>
      </UTable>
    </template>
    <template #info>
      <UTable
        :data="info.rows"
        :columns="info.columns"
        v-model:column-pinning="columnPinning"
        :key="infoKey"
      />
    </template>
    <template #script>
      <AppScript />
    </template>
  </UTabs>
</template>
