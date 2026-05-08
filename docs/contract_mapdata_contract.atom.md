---
id: contract_mapdata_contract
status: STABLE
human_name: UpsilonMapData Contract
type: CONTRACT
version: 1.0
priority: 1
parents:
  - [[shared:contract_upsilon_contract]]
dependents:
  - [[rule_mapdata_grid_standard]]
layer: BUSINESS
tags: [governance, contract, mapdata]
---

# New Atom

## INTENT
Establish the technical standards for geometric data structures and grid representation.

## THE RULE / LOGIC
- **Grid Standard:** Must support rectangular grids with support for height (Z-axis) and obstacle bitmasks.
- **Immutability:** Map data instances should be treated as immutable once generated to prevent race conditions in the engine.
- **Interoperability:** Data formats must be synchronized with `[[upsilontypes]]`.
- **Traceability:** Maintain links between data fields and the tactical rules they support (e.g., Jump limits, Range).

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mapdata_contract]]`
- **Related Atoms:** `[[shared:upsilon_contract]]`

## EXPECTATION
