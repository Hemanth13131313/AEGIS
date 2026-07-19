package aegis.policy.eu_ai_act_test

import rego.v1
import data.aegis.policy.eu_ai_act

test_compliant_full_implementation if {
    eu_ai_act.allow with input as {
        "risk_management": {"policies_active": true, "last_review_days": 30},
        "data_governance": {"no_raw_payload_storage": true},
        "transparency": {"policy_ref_on_decisions": true},
        "human_oversight": {"disagreement_escalation": true},
        "robustness": {"dual_model_ensemble": true}
    }
}

test_deny_review_overdue if {
    not eu_ai_act.allow with input as {
        "risk_management": {"policies_active": true, "last_review_days": 120},
        "data_governance": {"no_raw_payload_storage": true},
        "transparency": {"policy_ref_on_decisions": true},
        "human_oversight": {"disagreement_escalation": true},
        "robustness": {"dual_model_ensemble": true}
    }
    
    "ART9_REVIEW_OVERDUE" in eu_ai_act.deny with input as {
        "risk_management": {"policies_active": true, "last_review_days": 120},
        "data_governance": {"no_raw_payload_storage": true},
        "transparency": {"policy_ref_on_decisions": true},
        "human_oversight": {"disagreement_escalation": true},
        "robustness": {"dual_model_ensemble": true}
    }
}

test_deny_raw_payload_stored if {
    not eu_ai_act.allow with input as {
        "risk_management": {"policies_active": true, "last_review_days": 30},
        "data_governance": {"no_raw_payload_storage": false},
        "transparency": {"policy_ref_on_decisions": true},
        "human_oversight": {"disagreement_escalation": true},
        "robustness": {"dual_model_ensemble": true}
    }
}

test_deny_no_human_escalation if {
    not eu_ai_act.allow with input as {
        "risk_management": {"policies_active": true, "last_review_days": 30},
        "data_governance": {"no_raw_payload_storage": true},
        "transparency": {"policy_ref_on_decisions": true},
        "human_oversight": {"disagreement_escalation": false},
        "robustness": {"dual_model_ensemble": true}
    }
}

test_deny_single_model if {
    not eu_ai_act.allow with input as {
        "risk_management": {"policies_active": true, "last_review_days": 30},
        "data_governance": {"no_raw_payload_storage": true},
        "transparency": {"policy_ref_on_decisions": true},
        "human_oversight": {"disagreement_escalation": true},
        "robustness": {"dual_model_ensemble": false}
    }
}
