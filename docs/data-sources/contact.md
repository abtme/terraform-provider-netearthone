# netearthone_contact

Fetches the details of a NetearthOne contact record by its contact ID.

**API:** `GET /api/contacts/details.json`

## Example

```hcl
data "netearthone_contact" "owner" {
  contact_id = 67890
}

output "contact_name" {
  value = data.netearthone_contact.owner.name
}

output "contact_email" {
  value = data.netearthone_contact.owner.email
}
```

## Argument Reference

| Argument | Type | Required | Description |
|---|---|---|---|
| `contact_id` | Number | Yes | The NetearthOne contact ID to look up. |

## Attribute Reference

| Attribute | Type | Description |
|---|---|---|
| `name` | String | Full name of the contact. |
| `company` | String | Company name. |
| `type` | String | Contact type (e.g. `Contact`, `UkContact`). |
| `email` | String | Email address. |
| `address_line_1` | String | Primary address line. |
| `address_line_2` | String | Secondary address line. |
| `address_line_3` | String | Tertiary address line. |
| `city` | String | City. |
| `state` | String | State or province. |
| `country` | String | ISO 3166-1 alpha-2 country code. |
| `zipcode` | String | Postal/ZIP code. |
| `phone_cc` | String | Telephone country code. |
| `phone` | String | Telephone number. |
| `status` | String | Current status of the contact (e.g. `Active`, `InActive`). |
