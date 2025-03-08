<script setup lang="ts">
import { useUrlParams } from "@/composables/useUrlParams";
import { useWails } from "@/composables/useWails";
import { Import } from "_/go/app/App";

const { databaseId } = useUrlParams();
const wails = useWails();

async function importFile() {
  const result = await wails(() => Import(databaseId.value));
  if (result instanceof Error) {
    return;
  }
}
</script>

<template>
  <div class="p-2">
    <UButton icon="lucide:download" label="Import" @click="importFile" />
  </div>
</template>
