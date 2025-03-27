<script setup lang="ts">
import { useMagicKeys } from "@vueuse/core";
import { computed } from "vue";

const { kbds, placement } = defineProps<{
  kbds: string[];
  placement?: "top" | "bottom" | "left" | "right";
}>();

const keys = useMagicKeys();

const open = computed(() => keys.meta.value && keys["?"].value);
</script>
<template>
  <UTooltip
    v-model:open="open"
    :content="{ side: placement }"
    :kbds="kbds"
    :ui="{
      content: 'bg-secondary-400',
      kbds: 'before:content-[\'\'] before:me-0',
    }"
  >
    <slot />
  </UTooltip>
</template>
