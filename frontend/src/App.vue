<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

onMounted(async () => {
  await router.push({ name: "welcome" });
});
</script>

<template>
  <UApp>
    <div class="flex flex-auto w-screen h-screen overflow-hidden">
      <AppSidebar />
      <div class="flex flex-auto bg-neutral-900 overflow-hidden">
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
