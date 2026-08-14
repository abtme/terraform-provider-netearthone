resource "netearthone_contact" "owner" {
  customer_id = 12345
  type        = "Contact" # default, can be omitted

  name    = "Jane Smith"
  company = "N/A" # use "N/A" for individuals
  email   = "jane@example.com"

  address_line_1 = "123 Main St"
  city           = "London"
  state          = "England"
  country        = "GB"
  zipcode        = "EC1A 1BB"

  phone_cc = "44"
  phone    = "2012345678"
}

output "contact_id" {
  value = netearthone_contact.owner.id
}
