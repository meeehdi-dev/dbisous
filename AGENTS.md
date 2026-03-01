# Agent Instructions for DBisous

This file contains instructions for AI coding agents operating in this repository. 

## 1. Stack Overview
DBisous is a desktop application built using **Wails**. 
- **Backend**: Go with `database/sql` (SQLite, PostgreSQL, MySQL).
- **Frontend**: Vue 3 (Composition API), Vite, TypeScript, Nuxt UI, and Tailwind CSS.

## 2. Build, Lint, and Test Commands

### Backend (Go - Project Root)
- **Run Dev Mode**: `make dev` (or `wails dev`)
- **Build**: `make build` (or `wails build`)
- **Run all tests**: `make test` (or `go test ./app`)
- **Run a single test**: `go test -run <TestName> ./app`
- **Dependency Management**: `go mod tidy`

### Frontend (Vue/TS - `/frontend` Directory)
*Note: Run these commands inside the `frontend` directory.*
- **Install dependencies**: `npm install`
- **Build**: `npm run build`
- **Run all tests**: `npm run test` (uses Vitest)
- **Run a single test**: `npm run test -- -t "<test_name_or_pattern>"`
- **Linting**: `npm run lint` (ESLint)
- **Type Checking**: `npm run typecheck` (vue-tsc)

## 3. Code Style Guidelines

### Backend (Go)
- **Formatting**: Always use standard `go fmt`.
- **Naming Conventions**: Use `camelCase` for local variables/private functions, and `PascalCase` for exported identifiers.
- **Error Handling**: Explicitly check and handle errors (`if err != nil { ... }`). Do not suppress them.
- **Testing**: Use standard Go testing with `github.com/stretchr/testify/require` for assertions. Tests are located alongside code (e.g., `app/connection_test.go`).
- **Wails Integration**: Expose necessary backend methods to the Wails runtime for frontend consumption.

### Frontend (Vue 3 / TypeScript)
- **Components**: Use Vue 3 `<script setup lang="ts">` syntax. Component files must be in `PascalCase` and generally prefixed with `App` (e.g., `AppSidebar.vue`, `AppConnection.vue`).
- **Logic & State**: Heavily utilize Vue Composables (found in `frontend/src/composables/`) to encapsulate state and business logic. Keep UI components slim.
- **Styling**: Use Tailwind CSS utility classes directly in the `<template>`.
- **Typing**: Enforce strict TypeScript types. Define explicit interfaces/types for data structures and avoid `any`.
- **Testing**: Use Vitest (`.test.ts` files). Test utility functions and composables thoroughly.
- **Formatting**: Handled by Prettier (`prettier-plugin-tailwindcss`).

## 4. Workflows
- Always ensure both Go and frontend tests pass after making full-stack changes.
- Ensure frontend types are correct by running `npm run typecheck` in the `/frontend` directory.
- No existing Cursor rules (`.cursorrules`) or Copilot rules (`.github/copilot-instructions.md`) were found, so strictly adhere to the conventions observed in the existing codebase.