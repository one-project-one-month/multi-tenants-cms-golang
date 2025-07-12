#!/bin/bash

# Kong-Consul integration script
echo "Setting up Kong with Consul service discovery..."

# Wait for all services to be ready
sleep 30

# Function to check if Kong is ready
check_kong() {
  curl -s http://kong:8001/status > /dev/null 2>&1
  return $?
}

# Function to check if Consul is ready
check_consul() {
  curl -s http://consul:8500/v1/status/leader > /dev/null 2>&1
  return $?
}

# Function to check if service is registered in Consul
check_consul_service() {
  local service_name=$1
  curl -s http://consul:8500/v1/agent/services | grep -q "$service_name"
  return $?
}

# Function to get service details from Consul
get_service_from_consul() {
  local service_name=$1
  curl -s http://consul:8500/v1/agent/services | \
    python3 -c "
import json, sys
try:
    data = json.load(sys.stdin)
    for service_id, service in data.items():
        if service['Service'] == '$service_name':
            print(f\"{service['Address']}:{service['Port']}\")
            break
except:
    pass
"
}

# Function to register Kong service
register_kong_service() {
  local service_name=$1
  local service_host=$2
  local service_port=$3
  local route_path=$4
  
  echo "Registering $service_name in Kong..."
  
  # Create upstream for load balancing
  curl -i -X POST http://kong:8001/upstreams \
    --data name=${service_name}-upstream \
    --data host_header=$service_host \
    --data algorithm=round-robin
  
  # Add target to upstream
  curl -i -X POST http://kong:8001/upstreams/${service_name}-upstream/targets \
    --data target=${service_host}:${service_port} \
    --data weight=100
  
  # Create service pointing to upstream
  curl -i -X POST http://kong:8001/services/ \
    --data name=${service_name} \
    --data host=${service_name}-upstream \
    --data protocol=http \
    --data path=/
  
  # Create route
  curl -i -X POST http://kong:8001/services/${service_name}/routes \
    --data paths[]=${route_path} \
    --data strip_path=false \
    --data preserve_host=false
}

# Wait for Kong to be ready
echo "Waiting for Kong to be ready..."
while ! check_kong; do
  echo "Waiting for Kong Admin API..."
  sleep 5
done

# Wait for Consul to be ready
echo "Waiting for Consul to be ready..."
while ! check_consul; do
  echo "Waiting for Consul..."
  sleep 5
done

echo "Kong and Consul are ready!"

# Wait a bit more for services to register with Consul
echo "Waiting for services to register with Consul..."
sleep 15

# Configure CMS Service
echo "Configuring CMS Service..."
if check_consul_service "cms-multi-tenant-api"; then
  echo "Found CMS service in Consul"
  service_address=$(get_service_from_consul "cms-multi-tenant-api")
  if [ -n "$service_address" ]; then
    host=$(echo $service_address | cut -d: -f1)
    port=$(echo $service_address | cut -d: -f2)
    register_kong_service "cms-service" "$host" "$port" "/api/v1/cms"
  else
    echo "Using default CMS service configuration"
    register_kong_service "cms-service" "cms-main-system" "8081" "/api/v1/cms"
  fi
else
  echo "CMS service not found in Consul, using default configuration"
  register_kong_service "cms-service" "cms-main-system" "8081" "/api/v1/cms"
fi

# Configure Email Service
echo "Configuring Email Service..."
if check_consul_service "email-service"; then
  echo "Found Email service in Consul"
  service_address=$(get_service_from_consul "email-service")
  if [ -n "$service_address" ]; then
    host=$(echo $service_address | cut -d: -f1)
    port=$(echo $service_address | cut -d: -f2)
    register_kong_service "email-service" "$host" "$port" "/api/v1/email"
  else
    echo "Using default Email service configuration"
    register_kong_service "email-service" "email-service" "8080" "/api/v1/email"
  fi
else
  echo "Email service not found in Consul, using default configuration"
  register_kong_service "email-service" "email-service" "8080" "/api/v1/email"
fi

# Add some useful Kong plugins
echo "Adding Kong plugins..."

# Add rate limiting to CMS service
curl -i -X POST http://kong:8001/services/cms-service/plugins \
  --data name=rate-limiting \
  --data config.minute=100 \
  --data config.hour=1000

# Add rate limiting to Email service
curl -i -X POST http://kong:8001/services/email-service/plugins \
  --data name=rate-limiting \
  --data config.minute=50 \
  --data config.hour=500

# Add CORS plugin globally
curl -i -X POST http://kong:8001/plugins \
  --data name=cors \
  --data config.origins=* \
  --data config.methods=GET,POST,PUT,DELETE,OPTIONS \
  --data config.headers=Accept,Accept-Version,Content-Length,Content-MD5,Content-Type,Date,X-Auth-Token,Authorization

# Add request/response logging
curl -i -X POST http://kong:8001/plugins \
  --data name=file-log \
  --data config.path=/tmp/file.log

echo "Kong configuration with Consul discovery complete!"

# Show current configuration
echo "=== Current Kong Configuration ==="
echo "Services:"
curl -s http://kong:8001/services | python3 -m json.tool 2>/dev/null || echo "Services configured"

echo -e "\nRoutes:"
curl -s http://kong:8001/routes | python3 -m json.tool 2>/dev/null || echo "Routes configured"

echo -e "\nUpstreams:"
curl -s http://kong:8001/upstreams | python3 -m json.tool 2>/dev/null || echo "Upstreams configured"

echo -e "\nPlugins:"
curl -s http://kong:8001/plugins | python3 -m json.tool 2>/dev/null || echo "Plugins configured"

echo -e "\n=== Consul Services ==="
curl -s http://consul:8500/v1/agent/services | python3 -m json.tool 2>/dev/null || echo "Consul services available"

echo -e "\n=== Testing Setup ==="
echo "Test URLs:"
echo "- Kong Admin: http://localhost:8001"
echo "- Kong Manager: http://localhost:1337"
echo "- Consul UI: http://localhost:8500"
echo "- CMS API: http://localhost:8000/api/v1/cms/health"
echo "- Email API: http://localhost:8000/api/v1/email/health"
echo "- Dozzle Logs: http://localhost:8083"

echo "Kong-Consul integration setup completed!"