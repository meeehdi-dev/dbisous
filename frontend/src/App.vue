<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

onMounted(() => {
  router.push({ name: "welcome" });
});
</script>

<template>
  <UApp>
    <div class="flex h-screen w-screen">
      <AppSidebar />
      <div class="flex-1 bg-slate-900 overflow-auto">
        <div class="w-full flex flex-col h-full">
          <AppBreadcrumb />
          <Suspense>
            <RouterView v-slot="{ Component, route }">
              <Transition name="fade" mode="out-in">
                <component :is="Component" :key="route.path" />
              </Transition>
            </RouterView>
          </Suspense>
        </div>
      </div>
    </div>
  </UApp>
</template>
