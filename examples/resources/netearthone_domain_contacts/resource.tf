resource "netearthone_contact" "owner" {
  customer_id    = 12345
  name           = "Jane Smith"
  company        = "N/A"
  email          = "jane@example.com"
  address_line_1 = "123 Main St"
  city           = "London"
  country        = "GB"
  zipcode        = "EC1A 1BB"
  phone_cc       = "44"
  phone          = "2012345678"
}

data "netearthone_domain" "example" {
  domain_name = "example.com"
}

resource "netearthone_domain_contacts" "example" {
  order_id              = data.netearthone_domain.example.order_id
  registrant_contact_id = netearthone_contact.owner.id
  admin_contact_id      = netearthone_contact.owner.id
  tech_contact_id       = netearthone_contact.owner.id
  billing_contact_id    = netearthone_contact.owner.id
}
