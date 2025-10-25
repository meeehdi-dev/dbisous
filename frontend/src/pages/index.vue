<script setup vapor lang="ts"></script>

<template>
  <UDashboardPanel id="home">
    <template #header>
      <UDashboardNavbar title="Home" icon="lucide:house">
        <template #leading>
          <UDashboardSidebarCollapse />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body> coucou </template>
  </UDashboardPanel>
</template>
