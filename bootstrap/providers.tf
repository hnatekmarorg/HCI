terraform {
  required_providers {
    proxmox = {
      source = "Telmate/proxmox"
      version = "3.0.1-rc8"
    }
  }
}

provider "proxmox" {
  pm_api_url = var.proxmox_url
  pm_password = var.proxmox_api_password
  pm_user = var.proxmox_api_user
  pm_tls_insecure = true
}