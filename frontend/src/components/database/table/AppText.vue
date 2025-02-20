<script setup lang="ts">
import { computed } from "vue";
import { useCopy } from "../../../composables/useCopy";

const value = defineModel<string>();

const { defaultValue, nullable, initialValue, disabled } = defineProps<{
  defaultValue: string | undefined;
  nullable: boolean;
  initialValue: string;
  disabled: boolean;
}>();

const resetDisabled = computed(() => {
  return value.value === initialValue;
});
const defaultDisabled = computed(() => defaultValue === "NULL");

const { copy } = useCopy();
const items = computed(() => [
  {
    label: "Copy to clipboard",
    value: "copy",
    icon: "lucide:copy",
    onSelect: () => {
      copy(value.value as string);
    },
  },
  {
    label: "Reset",
    value: "reset",
    icon: "lucide:refresh-ccw",
    color: disabled || resetDisabled.value ? undefined : "primary",
    onSelect: () => {
      value.value = initialValue;
    },
    disabled: disabled || resetDisabled.value,
  },
  {
    label: "Default value",
    value: "default",
    icon: "lucide:refresh-ccw-dot",
    color: disabled || defaultDisabled.value ? undefined : "secondary",
    onSelect: () => {
      value.value = defaultValue;
    },
    disabled: disabled || defaultDisabled.value,
  },
  {
    label: "Set to NULL",
    value: "null",
    icon: "lucide:delete",
    color: !disabled && nullable ? "warning" : undefined,
    onSelect: () => {},
    disabled: disabled || !nullable,
  },
]);
</script>

<template>
  <UInput
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
</template>
