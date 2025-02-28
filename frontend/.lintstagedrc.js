import path from "path";
import process from "process";

/**
 * @param {string[]} filenames
 */
const buildEslintCommand = (filenames) =>
  `npm run lint --file ${filenames
    .map((f) => path.relative(process.cwd(), f))
    .join(" --file ")}`;

/**
 * @type {import('lint-staged').Configuration}
 */
export default {
  "*.{js,ts,vue}": [buildEslintCommand],
};
