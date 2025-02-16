<script setup lang="ts">
import { ref } from "vue";
import { useUrlParams } from "../../composables/useUrlParams";
import { Effect } from "effect";
import { formatColumns, FormattedQueryResult } from "./table";
import { useWails } from "../../wails";
import { GetTableInfo, GetTableRows } from "../../../wailsjs/go/app/App";

const wails = useWails();
const { databaseId, schemaId, tableId } = useUrlParams();

// const transactionQuery = ref(`SELECT * FROM ${tableId.value}`);
// const open = ref(false);

const tabs = [
  {
    label: "Rows",
    slot: "data",
    icon: "lucide:list",
  },
  {
    label: "Info",
    slot: "info",
    icon: "lucide:info",
  },
  {
    label: "Script",
    slot: "script",
    icon: "lucide:square-terminal",
  },
];

const data = ref<FormattedQueryResult>();
const fetchingData = ref(false);
async function fetchData(page = 1, itemsPerPage = 10) {
  fetchingData.value = true;
  await Effect.runPromise(
    wails(() =>
      GetTableRows(
        databaseId.value,
        page,
        itemsPerPage,
        schemaId.value,
        tableId.value,
      ),
    ).pipe(
      Effect.tap((result) => {
        data.value = {
          ...result,
          columns: formatColumns(result.columns),
        };
        fetchingData.value = false;
      }),
      Effect.catchTags({
        WailsError: Effect.succeed,
      }),
    ),
  );
}
fetchData();

const info = ref<FormattedQueryResult>();
const fetchingInfo = ref(false);
async function fetchInfo(page = 1, itemsPerPage = 10) {
  fetchingInfo.value = true;
  await Effect.runPromise(
    wails(() =>
      GetTableInfo(
        databaseId.value,
        page,
        itemsPerPage,
        schemaId.value,
        tableId.value,
      ),
    ).pipe(
      Effect.tap((result) => {
        info.value = {
          ...result,
          columns: formatColumns(result.columns, false),
        };
        fetchingInfo.value = false;
      }),
      Effect.catchTags({
        WailsError: Effect.succeed,
      }),
    ),
  );
}
fetchInfo();

const defaultQuery = ref(`SELECT * FROM ${tableId.value};`);
</script>

<template>
  <div class="flex flex-auto flex-col justify-between">
    <UTabs
      :items="tabs"
      variant="link"
      :ui="{
        root: 'flex flex-auto overflow-hidden',
        content: 'flex flex-auto flex-col gap-2 overflow-hidden',
      }"
    >
      <template #data>
        <AppRows
          :loading="fetchingData"
          :data="data"
          @pagination-change="fetchData"
        />
      </template>
      <template #info>
        <AppRows
          :loading="fetchingInfo"
          :data="info"
          @pagination-change="fetchInfo"
        />
      </template>
      <template #script>
        <AppScript v-model:default-query="defaultQuery" />
      </template>
    </UTabs>
    <!-- <div class="px-2 pb-2">
      <UAlert
        title="3 pending changes"
        icon="lucide:info"
        color="info"
        variant="soft"
        :actions="[
          {
            label: 'Apply',
            icon: 'lucide:check',
            color: 'warning',
            size: 'md',
            onClick() {
              open = true;
            },
          },
          {
            label: 'Cancel',
            icon: 'lucide:x',
            color: 'neutral',
            size: 'md',
          },
        ]"
      />
      <UModal
        v-model:open="open"
        title="Are you sure?"
        description="Please check the SQL script before applying."
      >
        <template #body>
          <div class="flex flex-col gap-8">
            <AppEditor v-model="transactionQuery" />
            <div class="flex justify-end gap-2">
              <UButton icon="lucide:check" label="Apply" />
              <UButton
                icon="lucide:x"
                color="neutral"
                label="Cancel"
                @click="open = false"
              />
            </div>
          </div>
        </template>
      </UModal>
    </div> -->
  </div>
</template>
