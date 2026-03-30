# netearthone_domains

Lists and filters NetearthOne domain registration orders. Useful for discovering order IDs or reporting on domain status across your account.

**API:** `GET /api/domains/search.json`

## Example — list all active domains

```hcl
data "netearthone_domains" "active" {
  status = ["Active"]
}

output "active_domains" {
  value = [for d in data.netearthone_domains.active.domains : d.domain_name]
}
```

## Example — filter by TLD and name

```hcl
data "netearthone_domains" "uk_domains" {
  domain_name = "example"       # substring match
  product_key = ["dotuk"]
  no_of_records = 100
}

output "uk_domain_list" {
  value = data.netearthone_domains.uk_domains.domains
}
```

## Argument Reference

All arguments are optional.

| Argument | Type | Description |
|---|---|---|
| `domain_name` | String | Filter by domain name substring. |
| `status` | List of String | Filter by order status. Valid values: `InActive`, `Active`, `Suspended`, `Pending Delete Restorable`, `Deleted`, `Archived`. |
| `product_key` | List of String | Filter by TLD product key (e.g. `"dotcom"`, `"dotnet"`, `"dotuk"`). |
| `no_of_records` | Number | Records per page (default `50`, max `500`). |
| `page_no` | Number | Page number for pagination (default `1`). |

## Attribute Reference

| Attribute | Type | Description |
|---|---|---|
| `total_records` | Number | Total number of matching records across all pages. |
| `domains` | List of Object | Matching domain orders (see below). |

### `domains` object attributes

| Attribute | Type | Description |
|---|---|---|
| `order_id` | String | Order ID of the domain registration. |
| `domain_name` | String | Fully-qualified domain name. |
| `status` | String | Current order status. |
| `product_key` | String | TLD product key. |
| `expiry_date` | String | Domain expiry timestamp (Unix epoch string). |
| `creation_date` | String | Domain creation timestamp (Unix epoch string). |
