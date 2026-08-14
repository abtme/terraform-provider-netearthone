# List all active domains
data "netearthone_domains" "active" {
  status = ["Active"]
}

output "active_domains" {
  value = [for d in data.netearthone_domains.active.domains : d.domain_name]
}

# Filter by TLD and name
data "netearthone_domains" "uk_domains" {
  domain_name   = "example" # substring match
  product_key   = ["dotuk"]
  no_of_records = 100
}

output "uk_domain_list" {
  value = data.netearthone_domains.uk_domains.domains
}
