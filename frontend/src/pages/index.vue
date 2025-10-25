<script setup vapor lang="ts">
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
  <UDashboardGroup>
    <UDashboardSidebar collapsible>
      <UNavigationMenu :items="items" orientation="vertical" />
    </UDashboardSidebar>
    <UDashboardPanel>
      <template #header>
        <UDashboardNavbar title="Home" icon="lucide:house" />
      </template>
    </UDashboardPanel>
  </UDashboardGroup>
</template>
