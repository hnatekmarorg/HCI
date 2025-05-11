variable "proxmox_url" {
  default = "https://localhost:8086"
  sensitive = true
}

variable "proxmox_api_user" {
  default = "root"
}

variable "proxmox_api_password" {
  sensitive = true
}

variable "public_ssh_key" {
  default = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDFPly2MYEeh4FtFtftOa0qasGW4VNIzYv/ZzheQ/dFs martin@fedora"
}

variable "template" {
  default = "local:vztmpl/debian-12-standard_12.7-1_amd64.tar.zst"
}