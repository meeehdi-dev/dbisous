import { resolve } from "path";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2024-11-01",
  ssr: false,
  modules: ["@nuxt/ui", "@vueuse/nuxt", "@nuxt/eslint"],
  css: ["~/style.css"],
  eslint: {
    config: {
      stylistic: false, // let prettier handle it
    },
  },
  ui: {
    // any Nuxt UI v4 options
  },
  srcDir: "src/",
  alias: {
    "@": resolve(__dirname, "src"),
    _: resolve(__dirname, "bindings"),
  },
  components: [
    {
      path: "~/components",
      pathPrefix: false, // so components are global, mimicking the auto import behavior of typical Vue setups if needed, though nuxt UI v4 auto imports them prefixed.
    },
  ],
  typescript: {
    tsConfig: {
      compilerOptions: {
        allowJs: true,
        strictNullChecks: false,
        verbatimModuleSyntax: false,
      },
      include: ["../bindings/**/*"],
    },
  },
});
