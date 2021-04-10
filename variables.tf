variable "env" {
  type = map
  default = {
    RECAPTCHA_SECRET = "6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe" # recaptcha test token, approve all
    site1 = "test@example.com"
  }
}