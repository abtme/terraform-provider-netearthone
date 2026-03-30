# netearthone_domain_contacts

Assigns registrant, admin, tech, and billing contacts to a NetearthOne domain registration order.

**API:** `POST /api/domains/modify-contact.json`
**Read:** `GET /api/domains/details.json` (options: ContactIds)

Use the [`netearthone_contact`](contact.md) resource to create the contact records themselves.

> Destroying this resource does **not** unassign the contacts — it stops Terraform from managing the assignments.

## Example

```hcl
resource "netearthone_contact" "owner" {
  customer_id    = 12345
  name           = "Jane Smith"
  company        = "N/A"
  email          = "jane@example.com"
  address_line_1 = "123 Main St"
  city           = "London"
  country        = "GB"
  zipcode        = "EC1A 1BB"
  phone_cc       = "44"
  phone          = "2012345678"
}

data "netearthone_domain" "example" {
  domain_name = "example.com"
}

resource "netearthone_domain_contacts" "example" {
  order_id              = data.netearthone_domain.example.order_id
  registrant_contact_id = netearthone_contact.owner.id
  admin_contact_id      = netearthone_contact.owner.id
  tech_contact_id       = netearthone_contact.owner.id
  billing_contact_id    = netearthone_contact.owner.id
}
```

## Argument Reference

| Argument | Type | Required | Description |
|---|---|---|---|
| `order_id` | Number | Yes | NetearthOne order ID of the domain. Changing this forces a new resource. |
| `registrant_contact_id` | Number | Yes | Contact ID to set as the domain registrant (owner). |
| `admin_contact_id` | Number | Yes | Contact ID for the administrative contact. |
| `tech_contact_id` | Number | Yes | Contact ID for the technical contact. |
| `billing_contact_id` | Number | Yes | Contact ID for the billing contact. |

## Attribute Reference

| Attribute | Description |
|---|---|
| `id` | The order ID (as a string). |
