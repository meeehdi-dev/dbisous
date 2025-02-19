<script setup lang="ts">
import { ref } from "vue";
import { Effect } from "effect";
import { useWails } from "../../wails";
import { FormattedQueryResult, RowAction } from "./table";
import { formatQueryResult } from "../../effects/columns";
import { client } from "../../../wailsjs/go/models";
import { TableData } from "@nuxt/ui/dist/module";

const { fetchFn, actions = [] } = defineProps<{
  fetchFn: (page: number, itemsPerPage: number) => Promise<client.QueryResult>;
  actions?: RowAction[];
}>();

const emit = defineEmits<{ view: [TableData] }>();

const wails = useWails();

const data = ref<FormattedQueryResult>();
const fetchingData = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  await Effect.runPromise(
    wails(() => fetchFn(page, itemsPerPage)).pipe(
      Effect.andThen(formatQueryResult),
      Effect.tap((result) => {
        data.value = result;
        fetchingData.value = false;
      }),
      Effect.catchTags({
        WailsError: Effect.succeed,
      }),
    ),
  );
}
fetchData();
</script>

<template>
  <AppRows
    :loading="fetchingData"
    :data="data"
    :actions="actions"
    @view="(row) => emit('view', row)"
    @pagination-change="fetchData"
  />
</template>
