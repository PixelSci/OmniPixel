# OmniPixel

OmniPixel is an all-in-one AI assistant platform for chat, image, video, audio, and other model-powered workflows.

The project is organized as a monorepo with a product app, a marketing site, an API server, and shared database initialization assets. The goal is to provide one workspace where users can manage model providers, start multimodal AI tasks, and keep their generated content and conversation history in one place.

## Apps

- `apps/omni-pixel` - the main Vue 3 + Vite product app.
- `apps/web` - the Nuxt 4 website.
- `apps/server` - the Go API server.
- `packages/db` - PostgreSQL initialization SQL.
