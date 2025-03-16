<script setup lang="ts">
import { client } from "_/go/models";
import { computed, ref } from "vue";

export type SortDirection = false | client.OrderDirection;

const emit = defineEmits<{ sort: [SortDirection]; filter: [string] }>();

const { label, sort, filter } = defineProps<{
  label: string;
  sort: SortDirection;
  filter: unknown;
}>();

const open = ref(false);
const filterValue = ref("");

const items = computed(() => [
  {
    label: "Ascending",
    icon: "lucide:arrow-up-narrow-wide",
    color:
      sort === client.OrderDirection.Ascending
        ? ("primary" as const)
        : undefined,
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
    color:
      sort === client.OrderDirection.Descending
        ? ("primary" as const)
        : undefined,
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

function onFilterCancel() {
  filterValue.value = "";
  emit("filter", filterValue.value);
  open.value = false;
}

function onFilterConfirm() {
  emit("filter", filterValue.value);
  open.value = false;
}
</script>

<template>
  <div class="flex justify-between">
    <div class="flex items-center gap-2">
      <UDropdownMenu :items="items">
        <UButton
          :color="sort ? 'primary' : 'neutral'"
          variant="ghost"
          :ui="{ base: '-mx-2.5' }"
          :icon="
            sort
              ? sort === client.OrderDirection.Ascending
                ? 'lucide:arrow-up-narrow-wide'
                : 'lucide:arrow-down-wide-narrow'
              : 'lucide:arrow-up-down'
          "
        />
      </UDropdownMenu>
      <span>
        {{ label }}
      </span>
    </div>
    <UPopover v-model:open="open" arrow>
      <UButton
        :color="filter ? 'primary' : 'neutral'"
        variant="ghost"
        icon="lucide:filter"
        @click="open = true"
      />

      <template #content>
        <UCard
          :ui="{
            footer: 'flex sm:p-2 gap-2 justify-end',
            header: 'sm:p-2',
            body: 'sm:p-2',
          }"
        >
          <UInput v-model="filterValue" />

          <template #footer>
            <UButton icon="lucide:x" color="error" @click="onFilterCancel" />
            <UButton icon="lucide:check" @click="onFilterConfirm" />
          </template>
        </UCard>
      </template>
    </UPopover>
  </div>
</template>
