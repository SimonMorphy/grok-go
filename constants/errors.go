package constants

const (
	ErrnoSuccess      = 0
	ErrnoUnknownError = -1

	ErrnoMarshalError   = 3
	ErrnoUnmarshalError = 4

	ErrnoInvalidApiKey          = 5
	ErrnoInvalidBaseUrl         = 6
	ErrnoInvalidEndpoint        = 7
	ErrnoInvalidValidationError = 8

	ErrnoBuildRequestError   = 101
	ErrnoInvalidRequestError = 102
	ErrnoInvalidJsonError    = 103
)

var ErrMsg = map[int]string{
	ErrnoSuccess:                "success",
	ErrnoUnknownError:           "unknown error",
	ErrnoMarshalError:           "marshal error",
	ErrnoUnmarshalError:         "unmarshal error",
	ErrnoInvalidApiKey:          "api_key is invalid",
	ErrnoInvalidBaseUrl:         "base_url is invalid",
	ErrnoInvalidEndpoint:        "endpoint is invalid",
	ErrnoInvalidValidationError: "validation failed",
	ErrnoBuildRequestError:      "build request error",
	ErrnoInvalidRequestError:    "invalid request",
	ErrnoInvalidJsonError:       "invalid json",
}
