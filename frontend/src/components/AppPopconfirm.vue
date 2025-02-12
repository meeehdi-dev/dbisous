<script setup lang="ts">
const emit = defineEmits<{ confirm: [] }>();

const { text } = defineProps<{ text: string }>();
const open = defineModel<boolean>({ default: false });
</script>

<template>
  <UPopover
    v-model:open="open"
    :content="{
      side: 'right',
      sideOffset: 4,
    }"
    arrow
  >
    <slot />

    <template #content>
      <UCard
        :ui="{
          footer: 'flex gap-2 justify-end',
        }"
      >
        {{ text }}

        <template #footer>
          <UButton icon="lucide:x" color="error" @click="open = false" />
          <UButton icon="lucide:check" @click="emit('confirm')" />
        </template>
      </UCard>
    </template>
  </UPopover>
</template>
