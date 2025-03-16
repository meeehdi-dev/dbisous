<script setup lang="ts">
import { useApp } from "@/composables/shared/useApp";
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
  <div class="flex flex-auto items-center justify-center">
    <UButton icon="lucide:download" label="Import" @click="importFile" />
  </div>
</template>
