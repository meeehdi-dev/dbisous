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
  <AppTypeSelect v-if="type === 'TYPE'" v-model="value as string" />
  <UCheckbox
    v-else-if="booleanTypes.includes(type)"
    v-model="value as boolean"
    :disabled="disabled"
  />
  <AppText
    v-else-if="textTypes.includes(type)"
    variant="ghost"
    v-model="value as string"
    :disabled="disabled"
    :initial-value="initialValue as string"
    :default-value="defaultValue as string"
    :nullable="nullable"
  />
  <UInput
    v-else-if="dateTypes.includes(type)"
    variant="ghost"
    :value="initialValue"
    :disabled="disabled"
    class="w-full"
  />
  <UInputNumber
    v-else-if="numberTypes.includes(type)"
    variant="ghost"
    :value="initialValue"
    :disabled="disabled"
    class="w-full"
  />
  <span v-else-if="type === ''" class="italic px-2.5">{{ initialValue }}</span>
  <span v-else class="font-bold text-red-400"
    >{{ initialValue }} ({{ type }})</span
  >
</template>
