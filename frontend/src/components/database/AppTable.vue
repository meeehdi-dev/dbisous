<script setup lang="ts">
import { TableColumn, TableData } from "@nuxt/ui/dist/runtime/types";
// import { Effect } from "effect";
// import { useWails } from "../../wails";
// import { ExecuteQuery, GetColumns, GetRows } from "../../../wailsjs/go/app/App";
import { ref } from "vue";
import { useUrlParams } from "../../composables/useUrlParams";
import { client } from "../../../wailsjs/go/models";
import { Effect } from "effect";
import { cell } from "./cell";
import { useWails } from "../../wails";
import { GetTableRows } from "../../../wailsjs/go/app/App";
// import { cell } from "./cell";

const wails = useWails();
const { databaseId, schemaId, tableId } = useUrlParams();

const query = ref(`SELECT * FROM ${tableId.value}`);
const columnsKey = ref(0);
const dataKey = ref(0);
const scriptKey = ref(0);
const open = ref(false);

const tabs = [
  {
    label: "Data",
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

const tableColumns = ref<{
  columns: Array<TableColumn<TableData>>;
  rows: Array<TableColumn<TableData>>;
  duration: string;
}>({ columns: [], rows: [], duration: "" });
// await Effect.runPromise(
//   wails(() =>
//     GetColumns(currentDatabase.value, currentSchema.value, currentTable.value),
//   ).pipe(
//     Effect.andThen((data) => ({
//       ...data,
//       columns: data.columns.map((column) => ({
//         accessorKey: column.name,
//         header: column.name,
//         cell: cell(mapColumnType(column.name)),
//       })),
//     })),
//     Effect.tap((data) => {
//       tableColumns.value = data;
//       columnsKey.value++;
//     }),
//   ),
// );

const rowsKey = ref(0);
const rows = ref<
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
  wails(() =>
    GetTableRows(databaseId.value, schemaId.value, tableId.value),
  ).pipe(
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
      rows.value = data;
      rowsKey.value++;
    }),
  ),
);

const queryData = ref<{
  columns: Array<TableColumn<TableData>>;
  rows: Array<TableColumn<TableData>>;
  duration: string;
}>({ columns: [], rows: [], duration: "" });
async function executeQuery() {
  // await Effect.runPromise(
  //   wails(() => ExecuteQuery(currentDatabase.value, query.value, [])).pipe(
  //     Effect.andThen((data) => ({
  //       ...data,
  //       columns: data.columns
  //         .map((column) => ({
  //           accessorKey: column.name,
  //           header: column.name,
  //           cell: cell(column.type),
  //         }))
  //         .concat([
  //           // @ts-expect-error tkt
  //           {
  //             accessorKey: "action",
  //             header: "Actions",
  //           },
  //         ]),
  //     })),
  //     Effect.tap((data) => {
  //       queryData.value = data;
  //       scriptKey.value++;
  //     }),
  //   ),
  // );
}

// function addRow() {
//   const obj: Record<string, unknown> = {};
//   for (const column of tableColumns.value.rows) {
//     obj[column.name] = column.dflt_value;
//   }
//   tableRows.value.rows.push(obj);
//   dataKey.value++;
// }

const columnPinning = ref({ right: ["action"] });
</script>

<template>
  <div class="flex flex-col justify-between h-full">
    <UTabs
      :items="tabs"
      variant="link"
      :ui="{ content: 'flex flex-col gap-2' }"
    >
      <template #data>
        <UTable
          :data="rows.rows"
          :columns="rows.columns"
          :key="dataKey"
          v-model:column-pinning="columnPinning"
        >
          <template #action-cell="{ row }">
            <UButton icon="lucide:copy" color="info" variant="ghost" />
            <UButton icon="lucide:trash" color="error" variant="ghost" />
          </template>
        </UTable>
        <div class="flex justify-center">
          <UButton icon="lucide:plus" variant="soft" label="Add row" />
        </div>
      </template>
      <template #info>
        <UTable
          :data="tableColumns.rows"
          :columns="tableColumns.columns"
          :key="columnsKey"
          v-model:column-pinning="columnPinning"
        />
      </template>
      <template #script>
        <div class="p-4 flex flex-col gap-4 w-full">
          <AppEditor v-model="query" :columns="tableColumns.rows" />
          <div class="flex gap-2 items-center">
            <UButton
              icon="lucide:terminal"
              label="Execute"
              @click="executeQuery"
            />
            <span v-if="queryData.duration" class="text-sm text-neutral-400">{{
              queryData.duration
            }}</span>
          </div>
          <USeparator
            :label="`${queryData.rows.length.toString()} result${queryData.rows.length > 1 ? 's' : ''}`"
          />
          <UTable
            :data="queryData.rows"
            :columns="queryData.columns"
            :key="scriptKey"
            v-model:column-pinning="columnPinning"
          >
            <template #action-cell="{ row }">
              <UButton icon="lucide:copy" color="info" variant="ghost" />
              <UButton icon="lucide:trash" color="error" variant="ghost" />
            </template>
          </UTable>
        </div>
      </template>
    </UTabs>
    <div class="p-2">
      <UAlert
        title="3 pending changes"
        icon="lucide:info"
        color="info"
        variant="soft"
        :actions="[
          {
            label: 'Apply',
            color: 'warning',
            size: 'md',
            onClick() {
              open = true;
            },
          },
          {
            label: 'Cancel',
            color: 'neutral',
            size: 'md',
          },
        ]"
      />
      <UModal
        v-model:open="open"
        title="Are you sure?"
        description="Please check the SQL script before applying."
      >
        <template #body>
          <div class="flex flex-col gap-8">
            <AppEditor v-model="query" :columns="tableColumns.rows" />
            <div class="flex justify-end gap-2">
              <UButton label="Apply" />
              <UButton color="neutral" label="Cancel" @click="open = false" />
            </div>
          </div>
        </template>
      </UModal>
    </div>
  </div>
</template>
