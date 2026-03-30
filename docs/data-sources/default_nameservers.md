# netearthone_default_nameservers

Fetches the default nameservers configured for a NetearthOne customer. These are the nameservers automatically assigned when a domain is registered without explicit nameserver values.

**API:** `GET /api/domains/customer-default-ns.json`

## Example

```hcl
data "netearthone_default_nameservers" "customer" {
  customer_id = 12345
}

# Use the customer's defaults for all domains
resource "netearthone_domain_nameservers" "example" {
  order_id    = data.netearthone_domain.example.order_id
  nameservers = data.netearthone_default_nameservers.customer.nameservers
}

output "default_ns" {
  value = data.netearthone_default_nameservers.customer.nameservers
}
```

## Argument Reference

| Argument | Type | Required | Description |
|---|---|---|---|
| `customer_id` | Number | Yes | The NetearthOne customer ID whose default nameservers to retrieve. |

## Attribute Reference

| Attribute | Type | Description |
|---|---|---|
| `nameservers` | List of String | The customer's configured default nameservers. |
