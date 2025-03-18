<script setup lang="ts">
import { ref } from "vue";
import { app } from "_/go/models";
import { useConnections } from "@/composables/shared/useConnections";

const { connections } = useConnections();

const slideoverOpen = ref(false);
const editedConnection = ref<Omit<app.Connection, "id">>();

function onConnectionAdd() {
  editedConnection.value = undefined;
  slideoverOpen.value = true;
}

function onConnectionAdded() {
  slideoverOpen.value = false;
  editedConnection.value = undefined;
}

function onConnectionEdit(id: string) {
  editedConnection.value = connections.value.find((c) => c.id === id);
  slideoverOpen.value = true;
}

function onConnectionDuplicate(id: string) {
  const dup = connections.value.find((c) => c.id === id);
  editedConnection.value = { ...dup, id: undefined } as Omit<
    app.Connection,
    "id"
  >;
  slideoverOpen.value = true;
}
</script>

<template>
  <div class="flex min-w-72 flex-col justify-between bg-neutral-800">
    <div
      class="flex flex-initial flex-col items-center gap-2 overflow-hidden px-2 py-2"
    >
      <div class="flex w-full flex-initial flex-col gap-2 overflow-auto">
        <AppSidebarConnectionCard
          v-for="connection in connections"
          :key="connection.id"
          :value="connection"
          @edit="onConnectionEdit"
          @duplicate="onConnectionDuplicate"
        />
      </div>

      <UButton
        icon="lucide:plus"
        label="Add connection"
        :ui="{ base: 'mt-2' }"
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
