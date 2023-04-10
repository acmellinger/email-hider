variable "env" {
  type = map(any)
  default = {
    RECAPTCHA_SECRET = "6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe" # recaptcha test token, approve all
    site1            = "test@example.com"
  }
  sensitive = true
}

variable "project" {
  type      = string
  default   = "example-project-123"
  sensitive = true
}
