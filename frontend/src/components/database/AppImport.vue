<script setup lang="ts">
import { useApp } from "@/composables/useApp";
import { useWails } from "@/composables/useWails";
import { Import } from "_/go/app/App";

const { database } = useApp();
const wails = useWails();
// eslint-disable-next-line no-undef
const toast = useToast();

async function importFile() {
  const result = await wails(() => Import(database.value));
  if (result instanceof Error) {
    return;
  }
  toast.add({
    title: "Successfully imported database!",
    description: result,
  });
}
</script>

<template>
  <div class="p-2">
    <UButton icon="lucide:download" label="Import" @click="importFile" />
  </div>
</template>
