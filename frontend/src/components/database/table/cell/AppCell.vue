<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from "vue";
import {
  booleanTypes,
  dateTypes,
  numberTypes,
  textTypes,
} from "@/components/database/table/table";
import { isDeleteChange, useTransaction } from "@/composables/useTransaction";

const {
  table,
  primaryKey,
  column,
  row,
  initialValue,
  type = "",
  defaultValue = undefined,
  nullable = false,
  disabled,
} = defineProps<{
  table?: string;
  primaryKey?: string;
  column?: string;
  row?: unknown;
  initialValue: unknown;
  type?: string;
  defaultValue?: unknown;
  nullable?: boolean;
  disabled: boolean;
}>();

const value = ref(initialValue);
watch(
  () => initialValue,
  () => {
    value.value = initialValue;
  },
);

const tx = useTransaction();
const abortListener = tx.addAbortListener(() => {
  value.value = initialValue;
});
onUnmounted(() => {
  tx.removeAbortListener(abortListener);
});

watch(value, () => {
  if (!table || !column || !primaryKey) {
    return;
  }
  // @ts-expect-error tkt
  const rowKey = row[primaryKey] as unknown;
  if (!rowKey) {
    // @ts-expect-error tkt
    tx.updateInsert(table, row.__key, column, value.value);
    return;
  }
  if (value.value === initialValue) {
    tx.removeUpdate(table, primaryKey, rowKey, column);
  } else {
    tx.addUpdate(table, primaryKey, rowKey, column, value.value);
  }
});

const isDeleted = computed(() => {
  // @ts-expect-error tkt
  let rowKey = row[primaryKey] as unknown;
  if (!rowKey) {
    // @ts-expect-error tkt
    rowKey = row.__key;
  }
  return tx.changes.value.some(
    (c) => isDeleteChange(c) && c.table === table && c.rowKey === rowKey,
  );
});

// @ts-expect-error tkt
const rowKey = row.__key;
</script>

<template>
  <div :class="`flex gap-1 group${isDeleted ? ' opacity-50' : ''}`">
    <AppTypeSelect
      v-if="type === 'TYPE'"
      v-model="value as string"
      :initial-value="initialValue as boolean"
      :default-value="defaultValue as boolean"
      :nullable="nullable"
      :disabled="disabled || isDeleted"
    />
    <AppCheckbox
      v-else-if="booleanTypes.includes(type)"
      v-model="value as boolean"
      :initial-value="initialValue as boolean"
      :default-value="defaultValue as boolean"
      :nullable="nullable"
      :disabled="disabled || isDeleted"
    />
    <AppText
      v-else-if="textTypes.includes(type)"
      v-model="value as string"
      :isNew="rowKey !== undefined"
      :initial-value="initialValue as string"
      :default-value="defaultValue as string"
      :nullable="nullable"
      :disabled="disabled || isDeleted"
    />
    <AppDatePicker
      v-else-if="dateTypes.includes(type)"
      v-model="value as string"
      :initial-value="initialValue as string"
      :default-value="defaultValue as string"
      :nullable="nullable"
      :disabled="disabled || isDeleted"
    />
    <AppInputNumber
      v-else-if="numberTypes.includes(type)"
      :initial-value="initialValue as number"
      :default-value="defaultValue as number"
      :nullable="nullable"
      :disabled="disabled || isDeleted"
    />
    <span v-else-if="type === ''" class="italic px-2.5">{{
      initialValue
    }}</span>
    <span v-else class="font-bold text-red-400"
      >{{ initialValue }} ({{ type }})</span
    >
    <AppCellActions
      v-model="value"
      :initial-value="initialValue"
      :default-value="defaultValue"
      :nullable="nullable"
      :disabled="disabled || isDeleted"
    />
  </div>
</template>
