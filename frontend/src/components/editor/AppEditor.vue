<script setup lang="ts">
import { onMounted, onUnmounted, useTemplateRef } from "vue";
import { useMonaco } from "./monaco";
import { editor } from "monaco-editor";

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
</script>

<template>
  <div class="bg-neutral-950 pl-2 py-2 rounded">
    <div ref="container" class="w-full h-[100px]" />
  </div>
</template>
