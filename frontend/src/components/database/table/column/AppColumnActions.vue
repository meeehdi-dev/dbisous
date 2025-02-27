<script setup lang="ts">
import { RowAction, RowEmits } from "@/components/database/table/table";
import { computed } from "vue";
import { useTransaction } from "@/composables/useTransaction";

const emit = defineEmits<RowEmits>();

const {
  row,
  table,
  primaryKey,
  actions = [],
} = defineProps<{
  row: Record<string, unknown>;
  table?: string;
  primaryKey?: string;
  actions?: RowAction[];
}>();

const tx = useTransaction();

const isDeleted = computed(() => {
  let rowKey = row.__key;
  if (rowKey === undefined && primaryKey) {
    rowKey = row[primaryKey];
  }
  return tx.deleteChanges.value.some(
    (c) => c.table === table && c.rowKey === rowKey,
  );
});
</script>

<template>
  <UButton
    v-if="actions.includes(RowAction.View)"
    icon="lucide:eye"
    color="primary"
    variant="ghost"
    @click="emit(RowAction.View, row)"
  />
  <UButton
    v-if="actions.includes(RowAction.Duplicate)"
    icon="lucide:copy"
    color="secondary"
    variant="ghost"
    @click="emit(RowAction.Duplicate, row)"
  />
  <UButton
    v-if="actions.includes(RowAction.Delete)"
    :icon="`lucide:${isDeleted ? 'step-back' : 'trash'}`"
    :color="isDeleted ? 'warning' : 'error'"
    variant="ghost"
    @click="emit(RowAction.Delete, row)"
  />
</template>
