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
  description  = "Email Receipt Template"
  subject      = "{{ store_name }} Order ({{ order_id }})"
  content_html = "<p>New email template</p>"
  content_text = "New email template"
}

resource "foxy_store_info" "default" {
  store_name = "Terraform Test"
  store_domain = "terraformtest"
  store_url = "http://www.example.com/"
  store_email = "test@example.com"
  locale_code = "en_US"
  region = ""
  language = "german"
  postal_code = "99999"
  country = "GB"

  require_signed_shipping_rates = true
  products_require_expires_property = false
  is_maintenance_mode = false
  is_active = false
  hide_decimal_characters = false
  hide_currency_symbol = false
  features_multiship = false
  customer_password_hash_type = "phpass"
  customer_password_hash_config = "8"
  checkout_type = "default_account"
  bcc_on_receipt_email = true
  app_session_time = 604800

  shipping_address_type = "residential"
  timezone = "America/Los_Angeles"
  webhook_key = "pyt3lDjdEx8Nl1HlL5AavPcxEBbVm8Ptn6hZGTZRdAHijStkUTNc5IawT1De"

  use_webhook = false
  use_single_sign_on = false
  use_email_dns = false
  use_cart_validation = false
  use_international_currency_symbol = false
  use_remote_domain = false
}