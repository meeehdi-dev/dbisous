<script setup lang="ts">
import { useCopy } from "@/composables/useCopy";
import { onClickOutside } from "@vueuse/core";
import { computed, ref, useTemplateRef } from "vue";

const value = defineModel<unknown>();

const { defaultValue, nullable, initialValue, disabled } = defineProps<{
  initialValue: unknown;
  defaultValue: unknown;
  nullable: boolean;
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
    onSelect: async () => {
      await copy(value.value as string);
    },
  },
  {
    label: "Reset",
    value: "reset",
    icon: "lucide:refresh-ccw",
    color: disabled || resetDisabled.value ? undefined : ("primary" as const),
    onSelect: () => {
      value.value = initialValue;
    },
    disabled: disabled || resetDisabled.value,
  },
  {
    label: "Default value",
    value: "default",
    icon: "lucide:refresh-ccw-dot",
    color:
      disabled || defaultDisabled.value ? undefined : ("secondary" as const),
    onSelect: () => {
      value.value = defaultValue;
    },
    disabled: disabled || defaultDisabled.value,
  },
  {
    label: "Set to NULL",
    value: "null",
    icon: "lucide:delete",
    color: !disabled && nullable ? ("warning" as const) : undefined,
    onSelect: () => {},
    disabled: disabled || !nullable,
  },
]);

const open = ref(false);
const container = useTemplateRef("container");

onClickOutside(null, (event) => {
  if (event.target !== container.value) {
    open.value = false;
  }
});
</script>

<template>
  <UDropdownMenu
    ref="container"
    v-model:open="open"
    :items="items"
    :content="{ side: 'right' }"
  >
    <UButton
      icon="lucide:ellipsis-vertical"
      color="neutral"
      variant="ghost"
      size="sm"
      :class="`${open ? 'opacity-100' : 'opacity-0'} transition-opacity group-hover:opacity-100`"
    />
  </UDropdownMenu>
</template>
