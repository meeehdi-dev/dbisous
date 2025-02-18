<script setup lang="ts">
import { ref, watch } from "vue";
import { useUrlParams } from "../../composables/useUrlParams";
import { Effect } from "effect";
import { formatColumns, FormattedQueryResult, RowAction } from "./table";
import { useWails } from "../../wails";
import {
  DeletePastQuery,
  ExecuteQuery,
  GetPastQueries,
} from "../../../wailsjs/go/app/App";
import { app } from "../../../wailsjs/go/models";

const defaultQuery = defineModel<string>("defaultQuery");

const wails = useWails();
const { databaseId } = useUrlParams();

const query = ref(defaultQuery.value ?? "");
const error = ref("");

watch(query, () => {
  error.value = "";
});

const data = ref<FormattedQueryResult>();
const fetchingData = ref(false);
async function fetchData() {
  fetchingData.value = true;
  await Effect.runPromise(
    wails(() => ExecuteQuery(databaseId.value, query.value)).pipe(
      Effect.tap((result) => {
        error.value = "";
        data.value = {
          ...result,
          columns: formatColumns(result.columns, false),
        };
        fetchingData.value = false;
        fetchPastQueries();
      }),
      Effect.catchTags({
        WailsError: (err) => {
          error.value = err.message;
          data.value = undefined;
          return Effect.succeed(err);
        },
      }),
    ),
  );
}

const pastQueries = ref<app.PastQuery[]>([]);
const fetchingPastQueries = ref(false);
async function fetchPastQueries() {
  fetchingPastQueries.value = true;
  await Effect.runPromise(
    wails(GetPastQueries).pipe(
      Effect.tap((result) => {
        pastQueries.value = result;
        fetchingPastQueries.value = false;
      }),
      Effect.catchTags({
        WailsError: Effect.succeed,
      }),
    ),
  );
}
fetchPastQueries();

function removePastQuery(pastQuery: app.PastQuery) {
  Effect.runPromise(
    wails(() => DeletePastQuery(pastQuery.id)).pipe(
      Effect.tap(fetchPastQueries),
    ),
  );
}

function setQuery(q: string, execute = false) {
  query.value = q;
  defaultQuery.value = q;
  if (execute) {
    fetchData();
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
            v-bind:key="past_query.id"
            class="w-full"
          >
            <UTooltip :text="past_query.query" :content="{ side: 'left' }">
              <UButton
                color="neutral"
                variant="soft"
                size="xs"
                :label="past_query.query"
                @click="setQuery(past_query.query)"
                class="w-full"
                :ui="{ label: 'flex flex-auto' }"
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
          :disabled="!query"
          :icon="error ? 'lucide:triangle-alert' : 'lucide:terminal'"
          label="Execute"
          @click="fetchData"
          :color="error ? 'warning' : 'primary'"
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
    <AppRows :loading="fetchingData" :data="data" />
  </div>
</template>
