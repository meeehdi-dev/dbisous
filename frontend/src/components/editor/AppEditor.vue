<script setup lang="ts">
import { useTemplateRef, watch } from "vue";
import { useMonaco } from "./monaco";

const { columns } = defineProps<{
  columns: Array<unknown>;
}>();
const value = defineModel<string>({ required: true });
const container = useTemplateRef("container");
const monaco = useMonaco(value, columns);

watch(container, (c) => {
  if (!c) {
    return;
  }

  monaco.init(c);
});
</script>

<template>
  <div class="bg-neutral-950 pl-2 py-2 rounded">
    <div ref="container" class="w-full h-[200px]" />
  </div>
</template>
