<script setup lang="ts">
import { computed, ref, watch } from "vue";
import {
  RowEmits,
  RowAction,
  FormattedQueryResult,
} from "@/components/connection/table/table";
import { useTransaction } from "@/composables/shared/useTransaction";
import { useWails } from "@/composables/useWails";
import { Execute } from "_/go/app/App";
import { useApp } from "@/composables/shared/useApp";

const emit = defineEmits<
  RowEmits<Record<string, unknown>> & { queryEdit: [string] }
>();
const tx = useTransaction();

const {
  loading,
  query,
  data,
  sorting,
  filtering,
  table,
  primaryKey,
  actions = [],
} = defineProps<{
  loading: boolean;
  query?: string;
  data?: FormattedQueryResult & { key: number };
  sorting: { id: string; desc: boolean }[];
  filtering: { id: string; value: unknown }[];
  table?: string;
  primaryKey?: string;
  actions?: RowAction[];
}>();

const page = ref(1);
const itemsPerPage = ref(10);

watch(page, () => {
  emit("paginationChange", page.value, itemsPerPage.value);
});
watch(itemsPerPage, () => {
  page.value = 1;
  emit("paginationChange", page.value, itemsPerPage.value);
});

const columnPinning = ref({ right: ["action"] });

const open = ref(false);
const txQuery = ref("");
function commit() {
  const sql = tx.commit();
  txQuery.value = sql;
  open.value = true;
}

const wails = useWails();
const { database } = useApp();
async function execute() {
  const db = database.value;
  if (!db) {
    return;
  }

  const result = await wails(() => Execute(db, txQuery.value));
  if (result instanceof Error) {
    return;
  }
  open.value = false;
  tx.abort();
  emit("paginationChange", page.value, itemsPerPage.value);
}

function abort() {
  if (data && tx.insertChanges.value.length > 0) {
    // TODO: emit for apptable to remove inserted rows
  }
  tx.abort();
}

const changesCount = computed(
  () =>
    tx.insertChanges.value.length +
    tx.updateChanges.value.length +
    tx.deleteChanges.value.length,
);

function onQueryEdit() {
  if (query) {
    emit("queryEdit", query);
  }
}
</script>

<template>
  <div class="flex flex-auto flex-col justify-between overflow-hidden">
    <div class="flex flex-auto flex-col gap-2 overflow-auto">
      <div v-if="query" class="mx-2 flex min-h-9 items-center gap-2">
        <AppEditor
          v-model="query!"
          :default-value="query"
          height="full"
          disabled
        />
        <UButton
          icon="lucide:edit"
          label="Edit query"
          :ui="{ base: 'h-8' }"
          @click="onQueryEdit"
        />
      </div>
      <UTable
        :key="data?.key"
        v-model:column-pinning="columnPinning"
        :sorting="sorting"
        :sorting-options="{ manualSorting: actions.length > 0 }"
        :column-filters="filtering"
        :column-filters-options="{ manualFiltering: actions.length > 0 }"
        :data="data?.rows"
        :columns="data?.columns"
        :loading="loading"
        :ui="{
          th: 'pt-0 pb-2 px-4',
          td: 'p-0 min-w-max',
          tbody:
            data && data.rows.length > 0 ? '[&>tr]:odd:bg-neutral-800' : '',
        }"
      >
        <template #action-cell="{ row: { original: row } }">
          <AppColumnActions
            :row="row"
            :actions="actions"
            :table="table"
            :primary-key="primaryKey"
            @view="emit(RowAction.View, row as Record<string, unknown>)"
            @duplicate="
              emit(RowAction.Duplicate, row as Record<string, unknown>)
            "
            @delete="emit(RowAction.Delete, row as Record<string, unknown>)"
          />
        </template>
      </UTable>
      <div
        v-if="actions.some((a) => a === RowAction.Insert)"
        class="flex flex-initial justify-center"
      >
        <UButton
          icon="lucide:plus"
          variant="soft"
          label="Add row"
          @click="emit(RowAction.Insert)"
        />
      </div>
    </div>
    <AppPagination
      v-model:page="page"
      v-model:items-per-page="itemsPerPage"
      :total="data?.total"
    />
    <div
      :class="`px-2 ${changesCount ? 'mb-2 h-16 opacity-100' : 'mb-0 h-0 opacity-0'} overflow-hidden transition-all duration-500`"
    >
      <UAlert
        color="neutral"
        variant="soft"
        icon="lucide:triangle-alert"
        :title="`${changesCount} pending change${changesCount > 1 ? 's' : ''}...`"
        orientation="horizontal"
        :actions="[
          {
            size: 'md',
            label: 'Cancel',
            color: 'secondary',
            variant: 'soft',
            icon: 'lucide:x',
            onClick: abort,
          },
          {
            size: 'md',
            label: 'Apply',
            color: 'warning',
            variant: 'soft',
            icon: 'lucide:check',
            onClick: commit,
          },
        ]"
      />
    </div>
    <UModal
      v-model:open="open"
      :title="`Apply ${changesCount} change${changesCount > 1 ? 's' : ''}`"
      description="Check the content of the SQL query before executing"
      :ui="{
        content: 'max-w-none w-[80%] h-[80%]',
        body: 'sm:p-0 p-0',
        footer: 'justify-end',
      }"
    >
      <template #body>
        <AppEditor v-model="txQuery" height="full" />
      </template>

      <template #footer>
        <UButton
          label="Cancel"
          variant="soft"
          color="neutral"
          icon="lucide:x"
          @click="open = false"
        />
        <UButton icon="lucide:check" label="Apply" @click="execute" />
      </template>
    </UModal>
  </div>
</template>
