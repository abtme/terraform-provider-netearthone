# netearthone_child_nameserver

Manages a child nameserver (glue record) for a NetearthOne domain. A child nameserver associates a hostname that lives under your domain (e.g. `ns1.example.com`) with one or more IP addresses.

**Create:** `POST /api/domains/add-cns.json`
**Update hostname:** `POST /api/domains/modify-cns-name.json`
**Update IP:** `POST /api/domains/modify-cns-ip.json`
**Delete:** `POST /api/domains/delete-cns.json`
**Read:** `GET /api/domains/details.json` (options: NsDetails)

## Example

```hcl
data "netearthone_domain" "example" {
  domain_name = "example.com"
}

resource "netearthone_child_nameserver" "ns1" {
  order_id  = data.netearthone_domain.example.order_id
  hostname  = "ns1.example.com"
  ip_addresses = ["203.0.113.10"]
}

resource "netearthone_child_nameserver" "ns2" {
  order_id  = data.netearthone_domain.example.order_id
  hostname  = "ns2.example.com"
  ip_addresses = ["203.0.113.11"]
}

# Then use these glue records as nameservers on the domain
resource "netearthone_domain_nameservers" "example" {
  order_id = data.netearthone_domain.example.order_id
  nameservers = [
    netearthone_child_nameserver.ns1.hostname,
    netearthone_child_nameserver.ns2.hostname,
  ]
}
```

## Argument Reference

| Argument | Type | Required | Description |
|---|---|---|---|
| `order_id` | Number | Yes | NetearthOne order ID of the domain. Changing this forces a new resource. |
| `hostname` | String | Yes | Full child nameserver hostname (e.g. `ns1.example.com`). |
| `ip_addresses` | List of String | Yes | One or more IPv4 or IPv6 addresses for this nameserver. |

## Attribute Reference

| Attribute | Description |
|---|---|
| `id` | Composite identifier: `"<order_id>/<hostname>"`. |
