## Context

The OmniPixel product app (`apps/omni-pixel`) has a chat sidebar (`HigChatSessionList`) that renders hardcoded session data. The backend (`apps/server`) has full session CRUD with a JWT-protected REST API, but the frontend has no API client layer and calls third-party LLM APIs directly from the browser. The backend controller has `CreateSession` and `GetSession` methods that are not yet exposed via routes.

## Goals / Non-Goals

**Goals:**
- Connect the session list sidebar to real backend data
- Route chat messages through the backend session chat endpoint
- Let users create and delete sessions from the UI
- Load persisted messages when switching sessions

**Non-Goals:**
- Session rename / title editing
- Session search or filtering
- Offline support
- Real-time collaboration / multi-device sync
- User authentication flow (assumes JWT is already handled)

## Decisions

### 1. API client: use `fetch` (not axios)

The app already has no HTTP library dependency. Using native `fetch` avoids adding a dependency for a handful of endpoints. The composable will wrap `fetch` with a thin helper for the backend base URL and JSON parsing. If the API surface grows later, introducing axios is a low-cost refactor.

### 2. Composable architecture: `useSessions` + inline chat in page

Session list state (CRUD, active session) lives in `useSessions`. Chat state (messages, send, loading) can either be a separate `useChat` composable or remain in the page component. Since chat is currently page-local and tightly coupled to the prompt input, we keep chat logic in `pages/index.vue` initially and extract a `useSessionChat` composable only if complexity warrants it. `useSessions` is the new piece — it owns session list fetching, create, delete, and active session ID.

```
┌─────────────────────────────────────────────────┐
│  pages/index.vue                                │
│  ┌──────────────┐  ┌──────────────────────────┐ │
│  │ HigSidebar   │  │ Chat area                │ │
│  │ ┌──────────┐ │  │ messages (local state)   │ │
│  │ │Session   │ │  │ HigPromptInput           │ │
│  │ │List      │ │  │                          │ │
│  │ │          │ │  │ on send:                 │ │
│  │ │on select │ │  │  → POST /sessions/:id/   │ │
│  │ │on delete │ │  │    chat                  │ │
│  │ │on new    │ │  │  → push assistant msg    │ │
│  │ └──────────┘ │  │                          │ │
│  └──────────────┘  └──────────────────────────┘ │
│       ↕                                        │
│  useSessions composable                        │
│   - sessions: Ref<SessionListItem[]>           │
│   - activeId: Ref<string | null>              │
│   - fetchSessions(), createSession(),          │
│     deleteSession(id)                          │
│   - calls GET/POST/DELETE /sessions            │
└─────────────────────────────────────────────────┘
```

### 3. Session creation: explicit "New Chat" creates a backend session

When the user clicks "New Chat", we call `POST /sessions` to create an empty session, then navigate to it. The first message sent in that session hits `POST /sessions/:id/chat`. Alternative considered: lazy creation (create session on first message). Rejected because the UI needs a session ID immediately for routing, and having an empty session in the list is better UX.

### 4. Backend routes: add the two missing routes

The controller already implements `CreateSession` and `GetSession`. We add them to the route file:

```
POST   /sessions       → CreateSession
GET    /sessions/:id    → GetSession
```

No controller or usecase changes needed. The existing `SendSessionPrompt` already supports auto-creating sessions server-side via `CreatedSession` in the response, but we route `POST /sessions` for the explicit "New Chat" path.

### 5. Chat content persistence: server-side, transparent to UI

The backend `SendSessionPrompt` already persists chat content after each exchange. When switching sessions, the frontend calls `GET /sessions/:id` which returns the session with its `messages` array. The UI loads these messages into the chat view. No additional client-side persistence.

## Risks / Trade-offs

- **Race condition on new session + first message**: User could click "New Chat" and immediately send a message before the `POST /sessions` completes. → Mitigation: Create session synchronously before rendering the chat input; show a loading state on the input until session ID is available.
- **JWT token management**: The backend requires valid JWT for all session endpoints. If the auth flow isn't fully wired on the frontend, requests will fail. → The existing user session token cookie (`session_token`) set by the backend should be included in `credentials: "include"` on fetch calls.
- **No optimistic updates on delete**: Deleting a session sends the request and removes from local state on success. If the network is slow, there's a brief delay. → Acceptable for v1; can add optimistic removal later.

## Open Questions

- How is the JWT token currently stored and sent by the frontend? (Cookie-based? LocalStorage?) This determines how we configure fetch credentials.
- Does the backend's `POST /sessions` need a request body beyond `{ title, model }`? The controller binds `CreateSessionRequest` which expects `session_id`, `title`, `preview`, `model` — the `session_id` is client-generated (UUID).
