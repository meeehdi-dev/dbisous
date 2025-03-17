<script setup lang="ts">
import { useConnections } from "@/composables/shared/useConnections";

const emit = defineEmits<{ edit: [string] }>();

const { connections } = useConnections();

function onConnectionEdit(id: string) {
  emit("edit", id);
}
</script>

<template>
  <AppSidebarConnectionCard
    v-for="connection in connections"
    :key="connection.id"
    :value="connection"
    @edit="onConnectionEdit"
  />
</template>
