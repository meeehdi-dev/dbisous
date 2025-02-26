<script setup lang="ts">
import { onMounted, onUnmounted, useTemplateRef, watch } from "vue";
import { useMonaco } from "@/composables/useMonaco";
import { editor } from "monaco-editor";

const { defaultValue = "", full = false } = defineProps<{
  defaultValue?: string;
  full?: boolean;
}>();
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
  <div class="flex flex-auto bg-neutral-950 pl-2 py-2 rounded h-full w-full">
    <div
      ref="container"
      :class="`flex flex-auto h-${full ? 'full' : '[100px]'} w-full`"
    />
  </div>
</template>
