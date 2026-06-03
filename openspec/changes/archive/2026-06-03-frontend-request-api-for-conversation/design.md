## Context

Backend routes: `GET /conversations`, `GET /conversations/:id`, `DELETE /conversations/:id`, `POST /conversation`. Frontend API modules currently use "session" naming (`session.ts`) that doesn't match. A dead `conversation.ts` duplicates the same logic with correct naming but is not imported anywhere.

## Goals / Non-Goals

**Goals:**
- Delete dead `conversation.ts`
- Rename `session.ts` → `conversation.ts` with identifier alignment
- Rename `useSessions.ts` → `useConversations.ts` with identifier alignment
- Update all import references in `index.vue` and `HigChatSessionList.vue`
- Zero functional changes — the app should work identically before and after

**Non-Goals:**
- Not renaming `HigChatSessionList.vue` component (cosmetic, would bloat the change)
- Not changing any API endpoint paths or backend code
- Not modifying `streamChat` behavior or SSE parsing
- Not adding new features (pagination, title editing, etc.)

## Decisions

### Rename rather than merge

There are two files with near-identical logic (`session.ts` used, `conversation.ts` dead). We rename the active one rather than switching to the dead one because the dead one lacks `deleteConversation` and has minor differences (imports `ApiError`, different `ChatRequest` field order). The active `session.ts` is known-good.

### Keep component name unchanged

`HigChatSessionList.vue` stays as-is. Renaming it to `HigConversationList.vue` is pure cosmetic churn with no user-facing impact and would require updating the import in `index.vue` as well. We only update the import path and type reference inside the component.

### Identifier mapping

| Old (`session.ts`) | New (`conversation.ts`) |
|---|---|
| `SessionItem` | `ConversationItem` |
| `SessionDetail` | `ConversationDetail` |
| `listSessions()` | `listConversations()` |
| `getSession()` | `getConversation()` |
| `deleteSession()` | `deleteConversation()` |
| `ChatMessage` | (unchanged) |
| `ChatRequest` | (unchanged) |
| `streamChat()` | (unchanged) |

### Composable renaming

| Old (`useSessions.ts`) | New (`useConversations.ts`) |
|---|---|
| `useSessions()` | `useConversations()` |
| `sessions` ref | `conversations` ref |
| `fetchSessions()` | `fetchConversations()` |
| `remove()` (calls `deleteSession`) | calls `deleteConversation` |

## Risks / Trade-offs

- Rename-then-update-imports is a mechanical refactor with near-zero risk
- The app uses `createMemoryHistory()` routing, so there are no URL-based links that could break
- No backend changes, no API contract changes
