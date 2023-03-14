Terraform provider for configuring the Foxy shopping cart (https://foxy.io/)

__This is experimental code. It may destroy any Foxy installation you run it against. DO NOT USE IT 
UNLESS YOU KNOW WHAT YOU ARE DOING AND ARE WILLING TO ACCEPT THE RISKS.__

## Setup

Because this is not yet published to the Terraform registry, the only way of using it is to 
build locally, and add the following to your `~/.terraformrc`:

```
provider_installation {
    dev_overrides {
        "67bricks.com/terraform/foxycart" = "[path to your go binaries output]"
    }
    direct {}
}
```

On MacOS, the path is typically `/Users/[your username]/go/bin`

You will also need to replace the `client_id`, `client_secret` and `refresh_token` values in the
`main.tf` file with appropriate values for your Foxy setup (or override them with environment 
variables). Setting FOXY_CLIENTSECRET appropriately will allow you to use the client details that 
are already present in the config.

## Current status

* Managing webhooks broadly seems to work
* Managing cart templates and checkout templates works - but while the Foxy API supports multiple
  cart and checkout templates, Foxy itself only uses one. So you need to:
  * Create a placeholder resource in your TF file for a cart_template (e.g. called `default`)
  * Find the ID of the existing cart_template via the Foxy API
  * `terraform import foxy_cart_template.default [the id]`

## MIT License

Copyright (c) 2023 Inigo Surguy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.