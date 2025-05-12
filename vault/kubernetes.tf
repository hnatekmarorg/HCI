variable "talos_version" {
  default = "v1.7.5"
}

data "talos_image_factory_extensions_versions" "gpu" {
  talos_version = var.talos_version
  filters = {
    names = [
      "nonfree-kmod-nvidia-production",
      "qemu-guest-agent",
      "nvidia-container-toolkit-production",
      "amd-ucode"
    ]
  }
}


resource "talos_image_factory_schematic" "gpu" {
  schematic = yamlencode(
    {
      customization = {
        systemExtensions = {
          officialExtensions = data.talos_image_factory_extensions_versions.gpu.extensions_info.*.name
        }
      }
    }
  )
}

resource "proxmox_vm_qemu" "first-gpu-node" {
  target_node = "donnager"
  pxe = true
  cpu = "host"
  boot                      = "order=net0"
  cores = 2
  network {
    id = 0
    model = "virtio"
    firewall = false
    bridge = "vmbr0"
  }
  agent = 1
  bios = "ovmf"
  efidisk {
    efitype = "4m"
    storage = "local-lvm"
  }
}

locals {
  node_config = jsonencode( {
    nodes = [
      {
        mac = proxmox_vm_qemu.first-gpu-node.network[0].macaddr
        type = "Talos"
        talos = {
          factoryHash = talos_image_factory_schematic.gpu.id
          version = var.talos_version
        }
      }
    ]
  })
}

resource "talos_machine_bootstrap" "gpu-bootstrap" {
  client_configuration = ""
  node                 = proxmox_vm_qemu.first-gpu-node.default_ipv4_address
}