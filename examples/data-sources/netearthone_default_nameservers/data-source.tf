data "netearthone_default_nameservers" "customer" {
  customer_id = 12345
}

data "netearthone_domain" "example" {
  domain_name = "example.com"
}

# Use the customer's defaults for all domains
resource "netearthone_domain_nameservers" "example" {
  order_id    = data.netearthone_domain.example.order_id
  nameservers = data.netearthone_default_nameservers.customer.nameservers
}

output "default_ns" {
  value = data.netearthone_default_nameservers.customer.nameservers
}
