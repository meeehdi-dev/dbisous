<script setup lang="ts">
import { useApp } from "@/composables/shared/useApp";
import { useWails } from "@/composables/useWails";
import { ImportDatabase } from "_/dbisous/app/app.js";

const { connection } = useApp();
const wails = useWails();

const toast = useToast();

async function importFile() {
  const result = await wails(() => ImportDatabase(connection.value));
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
