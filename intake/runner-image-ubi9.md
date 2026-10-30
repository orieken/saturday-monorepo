---
intake: 1
id: runner-image-ubi9
type: chore
plan: deliver-feature
status: draft
priority: 2
created: 2026-10-29
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: [ui-suite-real-chrome]
budget:
  max-attempts: 2
  max-usd: 25
  max-hours: 3
scope:
  paths: [containers/saturday-runner/**, docs/**]
---

# Add the UBI9 variant of the runner image

## Request

> Take both: UBI9 for teams that want a more secure container, and Debian for everyday users.

## Summary

The runner image's second supported variant, `ubi9`, with the same entrypoint, viewer and
`runner-check` as the Debian variant, its packages pinned, and WebGL working. Designed in
`docs/runner-image.md`, "The UBI9 variant".

## Acceptance Criteria

- [ ] Given `containers/saturday-runner/Containerfile`, then a `ubi9` target starts from
      `registry.access.redhat.com/ubi9/ubi` pinned by digest, and installs the browsers, display,
      viewer and Mesa packages `base` has, from UBI, CentOS Stream 9 and EPEL 9 with signature
      checks on every repository. It contains no `gpgcheck=0`, `--nogpgcheck` or `rpm -e --nodeps`.
- [ ] Given the Stream repositories, then they exclude `openssl*`; given the UBI repositories, then
      they exclude `mesa*`, so every Mesa package comes from Stream at one version.
- [ ] Given `containers/saturday-runner/rpm.lock`, then it lists every installed package at its exact
      version, `lock-packages.sh --ubi9` regenerates it, and the build fails if a locked version is
      unavailable.
- [ ] Given the EPEL packages (`x11vnc`, `openbox`, `novnc`, `python3-websockify`), then they are
      kept in the build context as RPMs, checked against `SHA256SUMS` and EPEL's signing key, and
      installed from there rather than from the EPEL mirror.
- [ ] Given the `ubi9` image, then `runner-check` passes in Chrome and Edge headed, with a non-null
      `webgl2` renderer recorded in the evidence, and the image runs as uid 1000.
- [ ] Given `scripts/ui-suite-in-chrome.sh --ubi9`, then the suite runs in the `ubi9` image exactly as
      it does in `base`.
- [ ] Given `docs/runner-image.md`, then it says how to choose a variant and how a browser bump
      updates both locks.
- [ ] Given the change is complete, then the Debian variant is unchanged.

## Out of Scope

- A UBI-only package set: not possible with a virtual display (see the design).
- testing-with-certs adopting this target.

## Domain Language

- No new terms introduced.

## Non-Functional Requirements

- None beyond `base`'s.

## Trust Boundaries

- Packages from UBI, Stream and EPEL are accepted only with valid signatures; vendored RPMs also
  match recorded SHA-256s.

## Test Approach

`runner-check` and the UI suite in the `ubi9` image.

## Open Questions

- None

## Infrastructure / Deploy Notes

- None
