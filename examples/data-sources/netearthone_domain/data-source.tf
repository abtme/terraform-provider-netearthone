data "netearthone_domain" "example" {
  domain_name = "example.com"
}

output "order_id" {
  value = data.netearthone_domain.example.order_id
}

output "current_nameservers" {
  value = data.netearthone_domain.example.nameservers
}
