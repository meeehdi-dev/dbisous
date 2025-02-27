<script setup lang="ts">
import { computed, ref } from "vue";
import {
  booleanTypes,
  dateTypes,
  numberTypes,
  textTypes,
} from "@/components/database/table/table";
import { useTransaction } from "@/composables/useTransaction";

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

const tx = useTransaction();

const valueRef = ref<unknown>(initialValue);
const value = computed({
  get: () => {
    // @ts-expect-error tkt
    let rowKey = row[primaryKey] as unknown;
    if (rowKey === "") {
      // @ts-expect-error tkt
      rowKey = row.__key;
      const change = tx.insertChanges.value.find(
        (c) => c.table === table && c.__key === rowKey,
      );
      if (change && column && change.values[column] !== undefined) {
        return change.values[column];
      } else {
        return initialValue;
      }
    } else {
      const change = tx.updateChanges.value.find(
        (c) => c.table === table && c.rowKey === rowKey,
      );
      if (change && column && change.values[column] !== undefined) {
        return change.values[column];
      } else {
        return initialValue;
      }
    }
  },
  set: (v) => {
    if (!table || !column || !primaryKey) {
      return;
    }
    // @ts-expect-error tkt
    const rowKey = row[primaryKey] as unknown;
    if (rowKey === "") {
      // @ts-expect-error tkt
      tx.updateInsert(table, row.__key, column, v);
      return;
    }
    if (v === initialValue) {
      tx.removeUpdate(table, primaryKey, rowKey, column);
    } else {
      tx.addUpdate(table, primaryKey, rowKey, column, v);
    }

    valueRef.value = v;
  },
});

const isDeleted = computed(() => {
  // @ts-expect-error tkt
  let rowKey = row[primaryKey] as unknown;
  if (rowKey === "") {
    // @ts-expect-error tkt
    rowKey = row.__key;
  }
  return tx.deleteChanges.value.some(
    (c) => c.table === table && c.rowKey === rowKey,
  );
});

// @ts-expect-error tkt
const rowKey = row.__key;
</script>

<template>
  <div
    :class="`p-1 flex gap-1 group transition-colors ${isDeleted ? 'opacity-20' : value === 'NULL' && !nullable ? 'bg-error-400/50' : rowKey !== undefined ? 'bg-warning-400/50' : value !== initialValue ? 'bg-primary-400/50' : ''}`"
  >
    <AppTypeSelect
      v-if="type.toLowerCase() === 'type'"
      v-model="value as string"
      :initial-value="initialValue as boolean"
      :default-value="defaultValue as boolean"
      :nullable="nullable"
      :disabled="disabled || isDeleted"
    />
    <AppCheckbox
      v-else-if="booleanTypes.includes(type.toLowerCase())"
      v-model="value as boolean"
      :disabled="disabled || isDeleted"
    />
    <AppText
      v-else-if="textTypes.includes(type.toLowerCase())"
      v-model="value as string"
      :disabled="disabled || isDeleted"
    />
    <AppDatePicker
      v-else-if="dateTypes.includes(type.toLowerCase())"
      v-model="value as string"
      :initial-value="initialValue as string"
      :default-value="defaultValue as string"
      :nullable="nullable"
      :disabled="disabled || isDeleted"
    />
    <AppInputNumber
      v-else-if="numberTypes.includes(type.toLowerCase())"
      v-model="value as number"
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
