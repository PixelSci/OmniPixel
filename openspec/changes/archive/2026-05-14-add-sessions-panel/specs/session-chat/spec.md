## ADDED Requirements

### Requirement: Send chat prompt through backend

The system SHALL route all chat prompts through the backend's session chat endpoint instead of calling third-party LLM APIs directly from the browser.

#### Scenario: User sends a message in an existing session
- **WHEN** the user types a message and presses send with an active session
- **THEN** a POST request is sent to `/sessions/:id/chat` with the prompt, provider, model, and api_key
- **AND** the user's message appears immediately in the chat
- **AND** the response is displayed when it arrives

#### Scenario: User sends first message in a new session
- **WHEN** the user sends a message in a session that was just created
- **THEN** the backend processes the prompt and persists the chat content
- **AND** the session title is auto-generated from the first message

#### Scenario: Chat request fails
- **WHEN** the `/sessions/:id/chat` request returns an error
- **THEN** the error message is displayed as an assistant message in the chat
- **AND** the loading state is cleared

#### Scenario: User stops a chat in progress
- **WHEN** the user clicks the stop button during an active chat request
- **THEN** the request is aborted
- **AND** the loading state is cleared
- **AND** any partial response is not displayed

### Requirement: Load persisted messages on session switch

The system SHALL load previously saved chat messages when switching to a session.

#### Scenario: User switches to a session with messages
- **WHEN** the user selects a session that has previous chat history
- **THEN** the `GET /sessions/:id` response includes the full message list
- **AND** all messages (user and assistant) are rendered in the chat area
- **AND** the messages include model attribution for assistant messages

#### Scenario: User switches to an empty session
- **WHEN** the user selects a session with no messages
- **THEN** the chat area shows an empty state ready for input

### Requirement: Chat messages persist after exchange

The system SHALL persist chat content after each prompt-response exchange.

#### Scenario: Successful prompt-response exchange
- **WHEN** the backend completes a chat response
- **THEN** the full message history (including the new exchange) is saved to the session
- **AND** the session's `updated_at` and `last_chat_at` timestamps are refreshed
- **AND** the session list order and grouping reflect the update
