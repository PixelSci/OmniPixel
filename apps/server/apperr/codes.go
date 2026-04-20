package apperr

// Error code convention: 5-digit numeric codes in the form MMBBB.
//
//	MM  — module / category
//	      10 = common / system
//	      11 = auth
//	      12 = user
//	      13.. = reserved for future modules
//	BBB — business error within the module (001–999)
//
// Code 0 is reserved for success responses.
const (
	CodeSuccess = 0

	// 10xxx — common / system
	CodeUnknown         = 10000
	CodeInternal        = 10001
	CodeInvalidParam    = 10002
	CodeValidation      = 10003
	CodeNotFound        = 10004
	CodeConflict        = 10005
	CodeTooManyRequests = 10006
	CodeTimeout         = 10007
	CodeUnavailable     = 10008

	// 11xxx — auth
	CodeUnauthorized = 11001
	CodeTokenExpired = 11002
	CodeTokenInvalid = 11003
	CodeForbidden    = 11004

	// 12xxx — user
	CodeUserNotFound   = 12001
	CodeUserExists     = 12002
	CodePasswordWrong  = 12003
	CodeUserDisabled   = 12004
)
