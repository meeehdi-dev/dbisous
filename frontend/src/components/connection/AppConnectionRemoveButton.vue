<script setup lang="ts">
import { ref } from "vue";
import { Effect } from "effect";
import { useWails } from "../../wails";
import { database } from "../../../wailsjs/go/models";
import { useConnections } from "../../composables/useConnections";
import { DeleteConnection } from "../../../wailsjs/go/app/App";

const { connection } = defineProps<{ connection: database.Connection }>();

const wails = useWails();
const { fetchConnections } = useConnections();

const open = ref(false);

function removeConnection(connection: database.Connection) {
  Effect.runPromise(
    wails(() => DeleteConnection(connection.id)).pipe(
      Effect.tap(fetchConnections),
    ),
  );
}
</script>

<template>
  <AppPopconfirm
    text="Are you sure?"
    v-model="open"
    @confirm="removeConnection(connection)"
  >
    <UTooltip text="Remove" :content="{ side: 'left' }">
      <UButton
        icon="lucide:trash"
        color="error"
        variant="soft"
        @click="open = true"
      />
    </UTooltip>
  </AppPopconfirm>
</template>
