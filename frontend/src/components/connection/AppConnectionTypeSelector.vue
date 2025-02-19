<script setup lang="ts">
import { app } from "../../../wailsjs/go/models";

const emit = defineEmits<{ select: [app.ConnectionType] }>();

const { value } = defineProps<{ value?: app.ConnectionType }>();
</script>

<template>
  <div class="flex justify-center">
    <div class="flex flex-col gap-4">
      <AppConnectionTypeButton
        label="SQLite"
        :value="app.ConnectionType.SQLite"
        :active="value === app.ConnectionType.SQLite"
        @click="emit('select', app.ConnectionType.SQLite)"
      />
      <AppConnectionTypeButton
        label="PostgreSQL"
        :value="app.ConnectionType.PostgreSQL"
        :active="value === app.ConnectionType.PostgreSQL"
        @click="emit('select', app.ConnectionType.PostgreSQL)"
      />
      <AppConnectionTypeButton
        label="MySQL / MariaDB"
        :value="app.ConnectionType.MySQL"
        :active="value === app.ConnectionType.MySQL"
        @click="emit('select', app.ConnectionType.MySQL)"
      />
    </div>
  </div>
</template>
