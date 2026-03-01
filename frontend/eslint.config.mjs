import withNuxt from "./.nuxt/eslint.config.mjs";
import eslintPluginPrettierRecommended from "eslint-plugin-prettier/recommended";

export default withNuxt(
  {
    ignores: ["bindings/**", "wailsjs/**", "dist/**", ".nuxt/**"],
  },
  eslintPluginPrettierRecommended,
  {
    rules: {
      "vue/multi-word-component-names": "off",
      "@typescript-eslint/no-explicit-any": "off",
    },
  },
);
