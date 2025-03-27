import { createSharedComposable } from "@vueuse/core";
import { app } from "_/go/models";
import { ref } from "vue";

export const useSidebar = createSharedComposable(() => {
  const slideoverOpen = ref(false);
  const editedConnection = ref<Partial<app.Connection>>();

  return { slideoverOpen, editedConnection };
});
