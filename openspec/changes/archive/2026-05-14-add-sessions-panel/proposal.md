## Why

The session list sidebar currently renders hardcoded mock data and the chat flow bypasses the backend entirely, calling OpenAI/DeepSeek directly from the browser. The backend already has full session CRUD and a session chat endpoint, but the frontend hasn't been wired to it yet. Connecting these pieces makes sessions persistent and enables server-side chat with proper session management.

## What Changes

- Expose the missing `POST /sessions` and `GET /sessions/:id` routes on the backend
- Create a `useSessions` composable to fetch, create, and delete sessions via the backend API
- Replace hardcoded `defaultSessions` in `HigChatSessionList` with live data from the composable
- Add session delete (with confirmation) to the session list UI
- Route chat sending through the backend's `POST /sessions/:id/chat` endpoint instead of calling LLM providers directly from the browser
- Load persisted chat messages when switching sessions
- Auto-create a backend session on first message of a new chat

## Capabilities

### New Capabilities
- `session-management`: CRUD operations for chat sessions — list, create, delete, and switch between sessions persisted on the server
- `session-chat`: Send chat prompts through the backend session endpoint, receiving AI responses that are persisted server-side

### Modified Capabilities
<!-- None - no existing spec files to modify -->

## Impact

- **Backend**: Add routes for `POST /sessions` and `GET /sessions/:id` (controller methods already exist)
- **Frontend components**: `HigChatSessionList.vue` (real data, delete action), `pages/index.vue` (chat routing, session loading)
- **New composable**: `useSessions.ts`
- **API client**: Need axios service or fetch wrapper for backend API calls
- **No breaking changes** — the mock data rendering is internal-only
