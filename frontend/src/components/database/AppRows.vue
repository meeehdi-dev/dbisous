<script setup lang="ts">
import type { TableColumn, TableData } from "@nuxt/ui/dist/module";
import { ref, useTemplateRef, watch } from "vue";
import { Emits, RowAction } from "./table";

const emit = defineEmits<Emits>();

const table = useTemplateRef("table");

const { rows, columns, actions } = defineProps<{
  rows: TableData[];
  columns: TableColumn<TableData>[];
  actions?: RowAction[];
}>();
const key = ref(0);

watch([rows, columns], () => {
  key.value++;
});

const columnPinning = ref({ right: ["action"] });

const pagination = ref({
  pageIndex: 0,
  pageSize: 5,
});
</script>

<template>
  <div class="flex flex-1 flex-col gap-4 justify-between">
    <div class="flex flex-col gap-4">
      <UTable
        ref="table"
        v-model:pagination="pagination"
        :data="rows"
        :columns="columns"
        v-model:column-pinning="columnPinning"
        :key="key"
      >
        <template #action-cell="{ row }">
          <UButton
            v-if="actions?.includes(RowAction.View)"
            icon="lucide:eye"
            color="primary"
            variant="ghost"
            @click="emit(RowAction.View, row)"
          />
          <UButton
            v-if="actions?.includes(RowAction.Copy)"
            icon="lucide:copy"
            color="secondary"
            variant="ghost"
            @click="emit(RowAction.Copy, row)"
          />
          <UButton
            v-if="actions?.includes(RowAction.Remove)"
            icon="lucide:trash"
            color="error"
            variant="ghost"
            @click="emit(RowAction.Remove, row)"
          />
        </template>
      </UTable>
      <div v-if="columns.length > 0" class="flex justify-center">
        <UButton icon="lucide:plus" variant="soft" label="Add row" />
      </div>
    </div>
    <AppPagination :tableApi="table?.tableApi" />
  </div>
</template>
