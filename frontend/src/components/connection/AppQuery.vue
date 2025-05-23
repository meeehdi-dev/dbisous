<script setup lang="ts">
import { ref, watch } from "vue";
import {
  formatColumns,
  FormattedQueryResult,
} from "@/components/connection/table/table";
import { useWails } from "@/composables/useWails";
import { DeletePastQuery, ExecuteQuery, GetPastQueries } from "_/go/app/App";
import { app, client } from "_/go/models";
import { SortDirection } from "@/components/connection/table/column/AppColumnHeader.vue";
import { useApp } from "@/composables/shared/useApp";

const defaultQuery = defineModel<string>("defaultQuery");

const wails = useWails();
const { connection } = useApp();

const query = ref(defaultQuery.value ?? "");
const error = ref("");

watch(query, () => {
  error.value = "";
});

const data = ref<FormattedQueryResult & { key: number }>();
const dataKey = ref(0);
const sorting = ref<Array<{ id: string; desc: boolean }>>([]);
const filtering = ref<Array<{ id: string; value: unknown }>>([]);
const fetchingData = ref(false);
async function fetchData(reload = true) {
  fetchingData.value = true;
  const result = await wails(() => ExecuteQuery(connection.value, query.value));
  fetchingData.value = false;
  if (result instanceof Error) {
    error.value = result.message;
    data.value = undefined;
    return;
  }
  data.value = {
    key: dataKey.value++,
    // eslint-disable-next-line @typescript-eslint/no-misused-spread
    ...result,
    columns: formatColumns(
      result.columns,
      (name: string, s: SortDirection) => {
        if (!s) {
          sorting.value = [];
        } else {
          sorting.value = [
            { id: name, desc: s === client.OrderDirection.Descending },
          ];
        }
      },
      async (name: string, f: unknown) => {
        if (!f) {
          filtering.value = [];
        } else {
          filtering.value = [{ id: name, value: f }];
        }
        return fetchData();
      },
      undefined,
      undefined,
      true,
    ),
  };
  if (reload) {
    await fetchPastQueries();
  }
}

const pastQueries = ref<Array<app.PastQuery>>([]);
const fetchingPastQueries = ref(false);
async function fetchPastQueries() {
  fetchingPastQueries.value = true;
  const result = await wails(GetPastQueries);
  if (result instanceof Error) {
    return;
  }
  pastQueries.value = result;
  fetchingPastQueries.value = false;
}
await fetchPastQueries();

async function removePastQuery(pastQuery: app.PastQuery) {
  const result = await wails(() => DeletePastQuery(pastQuery.id));
  if (result instanceof Error) {
    return;
  }
  await fetchPastQueries();
}

async function setQuery(q: string, execute = false) {
  query.value = q;
  defaultQuery.value = q;
  if (execute) {
    await fetchData(false);
  }
}

// eslint-disable-next-line no-undef
defineShortcuts({
  meta_enter: () => {
    void fetchData();
  },
});
</script>

<template>
  <div class="flex flex-auto flex-col gap-2 overflow-hidden">
    <div class="flex flex-col gap-2">
      <div class="flex gap-2 px-2">
        <!-- TODO: handle height here and children as h-full -->
        <div class="flex flex-1/2">
          <AppEditor v-model="query" :default-value="defaultQuery" />
        </div>
        <div
          class="flex h-[216px] flex-1/2 flex-col gap-2 overflow-auto rounded bg-neutral-950 p-2"
        >
          <div
            v-for="past_query in pastQueries"
            :key="past_query.id"
            class="w-full"
          >
            <UButton
              color="neutral"
              variant="soft"
              size="xs"
              :label="past_query.query"
              class="w-full"
              :ui="{ label: 'flex flex-auto', trailingIcon: 'flex gap-1' }"
              @click="setQuery(past_query.query)"
            >
              <template #trailing>
                <UButton
                  variant="ghost"
                  size="xs"
                  icon="lucide:play"
                  @click.stop="setQuery(past_query.query, true)"
                />
                <UButton
                  color="warning"
                  variant="ghost"
                  size="xs"
                  icon="lucide:trash"
                  @click.stop="removePastQuery(past_query)"
                />
              </template>
            </UButton>
          </div>
        </div>
      </div>
      <div class="flex items-center gap-2 px-2">
        <AppKbdTooltip :kbds="['meta', 'enter']" placement="right">
          <UButton
            :disabled="!query || error !== ''"
            :icon="error ? 'lucide:triangle-alert' : 'lucide:terminal'"
            label="Run query"
            :color="error ? 'warning' : 'primary'"
            @click="() => fetchData()"
          />
        </AppKbdTooltip>
        <!-- TODO: add button to execute script from sql file -->
        <span
          :class="`pointer-events-none text-sm text-neutral-400 transition-opacity ${data && data.duration ? 'opacity-100' : 'opacity-0'}`"
          >{{ data?.duration }}</span
        >
        <UBadge v-if="error" color="warning">
          {{ error }}
        </UBadge>
      </div>
    </div>
    <AppRows
      :loading="fetchingData"
      :data="data"
      :sorting="sorting"
      :filtering="filtering"
    />
  </div>
</template>
