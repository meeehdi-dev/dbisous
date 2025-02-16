<script setup lang="ts">
import { ref } from "vue";

const page = defineModel<number>("page");
const itemsPerPage = defineModel<number>("items-per-page");

const items = ref([10, 20, 50]);

const { total = 0 } = defineProps<{
  total?: number;
}>();
</script>

<template>
  <div class="flex flex-initial flex-col">
    <USeparator />
    <div class="flex justify-between px-2 py-4">
      <div class="flex items-center gap-2">
        <UTooltip text="Items per page" :content="{ side: 'right' }">
          <USelect v-model="itemsPerPage" :items="items" />
        </UTooltip>
      </div>
      <UPagination
        v-model:page="page"
        :items-per-page="itemsPerPage"
        :total="total"
      />
      <div class="flex items-center gap-1">
        <UIcon name="lucide:list-ordered" class="text-secondary-400" />
        <span class="text-secondary-400/80 text-sm">
          {{ total }} row{{ total > 1 ? "s" : "" }}
        </span>
      </div>
    </div>
  </div>
</template>
