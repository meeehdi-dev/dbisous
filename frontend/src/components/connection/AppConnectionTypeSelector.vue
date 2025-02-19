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
        value="sqlite"
        :active="value === 'sqlite'"
        @click="emit('select', app.ConnectionType.SQLite)"
      />
      <AppConnectionTypeButton
        label="PostgreSQL"
        value="postgresql"
        :active="value === 'postgresql'"
        @click="emit('select', app.ConnectionType.PostgreSQL)"
      />
      <AppConnectionTypeButton
        label="MySQL / MariaDB"
        value="mysql"
        :active="value === 'mysql'"
        @click="emit('select', app.ConnectionType.MySQL)"
      />
    </div>
  </div>
</template>
