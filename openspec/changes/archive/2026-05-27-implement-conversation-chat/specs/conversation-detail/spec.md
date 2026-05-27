## ADDED Requirements

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
