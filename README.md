# JSON:API

### TODOs

- many payload
- many payload of pointers
- inferred tags (e.g. no `jsonapi:"..."` tag, assume it is an attribute and infer name from field name)
- ignored fields (e.g. `jsonapi:"-"`)
- errors
- support all native types
- support structs
- support embedded structs
- support nested structs
- support custom types with custom un/marshallers
- jsonapi spec validation
- jsonapi settings (e.g.: spec version, error/warning on document validation, etc.)
- support omitempty tag
- add overflow check and tests for int, uint and float (both value and pointers)