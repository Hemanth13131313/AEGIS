# AEGIS Protobuf Contracts

This directory contains the protocol buffer definitions for AEGIS inter-service communication.

## Usage

To generate Go/Python stubs, use `buf generate` or `protoc`.

Example with buf:
```sh
buf generate
```

## Backward-Compatibility Rules

To maintain safe inter-service communication across deployments, observe the following rules:
- **No field renumbering:** Once a field is defined with a number, do not change or reuse it.
- Use `buf breaking` in CI to enforce backward compatibility.
