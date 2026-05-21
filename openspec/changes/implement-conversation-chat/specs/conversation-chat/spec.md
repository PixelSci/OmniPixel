## ADDED Requirements

### Requirement: Stream chat endpoint

The system SHALL provide `POST /api/v1/conversation` that accepts a JSON body with `conversation_id` (nullable UUID), `message` (string), and `model_id` (UUID string), and returns an SSE stream of AI tokens.

#### Scenario: New conversation
- **WHEN** the request contains `conversation_id: null` and a valid message
- **THEN** the system creates a new Conversation with title set to the first 60 characters of the message, stores the user message, calls the AI provider with the message, streams AI tokens via SSE, stores the complete AI response, and ends the stream with `{"done": true, "conversation_id": "...", "message_id": "..."}`

#### Scenario: Continue existing conversation
- **WHEN** the request contains a valid `conversation_id` belonging to the authenticated user
- **THEN** the system loads conversation history, appends the user message, calls AI with full context, streams tokens, and saves the AI response

#### Scenario: Conversation not owned by user
- **WHEN** the request's `conversation_id` belongs to a different user
- **THEN** the system returns HTTP 404 via the error response system

#### Scenario: AI provider fails
- **WHEN** the AI API returns an error
- **THEN** the system streams `{"error": "..."}` and closes the SSE connection

### Requirement: SSE stream format

The system SHALL stream each AI response token as an SSE event with JSON payload `{"token": "<text>"}` and close with `{"done": true, "conversation_id": "<uuid>", "message_id": "<uuid>"}`.

#### Scenario: Multi-token response
- **WHEN** the AI returns "你好世界" as three tokens
- **THEN** the SSE stream contains `data: {"token":"你"}\n\n`, `data: {"token":"好"}\n\n`, `data: {"token":"世界"}\n\n`, `data: {"done":true,...}\n\n`

### Requirement: AIProvider abstraction

The system SHALL define an `AIProvider` interface in the domain layer with a `ChatStream(messages []AIChatMessage, modelID string) (<-chan AIStreamChunk, error)` method.

#### Scenario: Provider accepts chat history
- **WHEN** `ChatStream` is called with a slice of `AIChatMessage` containing `[{"system": "..."}, {"user": "..."}, {"assistant": "..."}, {"user": "new"}]`
- **THEN** the provider sends the full context to the AI API and returns a channel of streaming chunks

### Requirement: StreamWriter abstraction

The system SHALL define a `StreamWriter` interface in the domain layer with `WriteToken(token string) error` and `WriteDone(conversationID, messageID uuid.UUID) error` methods, implemented by the controller layer as an SSE adapter.

#### Scenario: UseCase writes through StreamWriter
- **WHEN** the UseCase calls `streamWriter.WriteToken("你好")` then `streamWriter.WriteDone(convID, msgID)`
- **THEN** the controller's SSE adapter formats and writes the SSE data to the Fiber response
