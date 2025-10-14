#!/bin/bash
# V-UI Installation - Simple & Clean

echo "Installing Docker..."
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
rm get-docker.sh

systemctl start docker
systemctl enable docker

echo "Installing Docker Compose..."
curl -L "https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

echo "Cloning V-UI..."
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2

echo "Setup..."
cp .env.example .env
sed -i "s|changeme_32_char_random_secret_key_here_minimum|$(openssl rand -base64 32 | tr -d '\n')|g" .env

echo "Building..."
/usr/local/bin/docker-compose build --no-cache

echo "Starting..."
/usr/local/bin/docker-compose up -d

echo "Complete! Access at http://$(hostname -I | awk '{print $1}')"
