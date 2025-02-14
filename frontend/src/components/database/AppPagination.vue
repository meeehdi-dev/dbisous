<script setup lang="ts">
const { tableApi } = defineProps<{
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  tableApi: any | undefined;
}>();
</script>

<template>
  <div class="flex justify-center border-t border-(--ui-border) py-4">
    <UPagination
      :default-page="(tableApi?.getState().pagination.pageIndex || 0) + 1"
      :items-per-page="tableApi?.getState().pagination.pageSize"
      :total="tableApi?.getFilteredRowModel().rows.length"
      @update:page="(p) => tableApi?.setPageIndex(p - 1)"
    />
  </div>
</template>
