<script setup lang="ts">
import { computed } from "vue";

export type SortDirection = false | "desc" | "asc";

const emit = defineEmits<{ click: [SortDirection] }>();

const { label, sort } = defineProps<{
  label: string;
  sort: SortDirection;
}>();

const items = computed(() => [
  {
    label: "Asc",
    icon: "lucide:arrow-up-narrow-wide",
    color: sort === "asc" ? "primary" : undefined,
    onSelect: () => {
      emit("click", sort === "asc" ? false : "asc");
    },
  },
  {
    label: "Desc",
    icon: "lucide:arrow-down-narrow-wide",
    color: sort === "desc" ? "primary" : undefined,
    onSelect: () => {
      emit("click", sort === "desc" ? false : "desc");
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
          ? sort === 'asc'
            ? 'lucide:arrow-up-narrow-wide'
            : 'lucide:arrow-down-wide-narrow'
          : 'lucide:arrow-up-down'
      "
    />
  </UDropdownMenu>
</template>
