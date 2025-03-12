<script setup lang="ts">
import { ref } from "vue";
import { app } from "_/go/models";

const slideoverOpen = ref(false);
const editedConnection = ref<app.Connection>();

function onConnectionAdd() {
  editedConnection.value = undefined;
  slideoverOpen.value = true;
}

function onConnectionAdded() {
  slideoverOpen.value = false;
  editedConnection.value = undefined;
}

function onConnectionEdit(connection: app.Connection) {
  editedConnection.value = connection;
  slideoverOpen.value = true;
}
</script>

<template>
  <div class="min-w-72 bg-neutral-800 flex flex-col justify-between">
    <div
      class="flex flex-initial flex-col px-2 py-4 gap-4 items-center overflow-hidden"
    >
      <div class="w-full flex flex-initial flex-col gap-2 overflow-auto">
        <AppSidebarConnections @connection-edit="onConnectionEdit" />
      </div>

      <UButton
        icon="lucide:plus"
        label="Add connection"
        class="flex flex-initial"
        @click="onConnectionAdd"
      />

      <USlideover
        v-model:open="slideoverOpen"
        side="left"
        title="Add connection"
        description="Fill connection type and details to test and save"
      >
        <template #body>
          <AppSidebarConnectionForm
            v-model="editedConnection"
            @connection-added="onConnectionAdded"
          />
        </template>
      </USlideover>
    </div>
    <AppVersion />
  </div>
</template>
