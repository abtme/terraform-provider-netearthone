# netearthone_domain_nameservers

Sets the nameservers for an existing NetearthOne domain registration order.

**API:** `POST /api/domains/modify-ns.json`
**Read:** `GET /api/domains/details.json` (options: NsDetails)

> Removing this resource from Terraform state does **not** reset the nameservers — it simply stops managing them.

## Example

```hcl
data "netearthone_domain" "example" {
  domain_name = "example.com"
}

resource "netearthone_domain_nameservers" "example" {
  order_id = data.netearthone_domain.example.order_id

  nameservers = [
    "ns1.yourdns.com",
    "ns2.yourdns.com",
  ]
}
```

## Argument Reference

| Argument | Type | Required | Description |
|---|---|---|---|
| `order_id` | Number | Yes | NetearthOne order ID of the domain. Changing this forces a new resource. |
| `nameservers` | List of String | Yes | Nameservers to assign (minimum 2, maximum 13). |

## Attribute Reference

| Attribute | Description |
|---|---|
| `id` | The order ID (as a string). |
