echo "Installing systemd..."

#Create the service file
echo "[Unit]
  Description=kaching-go
  After=network-online.target

  [Service]
  User=ubuntu
  WorkingDirectory=/home/ubuntu/kaching-go
  ExecStart=/home/ubuntu/kaching-go/kaching-go
  StandardOutput=journal
  StandardError=journal
  Restart=always
  RestartSec=3
  LimitNOFILE=4096

  [Install]
  WantedBy=multi-user.target" > kaching-go.service

sudo mv kaching-go.service /etc/systemd/system/
sudo systemctl enable kaching-go.service
