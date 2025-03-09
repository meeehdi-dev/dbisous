<script setup lang="ts">
import { GetTableRows } from "_/go/app/App";
import { ref } from "vue";
import {
  formatColumns,
  FormattedQueryResult,
  RowAction,
} from "@/components/database/table/table";
import { client } from "_/go/models";
import { useWails } from "@/composables/useWails";
import { toSqlValue, useTransaction } from "@/composables/useTransaction";
import { SortDirection } from "@/components/database/table/column/AppColumnHeader.vue";
import { useApp } from "@/composables/useApp";

const { database, schema, table } = useApp();

const wails = useWails();
const tx = useTransaction();

const rows = ref<FormattedQueryResult & { key: number }>();
const rowsKey = ref(0);
const sorting = ref<Array<{ id: string; desc: boolean }>>([]);
const filtering = ref<Array<{ id: string; value: unknown }>>([]);
const columns = ref<Array<client.ColumnMetadata>>();
const primaryKey = ref<string>();
const fetchingData = ref(false);

async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  const result = await wails(() =>
    GetTableRows(
      database.value,
      new client.QueryParams({
        offset: (page - 1) * itemsPerPage,
        limit: itemsPerPage,
        filter: filtering.value.map((s) => ({
          column: s.id,
          value: toSqlValue(s.value),
        })),
        order: sorting.value.map((s) => ({
          column: s.id,
          direction: s.desc
            ? client.OrderDirection.Descending
            : client.OrderDirection.Ascending,
        })),
      }),
      schema.value,
      table.value,
    ),
  );
  fetchingData.value = false;
  if (result instanceof Error) {
    return;
  }
  columns.value = result.columns;
  primaryKey.value = result.columns.find((c) => c.primary_key)?.name;
  rows.value = {
    key: rowsKey.value++,
    // eslint-disable-next-line @typescript-eslint/no-misused-spread
    ...result,
    columns: formatColumns(
      result.columns,
      async (name: string, s: SortDirection) => {
        if (!s) {
          sorting.value = [];
        } else {
          sorting.value = [
            { id: name, desc: s === client.OrderDirection.Descending },
          ];
        }
        return fetchData();
      },
      async (name: string, f: unknown) => {
        if (!f) {
          filtering.value = [];
        } else {
          filtering.value = [{ id: name, value: f }];
        }
        return fetchData();
      },
      table.value,
      primaryKey.value,
      false,
    ),
  };
  // TODO: push tx insert changes
}
await fetchData();

function insertRow() {
  if (!rows.value) {
    return;
  }

  const row: Record<string, unknown> = {};
  columns.value?.forEach((c) => {
    row[c.name] = c.default_value;
  });
  const key = tx.addInsert(table.value, row);
  row.__key = key;

  rows.value.rows.push(row);
  rows.value.key++;
}

function duplicateRow(row: Record<string, unknown>) {
  if (!primaryKey.value || !rows.value) {
    return;
  }

  const dup = { ...row, [primaryKey.value]: "NULL" };
  const key = tx.addInsert(table.value, dup);
  dup.__key = key;

  rows.value.rows.push(dup);
  rows.value.key++;
}

function deleteRow(row: Record<string, unknown>) {
  if (!rows.value) {
    return;
  }

  let rowKey = row.__key;
  if (rowKey === undefined && primaryKey.value) {
    rowKey = row[primaryKey.value];
    tx.toggleDelete(table.value, primaryKey.value, rowKey);
  } else if (rowKey !== undefined) {
    tx.removeInsert(table.value, rowKey as number);

    rows.value.rows.splice(
      rows.value.rows.findIndex(
        (r: Record<string, unknown>) => r.__key === rowKey,
      ),
      1,
    );
    rows.value.key++;
  }
}
</script>

<template>
  <AppTabs :default-query="`SELECT * FROM ${table};`">
    <template #rows>
      <AppRows
        :loading="fetchingData"
        :data="rows"
        :sorting="sorting"
        :filtering="filtering"
        :table="table"
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
    <template #columns>
      <AppColumns
        :loading="fetchingData"
        :data="columns"
        :table="table"
        :primary-key="primaryKey"
      />
    </template>
  </AppTabs>
</template>
