package aegis.policy.tool_allowlist_test

import rego.v1

test_allow_if_in_allowlist if {
    allow with input as {
        "payload": {"tool_name": "search"},
        "allowed_tools": ["search", "calc"],
        "blocked_tools": ["shell"]
    }
}

test_deny_if_in_blocked if {
    deny["TOOL_BLOCKED"] with input as {
        "payload": {"tool_name": "shell"},
        "allowed_tools": ["search", "shell"],
        "blocked_tools": ["shell"]
    }
}

test_deny_if_not_in_allowlist if {
    deny["TOOL_NOT_IN_ALLOWLIST"] with input as {
        "payload": {"tool_name": "eval"},
        "allowed_tools": ["search", "calc"],
        "blocked_tools": []
    }
}
