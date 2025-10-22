<script setup vapor lang="ts">
const { to, label, onClick } = defineProps<{
  label?: string;
  to?: string;
  onClick?: () => void;
}>();
</script>

<template>
  <UButton :to="to" :label="label" @click="onClick" />
</template>
