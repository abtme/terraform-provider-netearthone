data "netearthone_domain_availability" "check" {
  domains = ["mysite", "mybusiness"]
  tlds    = ["com", "net", "co.uk"]
}

output "availability" {
  value = data.netearthone_domain_availability.check.results
  # e.g. { "mysite.com" = "available", "mysite.net" = "regthroughothers", ... }
}

# Only proceed if the domain is available
locals {
  mysite_com_available = data.netearthone_domain_availability.check.results["mysite.com"] == "available"
}
