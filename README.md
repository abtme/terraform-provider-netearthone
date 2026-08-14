# Terraform Provider for NetearthOne

A [Terraform](https://www.terraform.io) provider for managing domain
registration resources via the
[NetearthOne HTTP API](https://manage.netearthone.com/kb/answer/776) —
nameservers, WHOIS privacy, registrant contacts, child nameservers (glue
records), and domain lookup/search/availability checks.

## Using the provider

```hcl
terraform {
  required_providers {
    netearthone = {
      source = "abtme/netearthone"
    }
  }
}

provider "netearthone" {
  # auth_userid/api_key can also be set via the NETEARTHONE_AUTH_USERID/
  # NETEARTHONE_API_KEY environment variables
  auth_userid = 12345
  api_key     = "your-api-key"
}

data "netearthone_domain" "example" {
  domain_name = "example.com"
}

resource "netearthone_domain_nameservers" "example" {
  order_id    = data.netearthone_domain.example.order_id
  nameservers = ["ns1.yourdns.com", "ns2.yourdns.com"]
}
```

Full documentation, including the provider schema, every resource/data
source's arguments and attributes, and more usage examples, is published
on the
[Terraform Registry](https://registry.terraform.io/providers/abtme/netearthone/latest/docs).

## Resources and data sources

| Name | Description |
|---|---|
| [`netearthone_domain_nameservers`](docs/resources/domain_nameservers.md) | Sets the nameservers for an existing domain registration. |
| [`netearthone_domain_privacy`](docs/resources/domain_privacy.md) | Manages WHOIS privacy protection. |
| [`netearthone_domain_contacts`](docs/resources/domain_contacts.md) | Assigns registrant/admin/tech/billing contacts to a domain. |
| [`netearthone_child_nameserver`](docs/resources/child_nameserver.md) | Manages a child nameserver (glue record). |
| [`netearthone_contact`](docs/resources/contact.md) | Manages a registrant contact record. |
| [`netearthone_domain`](docs/data-sources/domain.md) (data source) | Looks up a domain by name. |
| [`netearthone_domains`](docs/data-sources/domains.md) (data source) | Lists/filters domain registration orders. |
| [`netearthone_domain_availability`](docs/data-sources/domain_availability.md) (data source) | Checks domain/TLD availability. |
| [`netearthone_contact`](docs/data-sources/contact.md) (data source) | Fetches a contact record by ID. |
| [`netearthone_default_nameservers`](docs/data-sources/default_nameservers.md) (data source) | Fetches a customer's default nameservers. |

See the [API coverage table](docs/index.md#api-coverage) for the full
mapping to NetearthOne API endpoints, including what's not yet implemented.

## Developing the provider

Requires [Go](https://go.dev/) (see `go.mod` for the version) and
[Terraform](https://www.terraform.io/downloads) locally.

```shell
go build ./...
```

### Generating docs

Resource/data-source docs under `docs/` are hand-written (not
tfplugindocs-generated) since they document NetearthOne API endpoint
mappings and behavioral notes that don't come from the provider schema
alone. The `examples/` directory holds the runnable `.tf` snippets shown
in those docs — keep both in sync when changing a resource's behavior.

### Releasing

Pushing a `v*` tag triggers `.github/workflows/release.yml`, which builds
and signs release artifacts with [GoReleaser](https://goreleaser.com/) and
publishes a GitHub Release. The Terraform Registry picks up new versions
automatically once connected.

## License

See [LICENSE](LICENSE).
