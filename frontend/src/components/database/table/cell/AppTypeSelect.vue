<script setup lang="ts">
import { ref } from "vue";

const value = defineModel<string>();

const { disabled } = defineProps<{ disabled?: boolean }>();

const types = ref<Array<string>>([]);
</script>

<template>
  <USelect
    variant="ghost"
    :items="types"
    v-model="value"
    :disabled="disabled"
  />
</template>
