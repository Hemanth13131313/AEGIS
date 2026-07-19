# Runbook: High Block Rate Spike

## Alert: HighBlockRateSpike
**Severity**: Warning
**Trigger**: Block rate > 20% for 3m

## Possible Causes
- Coordinated prompt injection / attack wave
- Bad policy deployment (too restrictive)
- Scanner misconfiguration (false positives)

## Immediate Actions (< 5 minutes)
1. Check `aegis-scanner` and `aegis-overview` Grafana dashboards
2. Check recent detections: which OWASP category is triggering?
3. Check policy versions: was a new policy deployed recently?

## Diagnostic Steps
- Is the pre-filter catch rate unusually high? Check pre-filter patterns.
- Are both Model A and Model B agreeing? Check ensemble disagreement rate.
- Has the Red Team pass rate regressed?

## Remediation
- If false positive: roll back policy or disable specific pre-filter regex.
- If attack wave: verify system stability, scale scanner replicas if CPU is high.
