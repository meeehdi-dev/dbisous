<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";
import { Route } from "./router";

const router = useRouter();

onMounted(async () => {
  await router.push({ name: Route.Welcome });
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
