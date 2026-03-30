# netearthone_domain

Looks up a NetearthOne domain by name, returning its order ID and current nameservers. The `order_id` output is the key input for all domain management resources.

**API:** `GET /api/domains/details-by-name.json`

## Example

```hcl
data "netearthone_domain" "example" {
  domain_name = "example.com"
}

output "order_id" {
  value = data.netearthone_domain.example.order_id
}

output "current_nameservers" {
  value = data.netearthone_domain.example.nameservers
}
```

## Argument Reference

| Argument | Type | Required | Description |
|---|---|---|---|
| `domain_name` | String | Yes | Fully-qualified domain name to look up (e.g. `"example.com"`). |

## Attribute Reference

| Attribute | Type | Description |
|---|---|---|
| `order_id` | Number | The NetearthOne order ID for this domain registration. |
| `nameservers` | List of String | Current nameservers assigned to the domain. |
