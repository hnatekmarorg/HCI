#!/bin/bash

export VAULT_ADDR=https://openbao.hnatekmar.xyz
bao operator init -key-shares=1 -key-threshold=1