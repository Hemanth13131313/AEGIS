package aegis.policy.input_sanitization

import rego.v1

default allow := false

allow if {
    not contains_pii
    not contains_secrets
}

deny contains "PII_DETECTED" if { contains_pii }
deny contains "SECRET_DETECTED" if { contains_secrets }

contains_pii if {
    input.payload.has_pii == true
}

contains_secrets if {
    input.payload.has_secrets == true
}
