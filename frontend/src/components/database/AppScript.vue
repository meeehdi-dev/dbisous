<script setup lang="ts">
import { ref, watch } from "vue";
import { useUrlParams } from "../../composables/useUrlParams";
import { Effect } from "effect";
import { formatColumns, FormattedQueryResult, RowAction } from "./table";
import { useWails } from "../../wails";
import { ExecuteQuery } from "../../../wailsjs/go/app/App";

const { defaultQuery } = defineProps<{ defaultQuery?: string }>();

const wails = useWails();
const { databaseId } = useUrlParams();

const query = ref(defaultQuery ?? "");
const error = ref("");

watch(query, () => {
  error.value = "";
});

const data = ref<FormattedQueryResult>();
async function executeQuery() {
  await Effect.runPromise(
    wails(() => ExecuteQuery(databaseId.value, query.value)).pipe(
      Effect.tap((result) => {
        error.value = "";
        data.value = {
          ...result,
          columns: formatColumns(result.columns),
        };
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
</script>

<template>
  <div class="flex flex-1 flex-col w-full">
    <div class="flex flex-col p-4 gap-4">
      <AppEditor v-model="query" />
      <div class="flex gap-2 items-center">
        <UButton
          :icon="error ? 'lucide:triangle-alert' : 'lucide:terminal'"
          label="Execute"
          @click="executeQuery"
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
    <AppRows :data="data" :actions="[RowAction.Copy, RowAction.Remove]" />
  </div>
</template>
