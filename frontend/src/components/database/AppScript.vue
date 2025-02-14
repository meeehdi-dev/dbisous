<script setup lang="ts">
import type { TableColumn, TableData } from "@nuxt/ui/dist/runtime/types";
import { ref } from "vue";
import { useUrlParams } from "../../composables/useUrlParams";
import { client } from "../../../wailsjs/go/models";
import { Effect } from "effect";
import { cell } from "./cell";
import { useWails } from "../../wails";
import { ExecuteQuery } from "../../../wailsjs/go/app/App";

const { defaultQuery } = defineProps<{ defaultQuery?: string }>();

const wails = useWails();
const { databaseId } = useUrlParams();

const query = ref(defaultQuery ?? "");

const queryResult = ref<
  Omit<client.QueryResult, "convertValues" | "columns"> & {
    columns: Array<TableColumn<TableData>>;
  }
>({
  columns: [],
  rows: [],
  sql_duration: "",
  total_duration: "",
});
async function executeQuery() {
  await Effect.runPromise(
    wails(() => ExecuteQuery(databaseId.value, query.value)).pipe(
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
        queryResult.value = data;
      }),
    ),
  );
}
</script>

<template>
  <div class="flex flex-1 flex-col gap-4 w-full">
    <div class="flex flex-col p-4 gap-4">
      <AppEditor v-model="query" />
      <div class="flex gap-2 items-center">
        <UButton icon="lucide:terminal" label="Execute" @click="executeQuery" />
        <span
          v-if="queryResult.sql_duration"
          class="text-sm text-neutral-400"
          >{{ queryResult.sql_duration }}</span
        >
      </div>
    </div>
    <USeparator
      :label="`${queryResult.rows.length.toString()} result${queryResult.rows.length > 1 ? 's' : ''}`"
    />
    <AppRows :rows="queryResult.rows" :columns="queryResult.columns" />
  </div>
</template>
