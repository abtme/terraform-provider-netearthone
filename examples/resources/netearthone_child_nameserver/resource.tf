data "netearthone_domain" "example" {
  domain_name = "example.com"
}

resource "netearthone_child_nameserver" "ns1" {
  order_id     = data.netearthone_domain.example.order_id
  hostname     = "ns1.example.com"
  ip_addresses = ["203.0.113.10"]
}

resource "netearthone_child_nameserver" "ns2" {
  order_id     = data.netearthone_domain.example.order_id
  hostname     = "ns2.example.com"
  ip_addresses = ["203.0.113.11"]
}

# Then use these glue records as nameservers on the domain
resource "netearthone_domain_nameservers" "example" {
  order_id = data.netearthone_domain.example.order_id
  nameservers = [
    netearthone_child_nameserver.ns1.hostname,
    netearthone_child_nameserver.ns2.hostname,
  ]
}
