import { defineConfig, PluginOption } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";
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
    // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
    "import.meta.env.PACKAGE_VERSION": JSON.stringify(packageJson.version),
  },
  plugins: [
    vue(),
    // eslint-disable-next-line @typescript-eslint/no-unsafe-call
    tailwindcss() as PluginOption,
    // eslint-disable-next-line @typescript-eslint/no-unsafe-call
    ui({
      ui: {
        colors: {
          primary: "indigo",
        },
      },
    }) as PluginOption,
  ],
});
