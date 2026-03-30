variable "auth_userid" {
  description = "NetearthOne auth-userid."
  type        = number
}

variable "api_key" {
  description = "NetearthOne API key."
  type        = string
  sensitive   = true
}

variable "base_url" {
  description = "NetearthOne API base URL. Use https://test.httpapi.com for sandbox testing."
  type        = string
  default     = "https://api.netearthone.com"
}

variable "domain_name" {
  description = "The domain name to manage nameservers for."
  type        = string
}

variable "nameservers" {
  description = "List of nameservers to assign to the domain."
  type        = list(string)
}
