terraform {
  required_providers {
    netearthone = {
      source = "abtme/netearthone"
    }
  }
}

# auth_userid/api_key can also be set via the NETEARTHONE_AUTH_USERID/
# NETEARTHONE_API_KEY environment variables instead of hardcoding them here.
provider "netearthone" {
  auth_userid = 12345
  api_key     = "your-api-key"
}
