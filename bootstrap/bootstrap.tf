// Basic infra for bootstrap
resource "proxmox_lxc" "bao-0" {
  hostname = "openbao.hnatekmar.xyz"
  ostemplate = var.template
  target_node = "donnager"
  memory = 1024
  unprivileged = false
  start = true
  features {
    nesting = true
  }
  password = var.proxmox_api_password
  ssh_public_keys = <<-EOT
    ${var.public_ssh_key}
  EOT
  rootfs {
    size    = "8G"
    storage = "storage"
  }
  network {
    name = "net0"
    bridge = "vmbr0"
    ip = "172.16.100.50/24"
    gw = "172.16.100.1"
  }
}
resource "terraform_data" "install_bao" {
  depends_on = [proxmox_lxc.bao-0]

  provisioner "local-exec" {
    command =<<-EOT
      echo [bao] > bao_inventory.ini
      echo "root@${proxmox_lxc.bao-0.network[0].ip}" | sed "s/\/.*$//g">> bao_inventory.ini
      ansible-playbook -i bao_inventory.ini ./ansible/bao.yaml
    EOT
  }
}