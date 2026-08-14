data "netearthone_contact" "owner" {
  contact_id = 67890
}

output "contact_name" {
  value = data.netearthone_contact.owner.name
}

output "contact_email" {
  value = data.netearthone_contact.owner.email
}
