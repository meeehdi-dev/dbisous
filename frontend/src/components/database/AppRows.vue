<script setup lang="ts">
import type { TableColumn, TableData } from "@nuxt/ui/dist/module";
import { ref, watch } from "vue";
import { RowEmits, RowAction } from "./table";

const emit = defineEmits<RowEmits>();

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
watch(
  () => data?.rows,
  () => {
    key.value++;
  },
);
watch(
  () => data?.columns,
  () => {
    key.value++;
  },
);

const columnPinning = ref({ right: ["action"] });
</script>

<template>
  <div class="flex flex-auto flex-col gap-4 justify-between overflow-hidden">
    <div class="flex flex-auto flex-col gap-4 overflow-auto">
      <UTable
        :data="data?.rows"
        :columns="data?.columns"
        v-model:column-pinning="columnPinning"
        :loading="loading"
        :key="key"
        sticky
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
        v-if="data?.columns && data.columns.length > 0"
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
  </div>
</template>
