<script setup lang="ts">
import { ref } from "vue";
import { Effect } from "effect";
import { useWails } from "../../wails";
import { app } from "../../../wailsjs/go/models";
import { useConnections } from "../../composables/useConnections";
import { DeleteConnection } from "../../../wailsjs/go/app/App";
import { useUrlParams } from "../../composables/useUrlParams";
import { useRouter } from "vue-router";

const { connection } = defineProps<{ connection: app.Connection }>();

const wails = useWails();
const { fetchConnections } = useConnections();
const { databaseId } = useUrlParams();
const router = useRouter();

const open = ref(false);

function removeConnection(connection: app.Connection) {
  Effect.runPromise(
    wails(() => DeleteConnection(connection.id)).pipe(
      Effect.tap(() => {
        if (connection.id === databaseId.value) {
          router.push("/");
        }
      }),
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
