<script setup lang="ts">
import { useUrlParams } from "@/composables/useUrlParams";
import { GetTableRows } from "_/go/app/App";
import { ref } from "vue";
import {
  formatColumns,
  FormattedQueryResult,
  RowAction,
} from "@/components/database/table/table";
import { client } from "_/go/models";
import { useWails } from "@/composables/useWails";
import { useTransaction } from "@/composables/useTransaction";

const { databaseId, schemaId, tableId } = useUrlParams();
const wails = useWails();

const data = ref<FormattedQueryResult & { key: number }>();
const dataKey = ref(0);
const columns = ref<client.ColumnMetadata[]>();
const primaryKey = ref<string>();
const fetchingData = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  const result = await wails(() =>
    GetTableRows(
      databaseId.value,
      page,
      itemsPerPage,
      schemaId.value,
      tableId.value,
    ),
  );
  if (result instanceof Error) {
    // TODO: specific error handling
  } else {
    columns.value = result.columns;
    primaryKey.value = result.columns.find((c) => c.primary_key)?.name;
    data.value = {
      key: dataKey.value++,
      ...result,
      columns: formatColumns(
        result.columns,
        tableId.value,
        primaryKey.value,
        false,
      ),
    };
  }
  fetchingData.value = false;
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
  data.value!.key++;
}

function deleteRow(row: unknown) {
  // @ts-expect-error tkt
  const rowKey = row[primaryKey.value] as unknown | undefined;
  if (rowKey === undefined) {
    // @ts-expect-error tkt
    const key = row.__key as number;
    if (key !== undefined) {
      tx.removeInsert(tableId.value, key);
      data.value!.rows.splice(
        // @ts-expect-error tkt
        data.value!.rows.findIndex((r: unknown) => r.__key === key),
        1,
      );
      data.value!.key++;
      return;
    }
  }
  tx.toggleDelete(tableId.value, primaryKey.value!, rowKey);
}
</script>

<template>
  <AppTabs :default-query="`SELECT * FROM ${tableId};`">
    <template #data>
      <AppRows
        :loading="fetchingData"
        :data="data"
        :table="tableId"
        :primary-key="primaryKey"
        :actions="
          primaryKey ? [RowAction.Duplicate, RowAction.Delete] : undefined
        "
        @duplicate="duplicateRow"
        @delete="deleteRow"
        @pagination-change="fetchData"
      />
    </template>
    <template #info>
      <AppColumns
        :loading="fetchingData"
        :data="columns"
        :table="tableId"
        :primary-key="primaryKey"
      />
    </template>
  </AppTabs>
</template>
