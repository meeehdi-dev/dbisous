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
const tx = useTransaction();

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
  fetchingData.value = false;
  if (result instanceof Error) {
    return;
  }
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
  // TODO: push tx insert changes
}
fetchData();

function insertRow() {
  const row: Record<string, unknown> = {};
  columns.value?.forEach((c) => {
    row[c.name] = c.default_value;
  });
  const key = tx.addInsert(tableId.value, row);
  row.__key = key;
  data.value!.rows.push(row);
  data.value!.key++;
}

function duplicateRow(row: Record<string, unknown>) {
  if (!primaryKey.value) {
    return;
  }

  const dup = { ...row, [primaryKey.value]: "NULL" };
  const key = tx.addInsert(tableId.value, dup);
  dup.__key = key;
  data.value!.rows.push(dup);
  data.value!.key++;
}

function deleteRow(row: Record<string, unknown>) {
  let rowKey = row.__key;
  if (rowKey === undefined && primaryKey.value) {
    rowKey = row[primaryKey.value];
    tx.toggleDelete(tableId.value, primaryKey.value!, rowKey);
  } else if (rowKey !== undefined) {
    tx.removeInsert(tableId.value, rowKey as number);
    data.value!.rows.splice(
      data.value!.rows.findIndex(
        (r: Record<string, unknown>) => r.__key === rowKey,
      ),
      1,
    );
    data.value!.key++;
  }
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
          primaryKey
            ? [RowAction.Insert, RowAction.Duplicate, RowAction.Delete]
            : undefined
        "
        @insert="insertRow"
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
