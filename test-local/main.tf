terraform {
  required_providers {
    netearthone = {
      source = "awxgit/netearthone"
    }
  }
}

provider "netearthone" {
  auth_userid = var.auth_userid
  api_key     = var.api_key
  base_url    = var.base_url
}

# Look up domain by name — gives us the order_id automatically.
data "netearthone_domain" "this" {
  domain_name = var.domain_name
}

resource "netearthone_domain_nameservers" "this" {
  order_id    = data.netearthone_domain.this.order_id
  nameservers = var.nameservers
}

output "order_id" {
  description = "The order ID for the domain."
  value       = data.netearthone_domain.this.order_id
}

output "nameservers_before" {
  description = "Nameservers on the domain at plan time (from the data source read)."
  value       = data.netearthone_domain.this.nameservers
}
