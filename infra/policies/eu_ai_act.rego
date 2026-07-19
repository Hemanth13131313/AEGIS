package aegis.policy.eu_ai_act

import rego.v1

# Input fields:
# input.risk_management.policies_active (boolean)
# input.risk_management.last_review_days (integer): days since last red team run
# input.transparency.policy_ref_on_decisions (boolean)
# input.record_keeping.immutable_audit_log (boolean)
# input.human_oversight.disagreement_escalation (boolean)
# input.robustness.dual_model_ensemble (boolean)
# input.data_governance.no_raw_payload_storage (boolean)

default allow := false
default compliant := false

allow if { compliant }

compliant if {
    article9_risk_management
    article10_data_governance
    article13_transparency
    article14_human_oversight
    article15_robustness
}

# Article 9: Risk management system
article9_risk_management if {
    input.risk_management.policies_active == true
    input.risk_management.last_review_days <= 90  # quarterly minimum
}

# Article 10: Data governance
article10_data_governance if {
    input.data_governance.no_raw_payload_storage == true
}

# Article 13: Transparency
article13_transparency if {
    input.transparency.policy_ref_on_decisions == true
}

# Article 14: Human oversight
article14_human_oversight if {
    input.human_oversight.disagreement_escalation == true
}

# Article 15: Accuracy and robustness
article15_robustness if {
    input.robustness.dual_model_ensemble == true
}

# Non-compliance reasons
deny contains "ART9_NO_ACTIVE_POLICIES" if {
    input.risk_management.policies_active != true
}

deny contains "ART9_REVIEW_OVERDUE" if {
    input.risk_management.last_review_days > 90
}

deny contains "ART10_RAW_PAYLOAD_STORED" if {
    input.data_governance.no_raw_payload_storage != true
}

deny contains "ART13_NO_POLICY_REF" if {
    input.transparency.policy_ref_on_decisions != true
}

deny contains "ART14_NO_HUMAN_ESCALATION" if {
    input.human_oversight.disagreement_escalation != true
}

deny contains "ART15_SINGLE_MODEL_ONLY" if {
    input.robustness.dual_model_ensemble != true
}
