<script setup lang="ts">
import { ref, useTemplateRef, watch } from "vue";
import { useElementSize } from "@vueuse/core";
import { useConnections } from "../composables/useConnections";
import { app } from "../../wailsjs/go/models";

const { connections } = useConnections();

const list = useTemplateRef("list");

const { height } = useElementSize(list);

const slideoverOpen = ref(false);
const secondaryAddButton = ref(false);
const editedConnection = ref<app.Connection>();

watch(height, () => {
  secondaryAddButton.value = height.value > window.outerHeight;
});

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
  <div class="w-72 bg-slate-800 flex flex-col justify-between">
    <div class="flex flex-col px-2 py-4 gap-4 items-center">
      <UButton
        v-if="secondaryAddButton"
        icon="lucide:plus"
        @click="slideoverOpen = true"
        label="Add connection"
      />

      <div ref="list" class="w-full flex flex-col gap-2">
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
    <div>
      <div class="p-4 flex flex-1 justify-center items-center">
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
