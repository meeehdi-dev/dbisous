<script setup lang="ts">
import type { TableColumn, TableData } from "@nuxt/ui/dist/module";
import { ref, watch } from "vue";
import { RowEmits, RowAction } from "@/components/database/table/table";
import { useTransaction } from "@/composables/useTransaction";

const emit = defineEmits<RowEmits>();
const tx = useTransaction();

const {
  loading,
  data,
  actions = [],
} = defineProps<{
  loading: boolean;
  data?: {
    rows?: TableData[];
    columns?: TableColumn<TableData>[];
    total?: number;
  };
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

const key = ref(0);
watch([() => data?.rows, () => data?.columns, () => loading], () => {
  key.value++;
});

const columnPinning = ref({ right: ["action"] });

function commit() {
  const sql = tx.commit();
  console.log(sql);
  // TODO: show script modal
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
        :key="key"
        :ui="{ td: 'p-1' }"
      >
        <template #action-cell="{ row }">
          <UButton
            v-if="actions.includes(RowAction.View)"
            icon="lucide:eye"
            color="primary"
            variant="ghost"
            @click="emit(RowAction.View, row)"
          />
          <UButton
            v-if="actions.includes(RowAction.Copy)"
            icon="lucide:copy"
            color="secondary"
            variant="ghost"
            @click="emit(RowAction.Copy, row)"
          />
          <UButton
            v-if="actions.includes(RowAction.Remove)"
            icon="lucide:trash"
            color="error"
            variant="ghost"
            @click="emit(RowAction.Remove, row)"
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
            label: 'Apply',
            color: 'warning',
            variant: 'soft',
            icon: 'lucide:check',
            onClick: commit,
          },
          {
            size: 'md',
            label: 'Cancel',
            color: 'secondary',
            variant: 'soft',
            icon: 'lucide:x',
            onClick: tx.abort,
          },
        ]"
      />
    </div>
  </div>
</template>
