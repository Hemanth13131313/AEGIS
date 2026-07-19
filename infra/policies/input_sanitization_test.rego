package aegis.policy.input_sanitization_test

import rego.v1

test_allow_clean_input if {
    allow with input as {
        "payload": {"has_pii": false, "has_secrets": false}
    }
}

test_deny_pii if {
    deny["PII_DETECTED"] with input as {
        "payload": {"has_pii": true, "has_secrets": false}
    }
}

test_deny_secrets if {
    deny["SECRET_DETECTED"] with input as {
        "payload": {"has_pii": false, "has_secrets": true}
    }
}
