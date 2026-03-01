<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

onMounted(async () => {
  await router.push({ path: "/" });
});
</script>

<template>
  <UApp>
    <div class="flex h-screen w-screen flex-auto overflow-hidden">
      <AppSidebar />
      <div class="flex flex-auto overflow-hidden bg-neutral-900">
        <div class="flex flex-auto flex-col overflow-hidden">
          <AppBreadcrumb />
          <Suspense>
            <NuxtPage />
          </Suspense>
        </div>
      </div>
    </div>
  </UApp>
</template>
