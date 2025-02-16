<script setup lang="ts">
import { ref } from "vue";
import { useConnections } from "../composables/useConnections";
import { app } from "../../wailsjs/go/models";

const { connections } = useConnections();

const slideoverOpen = ref(false);
const editedConnection = ref<app.Connection>();

function onConnectionAdded() {
  slideoverOpen.value = false;
  editedConnection.value = undefined;
}

function onConnectionEdit(connection: app.Connection) {
  editedConnection.value = connection;
  slideoverOpen.value = true;
}

const packageVersion = import.meta.env.PACKAGE_VERSION;
</script>

<template>
  <div class="min-w-72 bg-neutral-800 flex flex-col justify-between">
    <div
      class="flex flex-initial flex-col px-2 py-4 gap-4 items-center overflow-hidden"
    >
      <div class="w-full flex flex-initial flex-col gap-2 overflow-auto">
        <AppConnection
          v-for="connection in connections"
          v-bind:key="connection.id"
          :connection="connection"
          @connection-edit="onConnectionEdit"
        />
      </div>

      <UButton
        icon="lucide:plus"
        @click="slideoverOpen = true"
        label="Add connection"
        class="flex flex-initial"
      />

      <USlideover
        v-model:open="slideoverOpen"
        side="left"
        title="Add connection"
        description="Fill connection type and details to test and save"
      >
        <template #body>
          <AppConnectionForm
            v-model="editedConnection"
            @connection-added="onConnectionAdded"
          />
        </template>
      </USlideover>
    </div>
    <div class="flex flex-initial">
      <div class="p-4 flex flex-auto justify-center items-center">
        <UModal>
          <UButton
            icon="simple-icons:git"
            color="neutral"
            variant="soft"
            size="sm"
            :label="packageVersion"
          />

          <template #content>
            <Placeholder class="h-48 m-4" />
          </template>
        </UModal>
      </div>
    </div>
  </div>
</template>
