package aegis.policy.tool_allowlist

import rego.v1

# Input fields consumed:
# input.payload.tool_name (string)
# input.payload.tool_params (object): key-value parameter map
# input.org_id (string)
# input.allowed_tools (array[string]): configured allowed tools for this org/app
# input.blocked_tools (array[string]): explicitly blocked tools

default allow := false

allow if {
    input.payload.tool_name in input.allowed_tools
    not blocked
}

deny contains "TOOL_BLOCKED" if { blocked }
deny contains "TOOL_NOT_IN_ALLOWLIST" if { not input.payload.tool_name in input.allowed_tools }

blocked if { input.payload.tool_name in input.blocked_tools }
