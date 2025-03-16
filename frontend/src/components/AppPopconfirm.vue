<script setup lang="ts">
import { onClickOutside } from "@vueuse/core";
import { ref, useTemplateRef } from "vue";

const emit = defineEmits<{ confirm: [] }>();

const { text } = defineProps<{ text: string }>();
const open = ref(false);

const container = useTemplateRef("container");

onClickOutside(null, (event) => {
  if (event.target !== container.value) {
    open.value = false;
  }
});

function onConfirm() {
  emit("confirm");
  open.value = false;
}
</script>

<template>
  <UPopover
    ref="container"
    v-model:open="open"
    :content="{
      side: 'right',
      sideOffset: 4,
    }"
    arrow
  >
    <div @click="open = true">
      <slot />
    </div>

    <template #content>
      <UCard
        :ui="{
          footer: 'flex sm:p-2 gap-2 justify-end',
          header: 'sm:p-2',
          body: 'sm:p-2',
        }"
      >
        {{ text }}

        <template #footer>
          <UButton icon="lucide:x" color="error" @click="open = false" />
          <UButton icon="lucide:check" @click="onConfirm" />
        </template>
      </UCard>
    </template>
  </UPopover>
</template>
