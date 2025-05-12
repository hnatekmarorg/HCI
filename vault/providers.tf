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
    proxmox = {
      source = "Telmate/proxmox"
      version = "3.0.1-rc8"
    }
    null = {
      source  = "hashicorp/null"
      version = "3.2.3"
    }
    talos = {
      source = "siderolabs/talos"
      version = "0.8.0"
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

provider "proxmox" {
  pm_api_url = var.proxmox_url
  pm_password = var.proxmox_api_password
  pm_user = var.proxmox_api_user
  pm_tls_insecure = true
}