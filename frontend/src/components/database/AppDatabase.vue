<script setup lang="ts">
import { useRouter } from "vue-router";
import type { TableColumn, TableData } from "@nuxt/ui/dist/module";
import { ref } from "vue";
import { Effect } from "effect";
import { GetDatabaseInfo, GetSchemas } from "../../../wailsjs/go/app/App";
import { useUrlParams } from "../../composables/useUrlParams";
import { useWails } from "../../wails";
import { client } from "../../../wailsjs/go/models";
import { cell } from "./cell";

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
  wails(() => GetSchemas(databaseId.value)).pipe(
    Effect.andThen((result) => ({
      ...result,
      columns: result.columns
        .map((column) => ({
          accessorKey: column.name,
          header: column.name,
          cell: cell(""),
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
  wails(() => GetDatabaseInfo(databaseId.value)).pipe(
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

function navigateToSchema(schemaId: string) {
  router.push({ name: "schema", params: { schemaId } });
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
              navigateToSchema(row.original.schema_name || row.original.name)
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
