plugin "testing" {
  enabled = true
}

plugin "opentofu" {
  enabled = false
}

config {
  varfile = ["from_config.tfvars"]
}
