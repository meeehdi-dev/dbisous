<script setup lang="ts">
const { tableApi } = defineProps<{
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  tableApi: any | undefined;
}>();
</script>

<template>
  <div class="flex flex-col items-center py-4 gap-4">
    <USeparator />
    <UPagination
      :default-page="(tableApi?.getState().pagination.pageIndex || 0) + 1"
      :items-per-page="tableApi?.getState().pagination.pageSize"
      :total="tableApi?.getFilteredRowModel().rows.length"
      @update:page="(p) => tableApi?.setPageIndex(p - 1)"
    />
  </div>
</template>
