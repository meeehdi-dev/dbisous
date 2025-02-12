<script setup lang="ts">
import { TableColumn, TableData } from "@nuxt/ui/dist/runtime/types";
import { useRouter } from "vue-router";
import { ref } from "vue";
import { client } from "../../../wailsjs/go/models";
import { useUrlParams } from "../../composables/useUrlParams";
import { useWails } from "../../wails";
import { GetTables } from "../../../wailsjs/go/app/App";
import { Effect } from "effect";
import { cell } from "./cell";

const wails = useWails();
const router = useRouter();
const { databaseId, schemaId } = useUrlParams();

const tablesKey = ref(0);
const tables = ref<
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
    Effect.tap(Effect.log),
    Effect.andThen((data) => ({
      ...data,
      columns: data.columns
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
    Effect.tap((data) => {
      tables.value = data;
      tablesKey.value++;
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
  <UTable
    :data="tables.rows"
    :columns="tables.columns"
    v-model:column-pinning="columnPinning"
  >
    <template #action-cell="{ row }">
      <UButton
        icon="lucide:eye"
        variant="ghost"
        @click="redirectToTable(row.original.schema, row.original.name)"
      />
    </template>
  </UTable>
</template>
