# conversation-detail Specification

## Purpose
TBD - created by archiving change implement-conversation-detail-endpoint. Update Purpose after archive.
## Requirements
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

### Requirement: Chat request DTO

The system SHALL define a `ChatRequest` struct with `ConversationID` (nullable UUID pointer), `Message` (string), and `ModelID` (UUID string) fields.

#### Scenario: New conversation request
- **WHEN** a chat request has `ConversationID` set to nil
- **THEN** the system creates a new conversation before processing the message

#### Scenario: Continue conversation request
- **WHEN** a chat request has a valid `ConversationID`
- **THEN** the system appends to the existing conversation

### Requirement: StreamWriter interface

The system SHALL define a `StreamWriter` interface in the domain layer with `WriteToken(token string) error` and `WriteDone(conversationID, messageID uuid.UUID) error` methods.

#### Scenario: Streaming tokens
- **WHEN** the UseCase calls `WriteToken("你好")`
- **THEN** the implementation formats and writes the token to the SSE client

#### Scenario: Signaling completion
- **WHEN** the UseCase calls `WriteDone(convID, msgID)`
- **THEN** the implementation writes the final SSE event with done flag and IDs

### Requirement: AIProvider interface

The system SHALL define an `AIProvider` interface in the domain layer with `ChatStream(messages []AIChatMessage, modelID string) (<-chan AIStreamChunk, error)` method, where `AIChatMessage` has `Role` and `Content` fields, and `AIStreamChunk` has `Token` and `Done` fields.

#### Scenario: AI provider streams response
- **WHEN** `ChatStream` is called with message history
- **THEN** it returns a receive-only channel delivering chunks until `Done` is true

### Requirement: Conversation repository insert methods

The `ConversationRepository` interface SHALL include `Insert(conversation *Conversation) error` and `InsertMessage(message *Message) error` methods.

#### Scenario: Insert conversation
- **WHEN** `Insert` is called with a `*Conversation`
- **THEN** the implementation persists the record and returns the error

#### Scenario: Insert message
- **WHEN** `InsertMessage` is called with a `*Message`
- **THEN** the implementation persists the record and returns the error

