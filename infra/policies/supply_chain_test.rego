package aegis.policy.supply_chain_test

import rego.v1
import data.aegis.policy.supply_chain

test_allow_signed_attested_no_cves if {
    supply_chain.allow with input as {
        "image_signed": true,
        "sbom_attested": true,
        "critical_cves": 0,
        "base_image": "gcr.io/distroless/static-debian12:nonroot"
    }
}

test_deny_not_signed if {
    "IMAGE_NOT_SIGNED" in supply_chain.deny with input as {
        "image_signed": false,
        "sbom_attested": true,
        "critical_cves": 0,
        "base_image": "gcr.io/distroless/static-debian12:nonroot"
    }
}

test_deny_no_sbom if {
    "SBOM_NOT_ATTESTED" in supply_chain.deny with input as {
        "image_signed": true,
        "sbom_attested": false,
        "critical_cves": 0,
        "base_image": "gcr.io/distroless/static-debian12:nonroot"
    }
}

test_deny_critical_cves if {
    "CRITICAL_CVES_PRESENT" in supply_chain.deny with input as {
        "image_signed": true,
        "sbom_attested": true,
        "critical_cves": 1,
        "base_image": "gcr.io/distroless/static-debian12:nonroot"
    }
}

test_deny_invalid_base_image if {
    "INVALID_BASE_IMAGE" in supply_chain.deny with input as {
        "image_signed": true,
        "sbom_attested": true,
        "critical_cves": 0,
        "base_image": "ubuntu:latest"
    }
}

test_allow_distroless_nonroot if {
    supply_chain.allow with input as {
        "image_signed": true,
        "sbom_attested": true,
        "critical_cves": 0,
        "base_image": "gcr.io/distroless/python3-debian12:nonroot"
    }
}
