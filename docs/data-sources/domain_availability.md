# netearthone_domain_availability

Checks whether one or more domain name and TLD combinations are available for registration.

**API:** `GET /api/domains/available.json` (via `domaincheck` subdomain)

## Example

```hcl
data "netearthone_domain_availability" "check" {
  domains = ["mysite", "mybusiness"]
  tlds    = ["com", "net", "co.uk"]
}

output "availability" {
  value = data.netearthone_domain_availability.check.results
  # e.g. { "mysite.com" = "available", "mysite.net" = "regthroughothers", ... }
}

# Only proceed if the domain is available
locals {
  mysite_com_available = data.netearthone_domain_availability.check.results["mysite.com"] == "available"
}
```

## Argument Reference

| Argument | Type | Required | Description |
|---|---|---|---|
| `domains` | List of String | Yes | Second-level domain names to check, **without** TLD (e.g. `["mysite", "mybusiness"]`). |
| `tlds` | List of String | Yes | TLD extensions to check against (e.g. `["com", "net", "co.uk"]`). |

## Attribute Reference

| Attribute | Type | Description |
|---|---|---|
| `results` | Map of String | Map of `"domain.tld"` to status. See status values below. |

### Status Values

| Status | Meaning |
|---|---|
| `available` | The domain is available to register. |
| `regthroughus` | Already registered through your NetearthOne reseller account. |
| `regthroughothers` | Registered through another registrar. |
| `unknown` | Status could not be determined. |

> **Note:** The availability check API uses a `domaincheck` subdomain (e.g. `domaincheck.httpapi.com`). The provider derives this automatically from the configured `base_url` by replacing the `api.` prefix. If your setup uses a non-standard URL, this may need adjustment.
