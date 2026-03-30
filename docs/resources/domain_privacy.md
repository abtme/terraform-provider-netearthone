# netearthone_domain_privacy

Manages WHOIS privacy protection for a NetearthOne domain registration order.

**API:** `POST /api/domains/modify-privacy-protection.json`
**Read:** `GET /api/domains/details.json` (options: DomainStatus)

> Destroying this resource does **not** automatically disable privacy. Set `privacy_protected = false` before destroying if you want to explicitly disable it first.

## Example

```hcl
data "netearthone_domain" "example" {
  domain_name = "example.com"
}

resource "netearthone_domain_privacy" "example" {
  order_id          = data.netearthone_domain.example.order_id
  privacy_protected = true
  reason            = "Managed by Terraform"
}
```

## Argument Reference

| Argument | Type | Required | Description |
|---|---|---|---|
| `order_id` | Number | Yes | NetearthOne order ID of the domain. Changing this forces a new resource. |
| `privacy_protected` | Boolean | Yes | `true` to enable WHOIS privacy, `false` to disable. |
| `reason` | String | Yes | Reason for the change (required by the API). |

## Attribute Reference

| Attribute | Description |
|---|---|
| `id` | The order ID (as a string). |
