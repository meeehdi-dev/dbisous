<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui/dist/module";
import { ref, watch } from "vue";
import { RowEmits, RowAction, cell } from "@/components/database/table/table";
import { client } from "_/go/models";

const emit = defineEmits<RowEmits>();

const {
  loading,
  data,
  table,
  primaryKey,
  actions = [],
} = defineProps<{
  loading: boolean;
  data?: client.ColumnMetadata[];
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

const key = ref(0);
watch(
  () => data,
  () => {
    key.value++;
  },
);

const columns: TableColumn<client.ColumnMetadata>[] = [
  {
    accessorKey: "name",
    header: "Name",
    cell: cell({ type: "TEXT", disabled: true }),
  },
  {
    accessorKey: "type",
    header: "Type",
    cell: cell({ type: "TEXT", disabled: true }),
  },
  {
    accessorKey: "default_value",
    header: "Default value",
    cell: cell({ type: "TEXT", disabled: true }),
  },
  {
    accessorKey: "nullable",
    header: "Nullable",
    cell: cell({ type: "BOOL", disabled: true }),
  },
  {
    accessorKey: "primary_key",
    header: "Primary key",
    cell: cell({ type: "BOOL", disabled: true }),
  },
];

const columnPinning = ref({ right: ["action"] });
</script>

<template>
  <div class="flex flex-auto flex-col gap-4 justify-between overflow-hidden">
    <div class="flex flex-auto flex-col gap-4 overflow-auto">
      <UTable
        :key="key"
        v-model:column-pinning="columnPinning"
        :data="data"
        :columns="columns"
        :loading="loading"
        :ui="{ td: 'p-0' }"
      >
        <template #action-cell="{ row: { original: row } }">
          <AppColumnActions
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
    </div>
    <AppPagination
      v-model:page="page"
      v-model:items-per-page="itemsPerPage"
      :total="data?.length"
    />
  </div>
</template>
