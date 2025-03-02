<script setup lang="ts">
import { client } from "_/go/models";
import { computed } from "vue";

export type SortDirection = false | client.OrderDirection;

const emit = defineEmits<{ sort: [SortDirection] }>();

const { label, sort } = defineProps<{
  label: string;
  sort: false | SortDirection;
}>();

const items = computed(() => [
  {
    label: "Ascending",
    icon: "lucide:arrow-up-narrow-wide",
    color: sort === client.OrderDirection.Ascending ? "primary" : undefined,
    onSelect: () => {
      emit(
        "sort",
        sort === client.OrderDirection.Ascending
          ? false
          : client.OrderDirection.Ascending,
      );
    },
  },
  {
    label: "Descending",
    icon: "lucide:arrow-down-narrow-wide",
    color: sort === client.OrderDirection.Descending ? "primary" : undefined,
    onSelect: () => {
      emit(
        "sort",
        sort === client.OrderDirection.Descending
          ? false
          : client.OrderDirection.Descending,
      );
    },
  },
]);
</script>

<template>
  <UDropdownMenu :items="items">
    <UButton
      :color="sort ? 'primary' : 'neutral'"
      variant="ghost"
      :ui="{ base: '-mx-2.5' }"
      :label="label"
      :icon="
        sort
          ? sort === client.OrderDirection.Ascending
            ? 'lucide:arrow-up-narrow-wide'
            : 'lucide:arrow-down-wide-narrow'
          : 'lucide:arrow-up-down'
      "
    />
  </UDropdownMenu>
</template>
