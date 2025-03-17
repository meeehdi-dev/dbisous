<script setup lang="ts">
import { GetTableRows } from "_/go/app/App";
import { computed, ref, watch } from "vue";
import {
  formatColumns,
  FormattedQueryResult,
  RowAction,
} from "@/components/connection/table/table";
import { client } from "_/go/models";
import { useWails } from "@/composables/useWails";
import { useTransaction } from "@/composables/shared/useTransaction";
import { SortDirection } from "@/components/connection/table/column/AppColumnHeader.vue";
import { useApp } from "@/composables/shared/useApp";
import { toSqlValue } from "@/utils/transaction";
import { Tab } from "@/utils/tabs";

const { connection, schema, table } = useApp();

const wails = useWails();
const tx = useTransaction();

const rows = ref<FormattedQueryResult & { key: number }>();
const query = ref<string>();
const rowsKey = ref(0);
const sorting = ref<Array<{ id: string; desc: boolean }>>([]);
const filtering = ref<Array<{ id: string; value: unknown }>>([]);
const columns = ref<Array<client.ColumnMetadata>>();
const primaryKey = ref<string>();
const loading = ref(false);
const active = ref(Tab.Rows);
const defaultQuery = ref<string>();
watch(active, () => {
  if (active.value !== Tab.Query) {
    defaultQuery.value = undefined;
  }
});

const tableQuery = computed(() => `SELECT * FROM ${table.value};`);

async function fetchData(page = 1, itemsPerPage = 10) {
  loading.value = true;
  const result = await wails(() =>
    GetTableRows(
      connection.value,
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
  loading.value = false;
  if (result instanceof Error) {
    return;
  }
  query.value = result.query;
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
          sorting.value = sorting.value.filter((s) => s.id !== name);
        } else {
          sorting.value = [
            ...sorting.value.filter((s) => s.id !== name),
            { id: name, desc: s === client.OrderDirection.Descending },
          ];
        }
        return fetchData();
      },
      async (name: string, f: unknown) => {
        if (!f) {
          filtering.value = filtering.value.filter((f) => f.id !== name);
        } else {
          filtering.value = [
            ...filtering.value.filter((f) => f.id !== name),
            { id: name, value: f },
          ];
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

function onQueryEdit(query: string) {
  defaultQuery.value = query;
  active.value = Tab.Query;
}
</script>

<template>
  <AppTabs v-model="active" :default-query="defaultQuery ?? tableQuery">
    <template #rows>
      <AppRows
        :loading="loading"
        :query="query"
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
        @query-edit="onQueryEdit"
      />
    </template>
    <template #columns>
      <AppColumns
        :loading="loading"
        :data="columns"
        :table="table"
        :primary-key="primaryKey"
      />
    </template>
  </AppTabs>
</template>
