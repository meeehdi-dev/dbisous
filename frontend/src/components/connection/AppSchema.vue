<script setup lang="ts">
import { useRouter } from "vue-router";
import {
  formatColumns,
  FormattedQueryResult,
  RowAction,
} from "@/components/connection/table/table";
import { useWails } from "@/composables/useWails";
import { GetSchemaTables } from "_/go/app/App";
import { ref } from "vue";
import { client } from "_/go/models";
import { SortDirection } from "@/components/connection/table/column/AppColumnHeader.vue";
import { useApp } from "@/composables/shared/useApp";
import { Route } from "@/router";
import { toSqlValue } from "@/utils/transaction";

const wails = useWails();
const router = useRouter();
const { database, schema, table } = useApp();
table.value = ""; // FIXME: reset table var bc breadcrumb does not provide onclick

async function navigateToTable(t: string) {
  table.value = t;
  await router.push({
    name: Route.Table,
  });
}

const rows = ref<FormattedQueryResult & { key: number }>();
const rowsKey = ref(0);
const sorting = ref<Array<{ id: string; desc: boolean }>>([]);
const filtering = ref<Array<{ id: string; value: unknown }>>([]);
const columns = ref<Array<client.ColumnMetadata>>();
const loading = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  loading.value = true;
  const result = await wails(() =>
    GetSchemaTables(
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
    ),
  );
  loading.value = false;
  if (result instanceof Error) {
    return;
  }
  columns.value = result.columns;
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
      undefined,
      undefined,
      true,
    ),
  };
}
await fetchData();
</script>

<template>
  <AppTabs>
    <template #rows>
      <AppRows
        :loading="loading"
        :data="rows"
        :sorting="sorting"
        :filtering="filtering"
        :actions="[RowAction.View]"
        @view="
          (row) =>
            navigateToTable(
              (row.TABLE_NAME || row.table_name || row.name) as string,
            )
        "
        @pagination-change="fetchData"
      />
    </template>
    <template #columns>
      <AppColumns :loading="loading" :data="columns" />
    </template>
  </AppTabs>
</template>
