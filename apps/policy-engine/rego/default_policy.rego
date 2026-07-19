package aegis.policy.core

import rego.v1

# Input fields consumed by this policy:
# input.org_id (string): organization identifier
# input.app_id (string): application identifier
# input.auth.valid (boolean): true if the request has a valid auth token
# input.auth.user_id (string): authenticated user identifier
# input.payload.token_count (number): estimated token count of the input
# input.payload.charset_valid (boolean): true if input is valid UTF-8
# input.payload.tool_name (string): name of the tool being invoked (empty if not a tool call)
# input.payload.turn_count (number): number of conversation turns
# input.payload.role_sequence (array[string]): conversation role sequence
# input.allowed_tools (array[string]): org-level allowed tool names
# input.metadata (object): request metadata key-value pairs

default allow := false

# Allow if all conditions pass
allow if {
    input.auth.valid == true
    not exceeds_token_limit
    not invalid_charset
    not invalid_role_sequence
}

# Deny reasons (returned as array of strings)
deny contains "UNAUTHENTICATED" if {
    input.auth.valid != true
}

deny contains "TOKEN_LIMIT_EXCEEDED" if {
    exceeds_token_limit
}

deny contains "INVALID_CHARSET" if {
    invalid_charset
}

deny contains "INVALID_ROLE_SEQUENCE" if {
    invalid_role_sequence
}

deny contains "TOOL_NOT_ALLOWED" if {
    input.payload.tool_name != ""
    not tool_allowed
}

# Helper rules
max_tokens := 8192

exceeds_token_limit if {
    input.payload.token_count > max_tokens
}

invalid_charset if {
    input.payload.charset_valid == false
}

valid_first_roles := {"system", "user"}

invalid_role_sequence if {
    count(input.payload.role_sequence) > 0
    not (input.payload.role_sequence[0] in valid_first_roles)
}

tool_allowed if {
    input.payload.tool_name in input.allowed_tools
}
