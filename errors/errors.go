package errors

const (
	ErrInvalidRequestPayload    = "Invalid request payload"
	ErrMissingRequiredFields    = "Missing required fields: 'type' and 'text'"
	ErrInvalidRecipient         = "Invalid recipient"
	ErrFailedToParseRequestBody = "Failed to parse request body"
	ErrInvalidPhoneNumberTo     = "Invalid phone number to"
	ErrInvalidPhoneNumberSender = "Invalid phone number sender"
	ErrNotReadyOrNotAvailable   = "not ready or not available. Please pairing the device"
	ErrInvalidMessageType       = "Invalid message type"
	ErrFailedToSendMessage      = "Failed to send message"
	ErrFailedToMarshalResponse  = "Failed to marshal response"
	ErrPhoneNumberRequired      = "Phone number is required"
	ErrFailedToFetchGroups      = "Failed to fetch joined groups"
	ErrInvalidGroupID           = "Invalid group ID"
	ErrInvalidPhoneNumber       = "Invalid phone number"
	ErrFailedToJoinGroup        = "Failed to join group"
	ErrFailedToLeaveGroup       = "Failed to leave group"
	ErrFailedToCreateGroup      = "Failed to create group"
	ErrFailedToGetInviteLink    = "Failed to get group invite link"
)
