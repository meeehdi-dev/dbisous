<script setup lang="ts">
import { useApp } from "@/composables/shared/useApp";

const { database, schema, table } = useApp();
database.value = ""; // FIXME: reset schema var bc breadcrumb does not provide onclick
schema.value = ""; // FIXME: reset schema var bc breadcrumb does not provide onclick
table.value = ""; // FIXME: reset table var bc breadcrumb does not provide onclick
</script>

<template>
  <div class="flex flex-auto items-center justify-center">
    <UIcon name="noto:kiss-mark" class="text-primary-400 size-96 select-none" />
  </div>
</template>
