## MODIFIED Requirements

### Requirement: Repository methods use GORM query API

The system SHALL use GORM's query builder methods instead of raw SQL for all repository operations.

#### Scenario: Find user by email
- **WHEN** `UserRepository.FindByEmail` is called with an email address
- **THEN** the repository uses `db.Where("lower(email) = lower(?)", email).First(&user)` and returns the user or a not-found error

#### Scenario: Update user's last login
- **WHEN** `UserRepository.UpdateLastLogin` is called with a user ID
- **THEN** the repository uses `db.Model(&domain.User{}).Where("id = ?", userID).Updates(...)` to set `last_login_at` and `updated_at`

#### Scenario: List conversations by user ID
- **WHEN** `ConversationRepository.ListByUserID` is called with a user UUID
- **THEN** the repository uses `db.Where("user_id = ?", userID).Order("...").Find(&conversations)` and returns the slice

#### Scenario: Find conversation by ID and user ID
- **WHEN** `ConversationRepository.FindByID` is called with a conversation UUID and user UUID
- **THEN** the repository uses `db.Where("id = ? AND user_id = ? AND is_visible = ?", convID, userID, true).First(&conversation)` and returns the conversation or a not-found error

#### Scenario: List messages by conversation ID
- **WHEN** `ConversationRepository.ListMessagesByConversationID` is called with a conversation UUID
- **THEN** the repository uses `db.Where("conversation_id = ?", convID).Order("created_at ASC").Find(&messages)` and returns the slice

#### Scenario: Insert conversation
- **WHEN** `ConversationRepository.Insert` is called with a `*domain.Conversation`
- **THEN** the repository uses `db.Create(conversation)` and returns error if the insert fails

#### Scenario: Insert message
- **WHEN** `ConversationRepository.InsertMessage` is called with a `*domain.Message`
- **THEN** the repository uses `db.Create(message)` and returns error if the insert fails
