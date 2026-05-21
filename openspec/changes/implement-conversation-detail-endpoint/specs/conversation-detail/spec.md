## ADDED Requirements

### Requirement: Get conversation detail by ID

The system SHALL return a single conversation with its messages when a valid `conversation_id` is provided by an authenticated user.

#### Scenario: Successful retrieval
- **WHEN** an authenticated user requests `GET /api/v1/conversations/:conversation_id` with a valid UUID that belongs to them and `is_visible` is true
- **THEN** the system responds with HTTP 200 and a JSON body containing the conversation metadata (id, title, is_visible, is_archived, created_at, updated_at) and an ordered list of messages, excluding the `user_id` field from the conversation object

#### Scenario: Conversation not found or belongs to another user
- **WHEN** an authenticated user requests a conversation that does not exist, has `is_visible` set to false, or belongs to a different user
- **THEN** the system responds with HTTP 404 and a JSON error message

#### Scenario: Invalid conversation ID format
- **WHEN** an authenticated user requests `GET /api/v1/conversations/:conversation_id` with a value that is not a valid UUID
- **THEN** the system responds with HTTP 400 and a JSON error message

#### Scenario: Unauthenticated request
- **WHEN** a request is made without a valid JWT Bearer token
- **THEN** the JWT middleware responds with HTTP 401 before reaching the handler

### Requirement: Messages returned in chronological order

The system SHALL return messages associated with the conversation ordered by `created_at` in ascending order (oldest first).

#### Scenario: Conversation with multiple messages
- **WHEN** a conversation has 3 messages sent at 10:00, 10:01, and 10:02
- **THEN** the messages array in the response is ordered `[10:00, 10:01, 10:02]`
