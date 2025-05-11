resource "proxmox_lxc" "ipxe" {
  hostname = "ipxe.hnatekmar.xyz"
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
    size    = "2G"
    storage = "storage"
  }
  network {
    name = "net0"
    bridge = "vmbr0"
    ip = "172.16.100.51/24"
    gw = "172.16.100.1"
  }
}

resource "null_resource" "install_ipxe" {
  depends_on = [proxmox_lxc.ipxe]
  triggers = {
    always_run = timestamp()
  }
  provisioner "local-exec" {
    command =<<-EOT
      echo "root@${proxmox_lxc.ipxe.network[0].ip}" | sed "s/\/.*$//g"> inventory.ini
      ansible-playbook -i inventory.ini ./ansible/ipxe.yaml
    EOT
  }
}