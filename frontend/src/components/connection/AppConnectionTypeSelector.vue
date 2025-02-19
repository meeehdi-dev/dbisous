<script setup lang="ts">
const emit = defineEmits<{ select: ["sqlite" | "mysql" | "postgresql"] }>();

const { value } = defineProps<{ value: "sqlite" | "mysql" | "postgresql" }>();
</script>

<template>
  <div class="flex justify-center">
    <div class="flex flex-col gap-4">
      <AppConnectionTypeButton
        label="SQLite"
        value="sqlite"
        :active="value === 'sqlite'"
        @click="emit('select', 'sqlite')"
      />
      <AppConnectionTypeButton
        label="PostgreSQL"
        value="postgresql"
        :active="value === 'postgresql'"
        @click="emit('select', 'postgresql')"
      />
      <AppConnectionTypeButton
        label="MySQL / MariaDB"
        value="mysql"
        :active="value === 'mysql'"
        @click="emit('select', 'mysql')"
      />
    </div>
  </div>
</template>
