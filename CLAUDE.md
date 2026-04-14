This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Apps

- **apps/web** - official website, Nuxt 4
- **apps/omni-pixel** - product app, Vue 3 + Vite SPA

## UI

- **component** - [components](.claude/ui): Components that need implementation
- **style** - use tailwindcss v4

## Components

- Create application components in `components` or a relevant subfolder under `components`.
- Do not create application components in `components/ui`.
- `components/ui` is reserved for `shadcn-vue` component library source code and should only contain `shadcn-vue` primitives or generated library files.
