<script setup lang="ts">
import type { NavigationMenuItem } from "@nuxt/ui";

const items: NavigationMenuItem[] = [
  {
    label: "Home",
    icon: "lucide:house",
    active: true,
  },
];
</script>

<template>
  <UApp>
    <UHeader title="" :ui="{ center: 'flex', right: 'invisible' }"
      >DBisous</UHeader
    >

    <UMain>
      <UDashboardGroup class="top-(--ui-header-height)">
        <UDashboardSidebar collapsible>
          <template #default="{ collapsed }">
            <UDashboardSearchButton :collapsed="collapsed" />

            <UNavigationMenu
              :collapsed="collapsed"
              :items="items"
              orientation="vertical"
            />
          </template>
        </UDashboardSidebar>

        <UDashboardSearch />

        <RouterView />
      </UDashboardGroup>
    </UMain>
  </UApp>
</template>
