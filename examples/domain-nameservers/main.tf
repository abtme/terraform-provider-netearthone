terraform {
  required_providers {
    netearthone = {
      source  = "awxgit/netearthone"
      version = "~> 1.0"
    }
  }
}

provider "netearthone" {
  # auth_userid and api_key can also be set via environment variables:
  #   NETEARTHONE_AUTH_USERID
  #   NETEARTHONE_API_KEY
  auth_userid = 12345
  api_key     = "your-api-key-here"

  # Optional: override the API base URL (e.g. for testing)
  # base_url = "https://test.httpapi.com"
}

# Look up the domain to get its order_id automatically.
data "netearthone_domain" "example" {
  domain_name = "example.com"
}

# Set nameservers using the order_id from the data source.
resource "netearthone_domain_nameservers" "example" {
  order_id = data.netearthone_domain.example.order_id

  nameservers = [
    "ns1.yourdns.com",
    "ns2.yourdns.com",
  ]
}

output "current_nameservers" {
  description = "Nameservers currently assigned to the domain (before any apply)."
  value       = data.netearthone_domain.example.nameservers
}
