<script setup lang="ts">
import { useUrlParams } from "@/composables/useUrlParams";
import { GetTableRows } from "_/go/app/App";
import { ref } from "vue";
import {
  FormattedQueryResult,
  RowAction,
} from "@/components/database/table/table";
import { client } from "_/go/models";
import { Effect } from "effect";
import { useWails } from "@/composables/useWails";
import { formatQueryResult } from "@/effects/columns";
import { useTransaction } from "@/composables/useTransaction";

const { databaseId, schemaId, tableId } = useUrlParams();
const wails = useWails();

const data = ref<FormattedQueryResult>();
const columns = ref<client.ColumnMetadata[]>();
const primaryKey = ref("");
const fetchingData = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  await Effect.runPromise(
    wails(() =>
      GetTableRows(
        databaseId.value,
        page,
        itemsPerPage,
        schemaId.value,
        tableId.value,
      ),
    ).pipe(
      Effect.tap((result) => {
        columns.value = result.columns;
        primaryKey.value = result.primary_key;
      }),
      Effect.andThen(formatQueryResult),
      Effect.tap((result) => {
        data.value = result;
        fetchingData.value = false;
      }),
      Effect.catchTags({
        WailsError: Effect.succeed,
      }),
    ),
  );
}
fetchData();

const tx = useTransaction();

function duplicateRow(row: unknown) {
  // @ts-expect-error tkt
  const dup = { ...row };
  delete dup[primaryKey.value];
  const key = tx.addInsert(tableId.value, dup);
  dup.__key = key;
  data.value!.rows.push(dup);
}

function deleteRow(row: unknown) {
  // @ts-expect-error tkt
  const rowKey = row[primaryKey.value] as unknown | undefined;
  if (rowKey === undefined) {
    // @ts-expect-error tkt
    const key = row.__key as number;
    if (key !== undefined) {
      tx.removeInsert(tableId.value, key);
      // @ts-expect-error tkt
      data.value.rows = data.value.rows.filter((row) => row.__key !== key);
      return;
    }
  }
  tx.toggleDelete(tableId.value, primaryKey.value, rowKey);
}
</script>

<template>
  <AppTabs :default-query="`SELECT * FROM ${tableId};`">
    <template #data>
      <AppRows
        :loading="fetchingData"
        :data="data"
        :actions="[RowAction.Duplicate, RowAction.Delete]"
        @duplicate="duplicateRow"
        @delete="deleteRow"
        @pagination-change="fetchData"
      />
    </template>
    <template #info>
      <AppColumns :loading="fetchingData" :data="columns" />
    </template>
  </AppTabs>
</template>
