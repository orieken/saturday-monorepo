#!/usr/bin/env python3
"""Deterministic spec-critic checks for a by-hand spec package (prompt 02).

Usage: python3 docs/specs/spec-critic.py docs/specs/<id> features/<id>
Exits 1 if any check fails. Candidate for the Go spec-critic in deliver-spec.
"""
import json, re, sys, pathlib

IMPLEMENTATION = re.compile(r"\b(click\w*|selector|url|https?|endpoint|button|field|get|post|put|delete|json|css|xpath|api/)\b", re.I)
UNOBSERVABLE = re.compile(r"\b(the system|process(es|ed)?|internal|state is|is stored in|database)\b", re.I)

def load(spec, name):
    return json.loads((spec / f"{name}.json").read_text())

def check_example_map(emap, errors):
    for r in emap["rules"]:
        if not any(e["kind"] == "negative" for e in r["examples"]):
            errors.append(f"{r['id']}: no negative example")

def check_traceability(emap, dmodel, crit, errors):
    rules = {r["id"] for r in emap["rules"]}
    events = {e["name"] for e in dmodel["events"]}
    covered = set()
    for c in crit["criteria"]:
        if c["rule"] not in rules:
            errors.append(f"{c['id']}: traces to unknown rule {c['rule']}")
        if not c["events"] or not set(c["events"]) <= events:
            errors.append(f"{c['id']}: events {c['events']} not all in the domain model")
        covered.add(c["rule"])
        if UNOBSERVABLE.search(c["then"]):
            errors.append(f"{c['id']}: Then is not an observable outcome: {c['then']!r}")
    for r in sorted(rules - covered):
        errors.append(f"{r}: no criterion covers this rule")

def check_features(features, crit, terms, errors):
    ids = {c["id"] for c in crit["criteria"]}
    tagged = set()
    for f in sorted(features.glob("*.feature")):
        for n, line in enumerate(f.read_text().splitlines(), 1):
            s = line.strip()
            tagged |= set(re.findall(r"@(AC\d+)", s))
            if re.match(r"(Given|When|Then|And|But)\b", s):
                if IMPLEMENTATION.search(s):
                    errors.append(f"{f.name}:{n}: implementation detail: {s!r}")
                for quoted in re.findall(r'"([^"]+)"', s):
                    if quoted not in terms:
                        errors.append(f"{f.name}:{n}: quoted name not in the seed: {quoted!r}")
                for word in re.findall(r"\b[A-Z][a-z]+(?: [A-Z][a-z]+)*\b", re.sub(r'"[^"]*"', "", s)):
                    if word.split()[0] in ("Given", "When", "Then", "And", "But"):
                        word = " ".join(word.split()[1:])
                    if word and word not in terms:
                        errors.append(f"{f.name}:{n}: term not in the dictionary or seed: {word!r}")
            if s.startswith("Scenario") and not s.startswith("Scenario Outline"):
                pass
        for block in re.split(r"\n\s*\n", f.read_text()):
            if "Scenario:" in block:
                tags = " ".join(l for l in block.splitlines() if l.strip().startswith("@"))
                for need in ("@AC", "@level-", "@runner-"):
                    if need not in tags:
                        errors.append(f"{f.name}: scenario missing {need} tag: {block.strip().splitlines()[-1][:60]!r}")
    for i in sorted(ids - tagged, key=lambda x: int(x[2:])):
        errors.append(f"{i}: no scenario is tagged with it")

def main():
    spec, features = pathlib.Path(sys.argv[1]), pathlib.Path(sys.argv[2])
    emap, dmodel, crit = load(spec, "example-map"), load(spec, "domain-model"), load(spec, "acceptance-criteria")
    seed = crit["seed"]
    terms = {"Aurelis", "Elara", "Borin", "Expedition", "Reagent", "I"} | {s["name"] for s in seed["reagentSites"]} | {e["name"] for e in seed["seededExpeditions"]} | {x for a in seed["adventurers"] for x in a["savedExpeditions"]}
    errors = []
    check_example_map(emap, errors)
    check_traceability(emap, dmodel, crit, errors)
    if features.exists():
        check_features(features, crit, terms, errors)
    for e in errors:
        print(f"FAIL {e}")
    print(f"spec-critic: {len(errors)} finding(s)")
    sys.exit(1 if errors else 0)

if __name__ == "__main__":
    main()
