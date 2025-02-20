<script setup lang="ts">
import { ref, watch } from "vue";
import {
  booleanTypes,
  dateTypes,
  numberTypes,
  textTypes,
} from "@/components/database/table/table";

const {
  initialValue,
  type = "",
  defaultValue = undefined,
  nullable = false,
  disabled,
} = defineProps<{
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
</script>

<template>
  <div class="flex gap-1 group">
    <AppTypeSelect
      v-if="type === 'TYPE'"
      v-model="value as string"
      :initial-value="initialValue as boolean"
      :default-value="defaultValue as boolean"
      :nullable="nullable"
      :disabled="disabled"
    />
    <AppCheckbox
      v-else-if="booleanTypes.includes(type)"
      v-model="value as boolean"
      :initial-value="initialValue as boolean"
      :default-value="defaultValue as boolean"
      :nullable="nullable"
      :disabled="disabled"
    />
    <AppText
      v-else-if="textTypes.includes(type)"
      v-model="value as string"
      :initial-value="initialValue as string"
      :default-value="defaultValue as string"
      :nullable="nullable"
      :disabled="disabled"
    />
    <AppDatePicker
      v-else-if="dateTypes.includes(type)"
      v-model="value as string"
      :initial-value="initialValue as string"
      :default-value="defaultValue as string"
      :nullable="nullable"
      :disabled="disabled"
    />
    <AppInputNumber
      v-else-if="numberTypes.includes(type)"
      :initial-value="initialValue as number"
      :default-value="defaultValue as number"
      :nullable="nullable"
      :disabled="disabled"
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
      :disabled="disabled"
    />
  </div>
</template>
