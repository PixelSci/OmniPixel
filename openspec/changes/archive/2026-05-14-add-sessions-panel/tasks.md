## 1. Backend — expose missing session routes

- [x] 1.1 Add `POST /sessions` route for `CreateSession` in `apps/server/api/routes/session_route.go`
- [x] 1.2 Add `GET /sessions/:id` route for `GetSession` in `apps/server/api/routes/session_route.go`

## 2. Frontend — API client foundation

- [x] 2.1 Create `apps/omni-pixel/src/lib/api.ts` with a thin fetch wrapper: `BASE_URL` config, `request<T>(method, path, body?)` helper that includes `credentials: "include"` and handles JSON serialization
- [x] 2.2 Add `apps/omni-pixel/src/lib/session.ts` with typed API functions: `listSessions()`, `getSession(id)`, `createSession()`, `deleteSession(id)`, `sendChat(sessionId, payload)`

## 3. Frontend — useSessions composable

- [x] 3.1 Create `apps/omni-pixel/src/composables/useSessions.ts` with reactive `sessions`, `activeId`, `loading` state, and methods `fetchSessions()`, `createSession()`, `deleteSession(id)`, `setActive(id)`
- [x] 3.2 Generate a client-side UUID for new sessions (use `crypto.randomUUID()`)
- [x] 3.3 Re-fetch session list after create and delete operations

## 4. Frontend — HigChatSessionList with real data

- [x] 4.1 Remove `defaultSessions` and accept sessions as a required prop (no default)
- [x] 4.2 Add a delete affordance (e.g., right-click context menu or hover X button) that emits a `delete` event with the session ID
- [x] 4.3 Add a confirmation dialog (or native `confirm`) before emitting delete

## 5. Frontend — wire sessions panel in pages/index.vue

- [x] 5.1 Import and use `useSessions` composable, replacing local `activeSessionId` ref
- [x] 5.2 Call `fetchSessions()` on mount
- [x] 5.3 Replace the hardcoded `messages` seed data with messages loaded from `GET /sessions/:id` on session select
- [x] 5.4 Route `handleSend` through `POST /sessions/:id/chat` instead of direct provider API calls, using model info from `useModelSettings` and API key from `useModelApiKeys`
- [x] 5.5 Handle session creation in `handleNew`: call `POST /sessions`, then set as active
- [x] 5.6 Handle session deletion from the `delete` event on `HigChatSessionList`
- [x] 5.7 Remove the provider endpoint config and direct fetch logic (`chatProviderConfig`, the old fetch call)

## 6. Polish and edge cases

- [x] 6.1 Show a loading indicator in the session list while fetching
- [x] 6.2 Disable prompt input while a new session is being created
- [x] 6.3 Handle 401 responses by clearing sessions and showing "Please sign in" state
