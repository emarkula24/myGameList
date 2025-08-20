// @ts-check
import eslintJs from "@eslint/js";
import eslintReact from "@eslint-react/eslint-plugin";
import tseslint from "typescript-eslint";
export default tseslint.config(
  {
    ignores: ["**/dist/**", "eslint.config.mjs"],
    
  },
  eslintJs.configs.recommended,
  tseslint.configs.recommended,
  tseslint.configs.strict,
  tseslint.configs.stylistic,
  tseslint.configs.recommendedTypeChecked,
  eslintReact.configs["recommended-typescript"],
  {
    rules: {
      // "@eslint-react/no-missing-key": "warn",
      "@typescript-eslint/only-throw-error": "off"
    },
    files: ["**/*.ts", "**/*.tsx"],
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