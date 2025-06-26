#!/bin/bash

service_exists() {
 systemctl list-unit-files --type=service --all | grep -Fq "bender.service"
}

service_is_active() {
 systemctl is-active --quiet "bender.service"
}

if service_exists; then
 if service_is_active; then
  sudo systemctl stop "bender.service"
 fi
fi

go build .

sudo mkdir -p /usr/local/bin

sudo cp ./bender /usr/local/bin/bender

sudo chmod +x /usr/local/bin/bender

service="[Unit]
Description=Bender service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/bender
Restart=on-failure
SupplementaryGroups=docker

[Install]
WantedBy=multi-user.target"

echo "$service" | sudo tee /etc/systemd/system/bender.service > /dev/null

sudo chown root:docker /var/run/docker.sock
sudo chmod 0660 /var/run/docker.sock

sudo systemctl daemon-reload

sudo systemctl enable bender

sudo systemctl start bender
sudo systemctl restart bender
