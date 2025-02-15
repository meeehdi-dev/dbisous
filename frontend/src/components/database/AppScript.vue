<script setup lang="ts">
import { computed, ref } from "vue";
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

const results = computed(() => {
  if (!data.value) {
    return "No active query";
  }

  switch (data.value.rows.length) {
    case 0:
      return "No result";
    case 1:
      return "1 result";
    default:
      return data.value.rows.length + " results";
  }
});
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
      </div>
      <UBadge v-if="error" color="warning">
        {{ error }}
      </UBadge>
    </div>
    <USeparator :label="results" />
    <AppRows
      :rows="data?.rows"
      :columns="data?.columns"
      :actions="[RowAction.Copy, RowAction.Remove]"
    />
  </div>
</template>
