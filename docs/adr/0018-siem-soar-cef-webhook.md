# ADR 0018: SIEM/SOAR CEF & Webhook Export

## Status
Accepted

## Context
Enterprise security teams need AEGIS detections integrated into their existing Security Information and Event Management (SIEM) systems and Security Orchestration, Automation, and Response (SOAR) platforms.

## Decision
- Output detections in **CEF 0** (Common Event Format) over UDP/TCP Syslog for SIEM integration.
- Output detections as HMAC-SHA256 signed JSON payloads via Webhooks for SOAR integration.

## Consequences
- Provides out-of-the-box compatibility with major SIEMs like Splunk, Elastic, and IBM QRadar.
- Provides SOAR compatibility with PagerDuty, Tines, Shuffle, etc.
- Requires securing the webhook signing secret.
