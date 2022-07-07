provider "octal" {
  # example configuration here
}



resource "octal_cert_manager" "example" {
  name = "new-asdf"
  
  
  controller {}
  cainjector {}
  webhook {}
}