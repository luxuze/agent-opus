#!/bin/bash

# Initialize database script

set -e

echo "Initializing database..."

# Wait for MySQL to be ready
until mysql -h 127.0.0.1 -u root -p123456 -e "SELECT 1"; do
  echo "Waiting for MySQL to be ready..."
  sleep 2
done

# Create database if not exists
mysql -h 127.0.0.1 -u root -p123456 <<EOF
CREATE DATABASE IF NOT EXISTS agent_platform;
CREATE USER IF NOT EXISTS 'agent_user'@'%' IDENTIFIED BY 'agent_password';
GRANT ALL PRIVILEGES ON agent_platform.* TO 'agent_user'@'%';
FLUSH PRIVILEGES;
EOF

echo "Database initialized successfully!"
