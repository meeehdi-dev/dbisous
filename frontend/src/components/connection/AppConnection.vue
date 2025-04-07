<script setup lang="ts">
import { useRouter } from "vue-router";
import { GetConnectionDatabases, UseDatabase } from "_/go/app/App";
import {
  formatColumns,
  FormattedQueryResult,
  RowAction,
} from "@/components/connection/table/table";
import { useWails } from "@/composables/useWails";
import { ref, watch } from "vue";
import { app, client } from "_/go/models";
import { SortDirection } from "@/components/connection/table/column/AppColumnHeader.vue";
import { Route } from "@/router";
import { useApp } from "@/composables/shared/useApp";
import { toSqlValue } from "@/utils/transaction";
import { Tab } from "@/utils/tabs";
import { useConnections } from "@/composables/shared/useConnections";
import { parseConnectionString } from "@/utils/connection";

const router = useRouter();
const wails = useWails();
const { connections } = useConnections();
const { connection, database } = useApp();

const active = ref(Tab.Rows);
const defaultQuery = ref<string>();
watch(active, () => {
  if (active.value !== Tab.Query) {
    defaultQuery.value = undefined;
  }
});

const postgresPrefix = "postgres://";
async function navigateToDatabase(d: string) {
  const c = connections.value.find((c) => c.id === connection.value);
  if (!c) {
    // TODO: should add a toast to notify of an unexpected error?
    return;
  }

  if (c.type === app.ConnectionType.SQLite) {
    await router.push({ name: Route.Database });
    database.value = d;
    return;
  }

  const { host, port, user, pass, options } = parseConnectionString(
    c.connection_string,
  );

  const connectionString = `${c.type === app.ConnectionType.PostgreSQL ? postgresPrefix : ""}${user}:${pass}@${host}${port ? `:${port}` : ""}/${d}${options.length > 0 ? "?" : ""}${options.map((option) => [option.name, option.value].join(option.value ? "=" : "")).join("&")}`;

  const result = await wails(() =>
    UseDatabase(connection.value, connectionString),
  );
  if (result instanceof Error) {
    return;
  }
  await router.push({ name: Route.Database });
  database.value = d;
}

const rows = ref<FormattedQueryResult & { key: number }>();
const query = ref<string>();
const rowsKey = ref(0);
const filtering = ref<Array<{ id: string; value: unknown }>>([]);
const sorting = ref<Array<{ id: string; desc: boolean }>>([]);
const columns = ref<Array<client.ColumnMetadata>>();
const loading = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  loading.value = true;
  const result = await wails(() =>
    GetConnectionDatabases(
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
    ),
  );
  loading.value = false;
  if (result instanceof Error) {
    return;
  }
  query.value = result.query;
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
watch(connection, async () => {
  await fetchData();
});

function onQueryEdit(query: string) {
  defaultQuery.value = query;
  active.value = Tab.Query;
}
</script>

<template>
  <AppTabs v-model="active" :default-query="defaultQuery">
    <template #rows>
      <AppRows
        :loading="loading"
        :query="query"
        :data="rows"
        :sorting="sorting"
        :filtering="filtering"
        :actions="[RowAction.View]"
        @view="
          (row: Record<string, unknown>) =>
            navigateToDatabase(row.name as string)
        "
        @pagination-change="fetchData"
        @query-edit="onQueryEdit"
      />
    </template>
    <template #columns>
      <AppColumns :loading="loading" :data="columns" />
    </template>
  </AppTabs>
</template>
