<script setup lang="ts">
const packageVersion = import.meta.env.PACKAGE_VERSION as string;
</script>

<template>
  <div class="flex flex-initial">
    <div class="flex flex-auto items-center justify-center p-4">
      <UModal>
        <UButton
          icon="simple-icons:git"
          color="neutral"
          variant="soft"
          size="sm"
          :label="packageVersion"
        />

        <template #content>TODO</template>
      </UModal>
    </div>
  </div>
</template>
