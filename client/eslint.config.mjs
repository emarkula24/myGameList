// @ts-check
import eslintJs from "@eslint/js";
import eslintReact from "@eslint-react/eslint-plugin";
import tseslint from "typescript-eslint";
export default tseslint.config(
  eslintJs.configs.recommended,
  tseslint.configs.recommended,
  tseslint.configs.stylistic,
  tseslint.configs.recommendedTypeChecked,
  eslintReact.configs["recommended-typescript"],
  {
    rules: {
      "@eslint-react/no-missing-key": "warn",
    },
    files: ["**/*.ts", "**/*.tsx"],
    ignores: ["dist"],
  },
    {
    languageOptions: {
      parserOptions: {
        projectService: true,
        // @ts-ignore
        tsconfigRootDir: import.meta.dirname,
      },
    },
  },
)




// @ts-check
// import eslintJs from "@eslint/js";
// import eslintReact from "@eslint-react/eslint-plugin";
// import tseslint from "typescript-eslint";
// export default tseslint.config({
//   files: ["**/*.ts", "**/*.tsx"],
//   ignores: ["dist"],

  // Extend recommended rule sets from:
  // 1. ESLint JS's recommended rules
  // 2. TypeScript ESLint recommended rules
  // 3. ESLint React's recommended-typescript rules
  // extends: [
  //   eslintJs.configs.recommended,
  //   tseslint.configs.recommended,
  //   tseslint.configs.strict,
  //   tseslint.configs.stylistic,
  //   tseslint.configs.recommendedTypeChecked,
  //   eslintReact.configs["recommended-typescript"],

  // ],
  

  // Custom rule overrides (modify rule levels or disable rules)
//   rules: {
//     "@eslint-react/no-missing-key": "warn",
//   },
// });