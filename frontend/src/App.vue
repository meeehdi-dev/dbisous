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
    <UHeader title="DBisous" />

    <UMain>
      <UDashboardGroup class="top-(--ui-header-height)">
        <UDashboardSidebar>
          <UDashboardSearchButton />
          <UNavigationMenu :items="items" orientation="vertical" />
        </UDashboardSidebar>

        <UDashboardSearch />

        <RouterView />
      </UDashboardGroup>
    </UMain>
  </UApp>
</template>
