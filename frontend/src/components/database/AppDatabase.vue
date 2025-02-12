<script setup lang="ts">
import { useRouter } from "vue-router";
import type { TableColumn, TableData } from "@nuxt/ui/dist/module";
import { ref } from "vue";
import { Effect } from "effect";
import { GetSchemas } from "../../../wailsjs/go/app/App";
import { useUrlParams } from "../../composables/useUrlParams";
import { useWails } from "../../wails";
import { client } from "../../../wailsjs/go/models";
import { cell } from "./cell";

const wails = useWails();
const router = useRouter();
const { databaseId } = useUrlParams();

const schemasKey = ref(0);
const schemas = ref<
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
    Effect.tap(Effect.log),
    Effect.andThen((data) => ({
      ...data,
      columns: data.columns
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
    Effect.tap((data) => {
      schemas.value = data;
      schemasKey.value++;
    }),
  ),
);

function navigateToSchema(schemaId: string) {
  router.push({ name: "schema", params: { schemaId } });
}

const columnPinning = ref({ right: ["action"] });
</script>

<template>
  <UTable
    :data="schemas.rows"
    :columns="schemas.columns"
    v-model:column-pinning="columnPinning"
    :key="schemasKey"
  >
    <template #action-cell="{ row }">
      <UButton
        icon="lucide:eye"
        variant="ghost"
        @click="navigateToSchema(row.original.name)"
      />
    </template>
  </UTable>
</template>
