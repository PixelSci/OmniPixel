## ADDED Requirements

### Requirement: List user sessions

The system SHALL fetch and display all chat sessions belonging to the authenticated user, grouped by recency (Today, Yesterday, Previous 7 Days, Previous 30 Days, Older).

#### Scenario: Sessions load on page mount
- **WHEN** the user opens the app
- **THEN** the sidebar displays sessions fetched from `GET /sessions`, grouped by recency

#### Scenario: Empty session list
- **WHEN** the user has no sessions
- **THEN** the sidebar shows only the "New Chat" button with no groups listed

#### Scenario: Session re-fetches after delete
- **WHEN** a session is deleted
- **THEN** the session list refreshes to reflect the removal

### Requirement: Create a new session

The system SHALL allow creating a new empty session via the backend API.

#### Scenario: User clicks "New Chat"
- **WHEN** the user clicks the "New Chat" button
- **THEN** a new session is created via `POST /sessions` with a client-generated UUID
- **AND** the new session ID becomes the active session
- **AND** the chat area shows the empty-state message until a message is sent

#### Scenario: Session creation failure
- **WHEN** the `POST /sessions` request fails
- **THEN** the active session does not change
- **AND** an error toast or message is displayed to the user

### Requirement: Delete a session

The system SHALL allow deleting a session from the list.

#### Scenario: User deletes a session
- **WHEN** the user triggers delete on a session (e.g., right-click or a delete button)
- **THEN** a confirmation dialog appears
- **AND** on confirm, `DELETE /sessions/:id` is called
- **AND** the session is removed from the local list

#### Scenario: Deleting the active session
- **WHEN** the user deletes the currently active session
- **THEN** the active session is cleared
- **AND** the chat area returns to the empty-state ("Select or start a new chat")

#### Scenario: Delete failure
- **WHEN** the `DELETE /sessions/:id` request fails
- **THEN** the session remains in the list
- **AND** an error message is displayed

### Requirement: Switch between sessions

The system SHALL load the chat messages for a session when the user selects it.

#### Scenario: User selects a session
- **WHEN** the user clicks a session in the sidebar
- **THEN** the system fetches `GET /sessions/:id`
- **AND** the returned messages are displayed in the chat area
- **AND** the session is highlighted as active in the sidebar

#### Scenario: Selected session not found
- **WHEN** a `GET /sessions/:id` returns 404
- **THEN** the session is removed from the local list
- **AND** the active session is cleared
- **AND** a notification is shown
