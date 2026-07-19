package aegis.policy.supply_chain

import rego.v1

# Input fields:
# input.image_digest (string): SHA256 digest of the container image
# input.sbom_attested (boolean): whether a Cosign-attested SBOM exists
# input.image_signed (boolean): whether image is Cosign-signed
# input.critical_cves (integer): count of unfixed CRITICAL CVEs
# input.base_image (string): base image used
# input.model_digest (string): SHA256 of model weights (for AI model supply chain)

default allow := false

allow if {
    input.image_signed == true
    input.sbom_attested == true
    input.critical_cves == 0
    valid_base_image
}

deny contains "IMAGE_NOT_SIGNED" if {
    input.image_signed != true
}

deny contains "SBOM_NOT_ATTESTED" if {
    input.sbom_attested != true
}

deny contains "CRITICAL_CVES_PRESENT" if {
    input.critical_cves > 0
}

deny contains "INVALID_BASE_IMAGE" if {
    not valid_base_image
}

# Approved base images only
approved_bases := {
    "gcr.io/distroless/static-debian12",
    "gcr.io/distroless/static-debian12:nonroot",
    "gcr.io/distroless/python3-debian12",
    "gcr.io/distroless/python3-debian12:nonroot",
}

valid_base_image if {
    input.base_image in approved_bases
}
