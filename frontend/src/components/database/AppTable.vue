<script setup lang="ts">
import { ref } from "vue";
import { useUrlParams } from "../../composables/useUrlParams";
import { Effect } from "effect";
import { formatColumns, FormattedQueryResult } from "./table";
import { useWails } from "../../wails";
import { GetTableInfo, GetTableRows } from "../../../wailsjs/go/app/App";

const wails = useWails();
const { databaseId, schemaId, tableId } = useUrlParams();

const transactionQuery = ref(`SELECT * FROM ${tableId.value}`);
const open = ref(false);

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
const info = ref<FormattedQueryResult>();
async function getData(page = 1, itemsPerPage = 10) {
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
      }),
      Effect.catchTags({
        WailsError: Effect.succeed,
      }),
    ),
  );
}
async function getInfo(page = 1, itemsPerPage = 10) {
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
      }),
      Effect.catchTags({
        WailsError: Effect.succeed,
      }),
    ),
  );
}
getData();
getInfo();
</script>

<template>
  <div class="flex flex-col justify-between h-full">
    <UTabs
      :items="tabs"
      variant="link"
      :ui="{ root: 'h-full', content: 'flex flex-1 flex-col gap-2' }"
    >
      <template #data>
        <AppRows :data="data" @pagination-change="getData" />
      </template>
      <template #info>
        <AppRows :data="info" @pagination-change="getInfo" />
      </template>
      <template #script>
        <AppScript :default-query="`SELECT * FROM ${tableId};`" />
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
