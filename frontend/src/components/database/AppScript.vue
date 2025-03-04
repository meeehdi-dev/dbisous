<script setup lang="ts">
import { ref, watch } from "vue";
import { useUrlParams } from "@/composables/useUrlParams";
import {
  formatColumns,
  FormattedQueryResult,
} from "@/components/database/table/table";
import { useWails } from "@/composables/useWails";
import { DeletePastQuery, ExecuteQuery, GetPastQueries } from "_/go/app/App";
import { app, client } from "_/go/models";
import { SortDirection } from "@/components/database/table/column/AppColumnHeader.vue";

const defaultQuery = defineModel<string>("defaultQuery");

const wails = useWails();
const { databaseId } = useUrlParams();

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
  const result = await wails(() => ExecuteQuery(databaseId.value, query.value));
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
</script>

<template>
  <div class="flex flex-auto flex-col">
    <div class="flex flex-col p-4 gap-4">
      <div class="flex gap-2">
        <div class="flex flex-1/2">
          <AppEditor v-model="query" :default-value="defaultQuery" />
        </div>
        <div
          class="flex flex-1/2 flex-col gap-2 p-2 bg-neutral-950 rounded h-[116px] overflow-auto"
        >
          <div
            v-for="past_query in pastQueries"
            :key="past_query.id"
            class="w-full"
          >
            <UTooltip :text="past_query.query" :content="{ side: 'left' }">
              <UButton
                color="neutral"
                variant="soft"
                size="xs"
                :label="past_query.query"
                class="w-full"
                :ui="{ label: 'flex flex-auto' }"
                @click="setQuery(past_query.query)"
              >
                <template #trailing>
                  <div class="flex gap-1">
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
                  </div>
                </template>
              </UButton>
            </UTooltip>
          </div>
        </div>
      </div>
      <div class="flex gap-2 items-center">
        <UButton
          :disabled="!query || error !== ''"
          :icon="error ? 'lucide:triangle-alert' : 'lucide:terminal'"
          label="Execute"
          :color="error ? 'warning' : 'primary'"
          @click="() => fetchData()"
        />
        <span
          :class="`text-sm text-neutral-400 pointer-events-none transition-opacity ${data && data.duration ? 'opacity-100' : 'opacity-0'}`"
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
