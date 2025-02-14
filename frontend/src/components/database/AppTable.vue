<script setup lang="ts">
import type { TableColumn, TableData } from "@nuxt/ui/dist/runtime/types";
import { ref } from "vue";
import { useUrlParams } from "../../composables/useUrlParams";
import { client } from "../../../wailsjs/go/models";
import { Effect } from "effect";
import { cell } from "./cell";
import { useWails } from "../../wails";
import { GetTableInfo, GetTableRows } from "../../../wailsjs/go/app/App";

const wails = useWails();
const { databaseId, schemaId, tableId } = useUrlParams();

const transactionQuery = ref(`SELECT * FROM ${tableId.value}`);
const open = ref(false);

const tabs = [
  {
    label: "Rows",
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
  wails(() =>
    GetTableRows(databaseId.value, schemaId.value, tableId.value),
  ).pipe(
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
  wails(() =>
    GetTableInfo(databaseId.value, schemaId.value, tableId.value),
  ).pipe(
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
</script>

<template>
  <div class="flex flex-col justify-between h-full">
    <UTabs
      :items="tabs"
      variant="link"
      :ui="{ content: 'flex flex-col gap-2' }"
    >
      <template #data>
        <AppRows :rows="data.rows" :columns="data.columns" />
      </template>
      <template #info>
        <AppRows :rows="info.rows" :columns="info.columns" />
      </template>
      <template #script>
        <AppScript :default-query="`SELECT * FROM ${tableId};`" />
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
            icon: 'lucide:check',
            color: 'warning',
            size: 'md',
            onClick() {
              open = true;
            },
          },
          {
            label: 'Cancel',
            icon: 'lucide:x',
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
            <AppEditor v-model="transactionQuery" />
            <div class="flex justify-end gap-2">
              <UButton icon="lucide:check" label="Apply" />
              <UButton
                icon="lucide:x"
                color="neutral"
                label="Cancel"
                @click="open = false"
              />
            </div>
          </div>
        </template>
      </UModal>
    </div>
  </div>
</template>
