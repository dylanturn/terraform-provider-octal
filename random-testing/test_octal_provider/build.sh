clear
cd /Users/dylanturnbull/tmp/terraform-provider-octal
go build
mv /Users/dylanturnbull/tmp/terraform-provider-octal/terraform-provider-octal ~/.terraform.d/plugins/registry.terraform.io/hashicorp/octal/1.0.0/darwin_arm64/terraform-provider-octal
cd /Users/dylanturnbull/tmp/test_octal_provider
# pkill terraform -ls
rm -f .terraform.lock.hcl