<script setup lang="ts">
import { onMounted, onUnmounted, useTemplateRef, watch } from "vue";
import { useMonaco } from "@/composables/useMonaco";
import { editor } from "monaco-editor";

const {
  defaultValue = "",
  height = 100,
  disabled = false,
} = defineProps<{
  defaultValue?: string;
  height?: number | "full";
  disabled?: boolean;
}>();
const value = defineModel<string>({ required: true });
const container = useTemplateRef("container");
const { create } = useMonaco();

let e: editor.IStandaloneCodeEditor;
onMounted(() => {
  if (!container.value) {
    return;
  }

  e = create(container.value, value, disabled);
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
  <div class="flex h-full w-full flex-auto rounded bg-neutral-950 py-2 pl-2">
    <div
      ref="container"
      :class="[
        'flex w-full flex-auto',
        `h-${height === 'full' ? 'full' : `[${height}px]`}`,
      ]"
    />
  </div>
</template>
