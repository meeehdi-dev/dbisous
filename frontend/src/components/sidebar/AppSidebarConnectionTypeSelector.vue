<script setup lang="ts">
import { app } from "_/go/models";

const emit = defineEmits<{ select: [app.ConnectionType] }>();

const { value } = defineProps<{ value?: app.ConnectionType }>();

const connectionTypes = Object.entries(app.ConnectionType).map(
  ([key, value]) => ({
    label: key,
    value,
  }),
);
</script>

<template>
  <div class="flex justify-center">
    <div class="flex flex-col gap-2">
      <UButton
        v-for="connectionType in connectionTypes"
        :key="connectionType.label"
        size="xl"
        :ui="{ label: 'w-full text-center' }"
        :label="connectionType.label"
        :value="connectionType.value"
        :active="value === connectionType.value"
        :icon="`simple-icons:${connectionType.value}`"
        :trailing-icon="`lucide:square${value === connectionType.value ? '-check-big' : ''}`"
        @click="emit('select', connectionType.value)"
      />
    </div>
  </div>
</template>
