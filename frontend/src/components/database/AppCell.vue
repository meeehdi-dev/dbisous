<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { booleanTypes, dateTypes, numberTypes, textTypes } from "./table";

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

const resetDisabled = computed(() => {
  return value.value === initialValue;
});
const defaultDisabled = computed(() => defaultValue === undefined);

const bool = ref(!!initialValue);
const t = ref(initialValue as string);
const toast = useToast();

const types = ref<Array<string>>([]);
// await Effect.runPromise(
//   wails(() => GetTypes(currentDatabase.value)).pipe(
//     Effect.tap((data) => {
//       types.value = data.rows.map((d) => d.type);
//     }),
//   ),
// );
// TODO: SEPARATE EACH COMPONENT + EMIT CHANGES

const items = computed(() => [
  {
    label: "Copy to clipboard",
    value: "copy",
    icon: "lucide:copy",
    onSelect: () => {
      navigator.clipboard.writeText(value.value as string);
      toast.add({
        title: "Successfully copied to clipboard!",
        description: value.value as string,
      });
    },
  },
  {
    label: "Reset",
    value: "reset",
    icon: "lucide:refresh-ccw",
    color: resetDisabled.value ? undefined : "primary",
    onSelect: () => {
      value.value = initialValue;
    },
    disabled: disabled || resetDisabled.value,
  },
  {
    label: "Default value",
    value: "default",
    icon: "lucide:refresh-ccw-dot",
    color: defaultDisabled.value ? undefined : "secondary",
    onSelect: () => {
      value.value = defaultValue;
    },
    disabled: disabled || defaultDisabled.value,
  },
  {
    label: "Set to NULL",
    value: "null",
    icon: "lucide:delete",
    color: nullable ? "warning" : undefined,
    onSelect: () => {},
    disabled: disabled || !nullable,
  },
]);
</script>

<template>
  <USelect
    v-if="type === 'TYPE'"
    variant="ghost"
    :items="types"
    v-model="t"
    :disabled="disabled"
    class="w-full"
  />
  <UCheckbox
    v-else-if="booleanTypes.includes(type)"
    v-model="bool"
    :disabled="disabled"
  />
  <UInput
    v-else-if="textTypes.includes(type)"
    variant="ghost"
    v-model="value as string"
    :disabled="disabled"
    :highlight="value !== initialValue"
    :color="value !== initialValue ? 'warning' : undefined"
    :ui="{
      root: 'w-full group',
      base: 'pr-8 w-full overflow-ellipsis',
      trailing: 'pr-1',
    }"
  >
    <template #trailing>
      <UDropdownMenu :items="items" :content="{ side: 'right' }">
        <UButton
          icon="lucide:ellipsis-vertical"
          color="neutral"
          variant="ghost"
          size="sm"
          class="opacity-0 group-hover:opacity-100 transition-opacity"
        />
      </UDropdownMenu>
    </template>
  </UInput>
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
