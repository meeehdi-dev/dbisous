<script setup lang="ts">
import { useRouter } from "vue-router";
import { GetDatabaseSchemas } from "_/go/app/App";
import {
  formatColumns,
  FormattedQueryResult,
  RowAction,
} from "@/components/connection/table/table";
import { useWails } from "@/composables/useWails";
import { ref, watch } from "vue";
import { client } from "_/go/models";
import { SortDirection } from "@/components/connection/table/column/AppColumnHeader.vue";
import { Route } from "@/router";
import { useApp } from "@/composables/shared/useApp";
import { toSqlValue } from "@/utils/transaction";

const router = useRouter();
const { database, schema, table } = useApp();
schema.value = ""; // FIXME: reset schema var bc breadcrumb does not provide onclick
table.value = ""; // FIXME: reset table var bc breadcrumb does not provide onclick

async function navigateToSchema(s: string) {
  schema.value = s;
  table.value = "";
  await router.push({ name: Route.Schema });
}

const wails = useWails();

const rows = ref<FormattedQueryResult & { key: number }>();
const rowsKey = ref(0);
const filtering = ref<Array<{ id: string; value: unknown }>>([]);
const sorting = ref<Array<{ id: string; desc: boolean }>>([]);
const columns = ref<Array<client.ColumnMetadata>>();
const loading = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  loading.value = true;
  const result = await wails(() =>
    GetDatabaseSchemas(
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
watch(database, async () => {
  await fetchData();
});
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
          (row: Record<string, unknown>) =>
            navigateToSchema(
              (row.SCHEMA_NAME || row.schema_name || row.name) as string,
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
