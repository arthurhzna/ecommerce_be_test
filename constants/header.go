package constants

import "net/textproto"

var (
	XApiKey       = textproto.CanonicalMIMEHeaderKey("api-key")
	Authorization = textproto.CanonicalMIMEHeaderKey("authorization")
)
