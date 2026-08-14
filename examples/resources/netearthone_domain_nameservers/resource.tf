data "netearthone_domain" "example" {
  domain_name = "example.com"
}

resource "netearthone_domain_nameservers" "example" {
  order_id = data.netearthone_domain.example.order_id

  nameservers = [
    "ns1.yourdns.com",
    "ns2.yourdns.com",
  ]
}
