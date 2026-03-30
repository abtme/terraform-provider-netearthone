# netearthone_contact

Manages a NetearthOne registrant contact record. Contact records can be assigned to domain registrations using the [`netearthone_domain_contacts`](domain_contacts.md) resource.

**Create:** `POST /api/contacts/add.json`
**Update:** `POST /api/contacts/modify.json`
**Delete:** `POST /api/contacts/delete.json`
**Read:** `GET /api/contacts/details.json`

> The NetearthOne API may reject deletion of contacts that are currently assigned to active domain orders.

## Example

```hcl
resource "netearthone_contact" "owner" {
  customer_id    = 12345
  type           = "Contact"  # default, can be omitted

  name           = "Jane Smith"
  company        = "N/A"      # use "N/A" for individuals
  email          = "jane@example.com"

  address_line_1 = "123 Main St"
  city           = "London"
  state          = "England"
  country        = "GB"
  zipcode        = "EC1A 1BB"

  phone_cc       = "44"
  phone          = "2012345678"
}

output "contact_id" {
  value = netearthone_contact.owner.id
}
```

## Argument Reference

| Argument | Type | Required | Default | Description |
|---|---|---|---|---|
| `customer_id` | Number | Yes | — | Customer ID under which to create the contact. |
| `type` | String | No | `"Contact"` | Contact type. One of: `Contact`, `UkContact`, `EuContact`, `CnContact`, `CoContact`, `CaContact`, `DeContact`, `EsContact`. |
| `name` | String | Yes | — | Full name (max 255 characters). |
| `company` | String | Yes | — | Company name. Use `"N/A"` for natural persons. |
| `email` | String | Yes | — | Email address. |
| `address_line_1` | String | Yes | — | Primary address line (max 64 characters). |
| `address_line_2` | String | No | `""` | Secondary address line. |
| `address_line_3` | String | No | `""` | Tertiary address line. |
| `city` | String | Yes | — | City (max 64 characters). |
| `state` | String | No | `""` | State or province (max 64 characters). |
| `country` | String | Yes | — | ISO 3166-1 alpha-2 country code (e.g. `"GB"`, `"US"`). |
| `zipcode` | String | Yes | — | Postal/ZIP code. |
| `phone_cc` | String | Yes | — | Telephone country code, digits only (e.g. `"44"` for UK). |
| `phone` | String | Yes | — | Telephone number, 4–12 digits. |
| `fax_cc` | String | No | `""` | Fax country code. |
| `fax` | String | No | `""` | Fax number. |

## Attribute Reference

| Attribute | Description |
|---|---|
| `id` | The contact ID assigned by NetearthOne. |
