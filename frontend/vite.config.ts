import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";
import ui from "@nuxt/ui/vite";
import packageJson from "./package.json";

export default defineConfig({
  define: {
    "import.meta.env.PACKAGE_VERSION": JSON.stringify(packageJson.version),
  },
  plugins: [
    vue(),
    tailwindcss(),
    ui({
      ui: {
        card: {
          slots: {
            header: "p-2 sm:px-4",
            body: "p-2 sm:p-4",
            footer: "p-2 sm:px-4",
          },
        },
      },
    }),
  ],
});
