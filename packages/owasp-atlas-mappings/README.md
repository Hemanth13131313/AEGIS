# OWASP and MITRE ATLAS Mappings

This package contains the taxonomy mapping data for the AEGIS project. It maps internal detection categories to industry-standard taxonomies:
- **OWASP LLM Top 10**
- **MITRE ATLAS**

## Structure

The `mappings.json` file contains:
- `owasp_llm_top10`: Definitions for the OWASP LLM Top 10 vulnerabilities.
- `mitre_atlas_techniques`: Definitions for MITRE ATLAS techniques, along with `owasp_refs`.
- `detection_to_taxonomy_map`: Maps an AEGIS internal detection category (e.g., `prompt_injection`) to the corresponding OWASP and ATLAS identifiers.

## Usage

Services should import and use this mapping to attach taxonomy context to DetectionEvents.
