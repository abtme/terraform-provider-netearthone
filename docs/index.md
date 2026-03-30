# NetearthOne Provider

The NetearthOne provider manages domain registration resources via the [NetearthOne HTTP API](https://manage.netearthone.com/kb/answer/776).

It covers nameserver management, WHOIS privacy, contact records, child nameservers (glue records), and domain lookups.

## API Coverage

### Implemented

| NetearthOne API Endpoint | Method | Terraform Resource / Data Source |
|---|---|---|
| `/api/domains/modify-ns.json` | POST | `netearthone_domain_nameservers` (resource) |
| `/api/domains/details-by-name.json` | GET | `netearthone_domain` (data source) |
| `/api/domains/details.json` | GET | Used internally by multiple resources for Read |
| `/api/domains/modify-privacy-protection.json` | POST | `netearthone_domain_privacy` (resource) |
| `/api/domains/modify-contact.json` | POST | `netearthone_domain_contacts` (resource) |
| `/api/domains/add-cns.json` | POST | `netearthone_child_nameserver` (resource — Create) |
| `/api/domains/modify-cns-name.json` | POST | `netearthone_child_nameserver` (resource — Update hostname) |
| `/api/domains/modify-cns-ip.json` | POST | `netearthone_child_nameserver` (resource — Update IP) |
| `/api/domains/delete-cns.json` | POST | `netearthone_child_nameserver` (resource — Delete) |
| `/api/domains/search.json` | GET | `netearthone_domains` (data source) |
| `/api/domains/available.json` | GET | `netearthone_domain_availability` (data source) |
| `/api/domains/customer-default-ns.json` | GET | `netearthone_default_nameservers` (data source) |
| `/api/contacts/add.json` | POST | `netearthone_contact` (resource — Create) |
| `/api/contacts/modify.json` | POST | `netearthone_contact` (resource — Update) |
| `/api/contacts/delete.json` | POST | `netearthone_contact` (resource — Delete) |
| `/api/contacts/details.json` | GET | `netearthone_contact` (data source) and resource Read |

### Not Yet Implemented

| NetearthOne API Endpoint | Notes |
|---|---|
| `/api/domains/restore.json` | Restore a deleted domain — one-shot action, not well-suited to Terraform state |
| `/api/domains/modify-auth-code.json` | Modify EPP/transfer auth code |
| `/api/domains/orderid.json` | Get order ID by domain name (covered by `netearthone_domain` data source instead) |
| `/api/domains/tel/modify-whois-pref.json` | `.TEL` domain WHOIS preference — TLD-specific |
| `/api/contacts/search.json` | Search contacts — client method exists but no data source yet |
| `/api/contacts/default.json` | Get customer's default contacts |
| `/api/contacts/set-details.json` | Deprecated by NetearthOne |
| `/api/contacts/coop/add-sponsor.json` | `.COOP` sponsor management — TLD-specific |
| `/api/contacts/sponsors.json` | Get `.COOP` sponsors — TLD-specific |
| `/api/domains/available.json` (domaincheck host) | Uses a different API subdomain; the base URL substitution may need adjustment |

---

## Authentication

All API calls require an `auth-userid` and `api-key`. These can be set in the provider block or via environment variables.

| Environment Variable | Description |
|---|---|
| `NETEARTHONE_AUTH_USERID` | Your NetearthOne auth user ID |
| `NETEARTHONE_API_KEY` | Your NetearthOne API key |
| `NETEARTHONE_BASE_URL` | Override the API base URL (default: `https://api.netearthone.com`) |

## Provider Configuration

```hcl
terraform {
  required_providers {
    netearthone = {
      source  = "awxgit/netearthone"
      version = "~> 1.0"
    }
  }
}

provider "netearthone" {
  auth_userid = 12345        # or NETEARTHONE_AUTH_USERID env var
  api_key     = "your-key"  # or NETEARTHONE_API_KEY env var

  # Optional: use the sandbox for testing
  # base_url = "https://test.httpapi.com"
}
```

## Resources

- [netearthone_contact](resources/contact.md)
- [netearthone_child_nameserver](resources/child_nameserver.md)
- [netearthone_domain_contacts](resources/domain_contacts.md)
- [netearthone_domain_nameservers](resources/domain_nameservers.md)
- [netearthone_domain_privacy](resources/domain_privacy.md)

## Data Sources

- [netearthone_contact](data-sources/contact.md)
- [netearthone_default_nameservers](data-sources/default_nameservers.md)
- [netearthone_domain](data-sources/domain.md)
- [netearthone_domain_availability](data-sources/domain_availability.md)
- [netearthone_domains](data-sources/domains.md)
