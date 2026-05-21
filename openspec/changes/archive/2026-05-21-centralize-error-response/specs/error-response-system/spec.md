## ADDED Requirements

### Requirement: Error code definition

The system SHALL define numeric error codes grouped by module prefix: 10xxx for Account, 20xxx for Conversation, 30xxx for Message, and 90xxx for system-level errors.

#### Scenario: Error code uniqueness
- **WHEN** a new error code is defined
- **THEN** it SHALL have a unique integer value within its module range and a descriptive constant name of the form `ErrCode<Description>`

### Requirement: APIError construction

The system SHALL provide an `APIError` struct with `code` (int) and `message` (string) fields, and a `New(code, message)` constructor function.

#### Scenario: Build an APIError
- **WHEN** `New(10001, "邮箱或密码错误")` is called
- **THEN** it returns a pointer to an `APIError` with `Code: 10001` and `Message: "邮箱或密码错误"`

#### Scenario: APIError implements the error interface
- **WHEN** an `APIError` is returned as an `error` type
- **THEN** `Error()` returns its `message` string

### Requirement: Business error code maps to standard HTTP status code

The system SHALL provide an `HTTPStatus(code ErrorCode) int` function that maps each business error code to the corresponding HTTP status code.

#### Scenario: Invalid credentials error
- **WHEN** `HTTPStatus(10001)` is called
- **THEN** it returns `401`

#### Scenario: Conversation not found error
- **WHEN** `HTTPStatus(20002)` is called
- **THEN** it returns `404`

#### Scenario: Internal server error
- **WHEN** `HTTPStatus(90000)` is called
- **THEN** it returns `500`

### Requirement: Domain error to APIError mapping

The system SHALL maintain a mapping table that associates each domain error with a pre-built `*APIError`, and provide a `DomainError(c, err)` function that writes the mapped APIError as a JSON HTTP response.

#### Scenario: Domain error mapped successfully
- **WHEN** `DomainError(c, domain.ErrConversationNotFound)` is called
- **THEN** it writes an HTTP response with status 404 and JSON body `{"code": 20002, "message": "会话不存在"}`

#### Scenario: Unmapped domain error falls back
- **WHEN** `DomainError(c, unknownDomainErr)` is called and no mapping exists
- **THEN** it writes HTTP 500 with `{"code": 90000, "message": "服务器内部错误"}`

### Requirement: Generic error response writer

The system SHALL provide a `Write(c, apiErr)` function that writes any `*APIError` as a JSON response, deriving the HTTP status code via `HTTPStatus(apiErr.Code)`.

#### Scenario: Write a generic APIError
- **WHEN** `Write(c, internal.ErrInvalidConvID)` is called
- **THEN** it writes HTTP 400 with `{"code": 20001, "message": "会话 ID 格式错误"}`

### Requirement: Unified JSON error response format

The system SHALL return all error responses in the format `{"code": <int>, "message": "<string>"}`.

#### Scenario: Error response format
- **WHEN** any error response is sent
- **THEN** the JSON body contains exactly the keys `code` (integer) and `message` (string)
