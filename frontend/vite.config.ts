import { defineConfig, PluginOption } from "vite";
import vue from "@vitejs/plugin-vue";
import ui from "@nuxt/ui/vite";
import packageJson from "./package.json";
import { resolve } from "path";

export default defineConfig({
  resolve: {
    alias: [
      { find: "@", replacement: resolve(__dirname, "src") },
      { find: "_", replacement: resolve(__dirname, "wailsjs") },
    ],
  },
  define: {
    "import.meta.env.PACKAGE_VERSION": JSON.stringify(packageJson.version),
  },
  plugins: [
    vue(),
    ui({
      ui: {
        colors: {
          primary: "purple",
        },
      },
    }) as PluginOption,
  ],
});
