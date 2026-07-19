package aegis.policy.core_test

import rego.v1
import data.aegis.policy.core.allow
import data.aegis.policy.core.deny

# --- Allow path tests ---

test_allow_valid_authenticated_request if {
    allow with input as {
        "auth": {"valid": true, "user_id": "user-123"},
        "payload": {
            "token_count": 100,
            "charset_valid": true,
            "tool_name": "",
            "turn_count": 1,
            "role_sequence": ["user"]
        },
        "allowed_tools": []
    }
}

test_allow_with_allowed_tool if {
    allow with input as {
        "auth": {"valid": true, "user_id": "user-123"},
        "payload": {
            "token_count": 100,
            "charset_valid": true,
            "tool_name": "web_search",
            "turn_count": 1,
            "role_sequence": ["user"]
        },
        "allowed_tools": ["web_search", "calculator"]
    }
}

# --- Deny path tests ---

test_deny_unauthenticated if {
    deny["UNAUTHENTICATED"] with input as {
        "auth": {"valid": false, "user_id": ""},
        "payload": {"token_count": 100, "charset_valid": true, "tool_name": "", "turn_count": 1, "role_sequence": ["user"]},
        "allowed_tools": []
    }
}

test_deny_token_limit_exceeded if {
    deny["TOKEN_LIMIT_EXCEEDED"] with input as {
        "auth": {"valid": true, "user_id": "user-123"},
        "payload": {"token_count": 9000, "charset_valid": true, "tool_name": "", "turn_count": 1, "role_sequence": ["user"]},
        "allowed_tools": []
    }
}

test_deny_invalid_charset if {
    deny["INVALID_CHARSET"] with input as {
        "auth": {"valid": true, "user_id": "user-123"},
        "payload": {"token_count": 100, "charset_valid": false, "tool_name": "", "turn_count": 1, "role_sequence": ["user"]},
        "allowed_tools": []
    }
}

test_deny_tool_not_allowed if {
    deny["TOOL_NOT_ALLOWED"] with input as {
        "auth": {"valid": true, "user_id": "user-123"},
        "payload": {"token_count": 100, "charset_valid": true, "tool_name": "shell_exec", "turn_count": 1, "role_sequence": ["user"]},
        "allowed_tools": ["web_search"]
    }
}

test_deny_invalid_role_sequence if {
    deny["INVALID_ROLE_SEQUENCE"] with input as {
        "auth": {"valid": true, "user_id": "user-123"},
        "payload": {"token_count": 100, "charset_valid": true, "tool_name": "", "turn_count": 1, "role_sequence": ["assistant", "user"]},
        "allowed_tools": []
    }
}
