#!/bin/bash
go build forge.go
mkdir ~/.forge
cp forge ~/.forge
echo "export PATH='$HOME/.forge:$PATH'" >> ~/.bashrc
echo "Forge Installed!"
echo "Install location: ~/.forge"
