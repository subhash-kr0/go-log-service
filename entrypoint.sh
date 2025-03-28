#!/bin/sh

# Ensure logs directory exists with correct permissions
mkdir -p logs
chmod 777 logs

# Ensure logrotate status file is accessible
logrotate_status_file="/tmp/logrotate.status"
touch $logrotate_status_file
chmod 666 $logrotate_status_file

# Start log rotation in the background
crond &

# Wait for logrotate service to be ready
sleep 2

# Force log rotation to start immediately
logrotate -s $logrotate_status_file -f /etc/logrotate.d/app

# Start the application
exec /app/app
