<script setup lang="ts">
import { useConnections } from "@/composables/shared/useConnections";
import { app } from "_/go/models";

const emit = defineEmits<{ connectionEdit: [app.Connection] }>();

const { connections } = useConnections();

function onConnectionEdit(connection: app.Connection) {
  emit("connectionEdit", connection);
}
</script>

<template>
  <AppSidebarConnectionCard
    v-for="connection in connections"
    :key="connection.id"
    :connection="connection"
    @connection-edit="onConnectionEdit"
  />
</template>
