terraform {
  required_providers {
    foxy = {
      source = "67bricks.com/terraform/foxycart"
    }
  }
}

provider "foxy" {
  base_url      = "https://api.foxycart.com"
  client_id     = "client_1Q6iX3A1UjKNUZxEeV7P"
  # Either add client_secret here, or override by setting the environment variable FOXY_CLIENTSECRET
  refresh_token = "W2B9upaRrfwQtGP4VPLTR4QLlMdC6btP4Qy9UUGY"
}

resource "foxy_webhook" "example" {
  format        = "json"
  name          = "New webhook"
  url           = "https://example.com/new"
  event_resource = "transaction"
}

resource "foxy_cart_template" "default" {
  description      = "This is the default Cart template"
  content          = "<p>New improved cart template</p>"
}

resource "foxy_checkout_template" "default" {
  description      = "Checkout Template"
  content          = "<p>New checkout template</p>"
}

resource "foxy_email_template" "default" {
  description  = "Email ReceiptTemplate"
  subject      = "{{ store_name }} Order ({{ order_id }})"
  content_html = "<p>New email template</p>"
  content_text = "New email template"
}
