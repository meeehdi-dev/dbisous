<script setup lang="ts">
import { app } from "_/go/models";
import { useCopy } from "@/composables/useCopy";

const { connection } = defineProps<{ connection: app.Connection }>();

const { copy } = useCopy();
</script>

<template>
  <UPopover mode="hover" :content="{ align: 'start', side: 'right' }" arrow>
    <UButton icon="lucide:info" variant="soft" color="secondary" />
    <template #content>
      <div class="flex flex-col gap-2 p-2 text-gray-400">
        <UTooltip text="Connection string" :content="{ side: 'left' }">
          <div class="flex flex-row items-center gap-2">
            <UIcon name="lucide:link" class="text-secondary-400" />
            <UButton
              color="neutral"
              variant="ghost"
              trailing-icon="lucide:copy"
              :ui="{ base: 'px-1' }"
              :label="connection.connection_string"
              @click="copy(connection.connection_string)"
            />
          </div>
        </UTooltip>
        <UTooltip text="Creation date" :content="{ side: 'left' }">
          <div class="flex flex-row items-center gap-2">
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
            class="flex flex-row items-center gap-2"
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
</template>
