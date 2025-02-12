<script setup lang="ts">
import { ref } from "vue";
import { Effect } from "effect";
import { useWails } from "../../wails";
import { DeleteDatabase } from "../../../wailsjs/go/app/App";
import { database } from "../../../wailsjs/go/models";
import { useDatabase } from "../../composables/useDatabase";

const { connection } = defineProps<{ connection: database.Database }>();

const wails = useWails();
const { fetchDatabases } = useDatabase();

const open = ref(false);

function removeConnection(connection: database.Database) {
  Effect.runPromise(
    wails(() => DeleteDatabase(connection.id)).pipe(Effect.tap(fetchDatabases)),
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
