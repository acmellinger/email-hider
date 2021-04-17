variable "env" {
  type = map(any)
  default = {
    RECAPTCHA_SECRET = "6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe" # recaptcha test token, approve all
    site1            = "test@example.com"
  }
  sensitive = true
}

variable "cors_origins" {
  type    = list(any)
  default = []
}
