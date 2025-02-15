<script setup lang="ts">
import { useRouter } from "vue-router";
import { ref } from "vue";
import { Effect } from "effect";
import { GetSchemas } from "../../../wailsjs/go/app/App";
import { useUrlParams } from "../../composables/useUrlParams";
import { useWails } from "../../wails";
import { formatColumns, FormattedQueryResult, RowAction } from "./table";

const wails = useWails();
const router = useRouter();
const { databaseId } = useUrlParams();

const tabs = [
  {
    label: "Schemas",
    slot: "data",
    icon: "lucide:list",
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
        :data="data"
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
      <AppRows :data="info" />
    </template>
    <template #script>
      <AppScript />
    </template>
  </UTabs>
</template>
