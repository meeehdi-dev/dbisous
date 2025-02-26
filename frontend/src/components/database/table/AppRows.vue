<script setup lang="ts">
import { ref, watch } from "vue";
import {
  RowEmits,
  RowAction,
  FormattedQueryResult,
} from "@/components/database/table/table";
import {
  InsertChange,
  isInsertChange,
  useTransaction,
} from "@/composables/useTransaction";
import { useWails } from "@/composables/useWails";
import { Execute } from "_/go/app/App";
import { useUrlParams } from "@/composables/useUrlParams";

const emit = defineEmits<RowEmits>();
const tx = useTransaction();

const {
  loading,
  data,
  table,
  primaryKey,
  actions = [],
} = defineProps<{
  loading: boolean;
  data?: FormattedQueryResult & { key: number };
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
const query = ref("");
function commit() {
  const sql = tx.commit();
  query.value = sql;
  open.value = true;
}

const wails = useWails();
const { databaseId } = useUrlParams();
async function execute() {
  const result = await wails(() => Execute(databaseId.value, query.value));
  if (result instanceof Error) {
    // TODO: specific error handling
  } else {
    open.value = false;
    tx.abort();
    emit("paginationChange", page.value, itemsPerPage.value);
  }
}

function abort() {
  if (data) {
    const inserted = tx.changes.value.filter((c) =>
      isInsertChange(c),
    ) as InsertChange[];
    if (inserted.length > 0) {
      console.log(inserted);
      // TODO: emit for apptable to remove inserted rows
    }
  }
  tx.abort();
}
</script>

<template>
  <div class="flex flex-auto flex-col justify-between overflow-hidden">
    <div class="flex flex-auto flex-col gap-4 overflow-auto">
      <UTable
        :data="data?.rows"
        :columns="data?.columns"
        v-model:column-pinning="columnPinning"
        :loading="loading"
        :key="data?.key"
        :ui="{ td: 'p-0 min-w-max', tbody: '[&>tr]:odd:bg-neutral-800' }"
      >
        <template #action-cell="{ row: { original: row } }">
          <AppActionsColumn
            :row="row"
            :actions="actions"
            :table="table"
            :primary-key="primaryKey"
            @view="emit(RowAction.View, row)"
            @duplicate="emit(RowAction.Duplicate, row)"
            @delete="emit(RowAction.Delete, row)"
          />
        </template>
      </UTable>
      <div
        v-if="data?.columns && data.columns.some((c) => c.id === 'actions')"
        class="flex flex-initial justify-center"
      >
        <UButton icon="lucide:plus" variant="soft" label="Add row" />
      </div>
    </div>
    <AppPagination
      v-model:page="page"
      v-model:items-per-page="itemsPerPage"
      :total="data?.total"
    />
    <div
      class="px-2 pb-2"
      :class="`${tx.changes.value.length > 0 ? 'h-18 opacity-100' : 'h-0 opacity-0'} transition-all`"
    >
      <UAlert
        color="neutral"
        variant="soft"
        icon="lucide:triangle-alert"
        :title="`${tx.changes.value.length} pending change${tx.changes.value.length > 1 ? 's' : ''}...`"
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
      :title="`Apply ${tx.changes.value.length} change${tx.changes.value.length > 1 ? 's' : ''}`"
      description="Check the content of the SQL query before executing"
      :ui="{
        content: 'max-w-none w-[80%] h-[80%]',
        body: 'sm:p-0 p-0',
        footer: 'justify-end',
      }"
    >
      <template #body>
        <AppEditor v-model="query" full />
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
