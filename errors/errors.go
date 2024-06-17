package errors

const (
    ErrInvalidRequestPayload      = "Invalid request payload"
    ErrMissingRequiredFields      = "Missing required fields: 'type' and 'text'"
    ErrInvalidRecipient           = "Invalid recipient"
    ErrFailedToParseRequestBody   = "Failed to parse request body"
    ErrInvalidPhoneNumberTo       = "Invalid phone number to"
    ErrInvalidPhoneNumberSender   = "Invalid phone number sender"
    ErrNotReadyOrNotAvailable     = "not ready or not available. Please pairing the device"
    ErrInvalidMessageType         = "Invalid message type"
    ErrFailedToSendMessage        = "Failed to send message"
    ErrFailedToMarshalResponse    = "Failed to marshal response"
)
