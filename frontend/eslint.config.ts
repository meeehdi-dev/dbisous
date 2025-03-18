import globals from "globals";
import eslint from "@eslint/js";
import tseslint from "typescript-eslint";
import pluginVue from "eslint-plugin-vue";
import eslintPluginPrettierRecommended from "eslint-plugin-prettier/recommended";

/** @type {import('eslint').Linter.Config[]} */
export default [
  { files: ["src/**/*.{ts,vue}"] },
  {
    languageOptions: {
      globals: globals.browser,
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
        extraFileExtensions: [".vue", ".ts", ".js"],
      },
    },
  },
  eslint.configs.recommended,
  ...tseslint.configs.strictTypeChecked,
  ...pluginVue.configs["flat/recommended"],
  eslintPluginPrettierRecommended,
  {
    files: ["src/**/*.vue", "src/**/*.ts"],
    languageOptions: { parserOptions: { parser: tseslint.parser } },
  },
  {
    ignores: [
      "*.d.ts",
      "*.js",
      "dist/",
      "node_modules/",
      "src/assets/",
      "src/vite-env.d.ts",
      "wailsjs/",
    ],
  },
];
