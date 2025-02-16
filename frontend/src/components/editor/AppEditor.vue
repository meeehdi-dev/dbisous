<script setup lang="ts">
import { onMounted, onUnmounted, useTemplateRef, watch } from "vue";
import { useMonaco } from "./monaco";
import { editor } from "monaco-editor";

const { defaultValue = "" } = defineProps<{ defaultValue?: string }>();
const value = defineModel<string>({ required: true });
const container = useTemplateRef("container");
const { create } = useMonaco();

let e: editor.IStandaloneCodeEditor;
onMounted(() => {
  e = create(container.value!, value);
});
onUnmounted(() => {
  e.dispose();
});

watch(
  () => defaultValue,
  () => {
    e.setValue(defaultValue);
  },
);
</script>

<template>
  <div class="flex flex-auto bg-neutral-950 pl-2 py-2 rounded">
    <div ref="container" class="flex flex-auto h-[100px] w-full" />
  </div>
</template>
