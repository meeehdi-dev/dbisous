<script setup lang="ts">
import { computed } from "vue";
import { useConnections } from "../../composables/useConnections";
import { useUrlParams } from "../../composables/useUrlParams";
import { app } from "../../../wailsjs/go/models";
import { Effect } from "effect";
import { useWails } from "../../wails";
import { DeleteConnection } from "../../../wailsjs/go/app/App";
import { useRouter } from "vue-router";

const { connection } = defineProps<{ connection: app.Connection }>();
const emit = defineEmits<{ connectionEdit: [app.Connection] }>();

const { isConnected, connect, disconnect, select, fetchConnections } =
  useConnections();
const { databaseId } = useUrlParams();
const toast = useToast();
const wails = useWails();
const router = useRouter();

const connected = computed(() => isConnected(connection.id));

function getConnectionName(connection: app.Connection) {
  if (connection.name) {
    return connection.name;
  }
  const parts = connection.connection_string.replace(/\\/g, "/").split("/");
  return parts[parts.length - 1];
}

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text);
  toast.add({
    title: "Successfully copied to clipboard!",
    description: text,
  });
}

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
  <UCard
    :ui="{
      root:
        connection.id === databaseId
          ? 'border-r-2 border-r-primary-400 transition-colors'
          : connected
            ? 'cursor-pointer border-r-2 border-r-primary-400/50 hover:border-r-primary-400 transition-colors'
            : 'border-r-2 border-r-transparent transition-colors',
    }"
    @click="select(connection.id)"
  >
    <div class="flex gap-4 items-center">
      <div class="flex flex-initial">
        <UIcon :name="`simple-icons:${connection.type}`" class="size-8" />
      </div>
      <div class="flex flex-1 flex-row gap-2 justify-between">
        <UTooltip :text="connection.name" :content="{ side: 'right' }">
          <span class="line-clamp-1 text-ellipsis">
            {{ getConnectionName(connection) }}
          </span>
        </UTooltip>
        <UPopover mode="hover" :content="{ side: 'right' }">
          <div class="flex">
            <UIcon
              name="lucide:info"
              class="size-6 text-secondary-400/50 hover:text-secondary-400 transition-colors"
            />
          </div>
          <template #content>
            <div class="p-2 flex flex-col gap-2 text-gray-400">
              <UTooltip text="Connection string" :content="{ side: 'left' }">
                <div class="flex flex-row gap-2 items-center">
                  <UIcon name="lucide:link" class="text-secondary-400" />
                  <UButton
                    color="neutral"
                    variant="ghost"
                    trailing-icon="lucide:copy"
                    :ui="{ base: 'px-1' }"
                    :label="connection.connection_string"
                    @click="copyToClipboard(connection.connection_string)"
                  />
                </div>
              </UTooltip>
              <UTooltip text="Creation date" :content="{ side: 'left' }">
                <div class="flex flex-row gap-2 items-center">
                  <UIcon
                    name="lucide:calendar"
                    :class="
                      connection.created_at !== connection.updated_at
                        ? 'text-primary-400/50'
                        : 'text-primary-400'
                    "
                  />
                  <span class="text-sm">
                    {{ new Date(connection.created_at).toLocaleString() }}
                  </span>
                </div>
              </UTooltip>
              <UTooltip text="Last update" :content="{ side: 'left' }">
                <div
                  v-if="connection.created_at !== connection.updated_at"
                  class="flex flex-row gap-2 items-center"
                >
                  <UIcon name="lucide:calendar-sync" class="text-primary-400" />
                  <span class="text-sm">
                    {{ new Date(connection.updated_at).toLocaleString() }}
                  </span>
                </div>
              </UTooltip>
            </div>
          </template>
        </UPopover>
      </div>
    </div>

    <template #footer>
      <div class="flex gap-2 justify-end">
        <AppPopconfirm
          text="Are you sure?"
          @confirm="removeConnection(connection)"
        >
          <UTooltip text="Remove" :content="{ side: 'left' }">
            <UButton icon="lucide:trash" color="error" variant="soft" />
          </UTooltip>
        </AppPopconfirm>
        <UTooltip text="Edit" :content="{ side: 'top' }">
          <UButton
            icon="lucide:edit"
            color="neutral"
            variant="soft"
            @click="emit('connectionEdit', connection)"
          />
        </UTooltip>
        <UTooltip
          :text="connected ? 'Disconnect' : 'Connect'"
          :content="{ side: 'right' }"
        >
          <UButton
            :icon="connected ? 'lucide:unplug' : 'lucide:plug'"
            :color="connected ? 'warning' : 'primary'"
            variant="soft"
            @click.prevent="
              () => {
                connected ? disconnect(connection.id) : connect(connection.id);
              }
            "
          />
        </UTooltip>
      </div>
    </template>
  </UCard>
</template>
