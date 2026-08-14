data "netearthone_domain" "example" {
  domain_name = "example.com"
}

resource "netearthone_domain_privacy" "example" {
  order_id          = data.netearthone_domain.example.order_id
  privacy_protected = true
  reason            = "Managed by Terraform"
}
