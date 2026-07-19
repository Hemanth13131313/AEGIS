# ADR 0016: Provider Adapter Pattern

## Status
Accepted

## Context
Different AI providers (OpenAI, Anthropic, Gemini, Bedrock) have widely different API schemas for chat completions and generation. Writing custom inspection logic for each provider's schema in the scanner and policy engine would lead to high complexity and tight coupling.

## Decision
Implement an Adapter pattern in the Go Gateway. The adapter interface normalizes all provider-specific HTTP requests and responses to a standard `CanonicalRequest` and `CanonicalResponse` structure before AEGIS performs inspection.

## Consequences
- New providers require only a new adapter implementation.
- Existing scanner and policy logic remain completely unchanged and agnostic of the underlying provider.
- Raw provider-specific formats are no longer passed to the scanner, ensuring consistent security evaluation.
