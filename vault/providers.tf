terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "4.5.0"
    }
    authentik = {
      source = "goauthentik/authentik"
      version = "2025.4.0"
    }
  }
}


provider "authentik" {
  # Configuration options
  url = "https://authentik.hnatekmar.xyz"
}

provider "vault" {
  address = "https://openbao.hnatekmar.xyz"
}

