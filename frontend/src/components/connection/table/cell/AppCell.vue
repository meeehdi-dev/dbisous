<script setup lang="ts">
import { computed, ref } from "vue";
import {
  booleanTypes,
  dateTypes,
  enumTypes,
  numberTypes,
  textTypes,
} from "@/components/connection/table/table";
import { useTransaction } from "@/composables/shared/useTransaction";

const {
  table,
  primaryKey,
  column,
  row,
  initialValue,
  type = "",
  items = [],
  defaultValue = undefined,
  nullable = false,
  disabled,
} = defineProps<{
  table?: string;
  primaryKey?: string;
  column?: string;
  row?: Record<string, unknown>;
  initialValue: unknown;
  type?: string;
  items?: string[];
  defaultValue?: unknown;
  nullable?: boolean;
  disabled: boolean;
}>();

const tx = useTransaction();

const valueRef = ref<unknown>(initialValue);
const value = computed({
  get: () => {
    if (!row) {
      return initialValue;
    }

    let rowKey = row.__key;
    if (rowKey !== undefined) {
      const change = tx.insertChanges.value.find(
        (c) => c.table === table && c.__key === rowKey,
      );
      if (change && column && change.values[column] !== undefined) {
        return change.values[column];
      } else {
        return initialValue;
      }
    } else if (primaryKey) {
      rowKey = row[primaryKey];
      const change = tx.updateChanges.value.find(
        (c) => c.table === table && c.rowKey === rowKey,
      );
      if (change && column && change.values[column] !== undefined) {
        return change.values[column];
      } else {
        return initialValue;
      }
    } else {
      return initialValue;
    }
  },
  set: (v) => {
    if (!table || !column || !primaryKey || !row) {
      return;
    }
    let rowKey = row.__key;
    if (rowKey !== undefined) {
      tx.updateInsert(table, rowKey as number, column, v);
      return;
    }
    rowKey = row[primaryKey];
    if (v === initialValue) {
      tx.removeUpdate(table, primaryKey, rowKey, column);
    } else {
      tx.addUpdate(table, primaryKey, rowKey, column, v);
    }

    valueRef.value = v;
  },
});

const isDeleted = computed(() => {
  if (!row) {
    return false;
  }

  let rowKey = row.__key;
  if (rowKey !== undefined || primaryKey === undefined) {
    return false;
  }

  rowKey = row[primaryKey];
  return tx.deleteChanges.value.some(
    (c) => c.table === table && c.rowKey === rowKey,
  );
});

const isNullError = computed(() => value.value === "NULL" && !nullable);
const isNew = computed(() => row?.__key !== undefined);
const isDirty = computed(() => value.value !== initialValue);
</script>

<template>
  <div
    :class="[
      'group flex gap-1 p-1 transition-colors',
      disabled
        ? ''
        : isDeleted
          ? 'opacity-20'
          : isNullError
            ? 'bg-error-400/50'
            : isNew
              ? 'bg-warning-400/50'
              : isDirty
                ? 'bg-primary-400/50'
                : '',
    ]"
  >
    <AppTypeSelect
      v-if="type.toLowerCase() === 'type'"
      v-model="value as string"
      :disabled="disabled || isDeleted"
    />
    <AppSelect
      v-else-if="enumTypes.includes(type.toLowerCase())"
      v-model="value as string"
      :items="items"
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
      :disabled="disabled || isDeleted"
    />
    <span v-else-if="type === ''" class="px-2.5 italic">{{
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
