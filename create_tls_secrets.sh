#!/bin/bash

set -e

# Output directory
OUTPUT_DIR="./tls-secrets"
mkdir -p "$OUTPUT_DIR"

# List of secrets and durations: name duration
SECRETS=(
  tls-secret-365-1 365
  tls-secret-365-2 365
  tls-secret-10-1 10
  tls-secret-10-2 10
  tls-secret-expired-1 -1
  tls-secret-expired-2 -1
)

# Check if faketime is available
if ! command -v faketime >/dev/null 2>&1; then
  echo "❌ Error: 'faketime' is not installed. Please install it and try again."
  exit 1
fi

# Process secrets
i=0
while [ $i -lt ${#SECRETS[@]} ]; do
  secret_name="${SECRETS[$i]}"
  days="${SECRETS[$((i+1))]}"
  echo "🔁 Re-creating TLS secret '$secret_name' (valid for $days days)"

  key_file="${OUTPUT_DIR}/${secret_name}.key"
  crt_file="${OUTPUT_DIR}/${secret_name}.crt"

  # Delete existing secret if it exists
  if kubectl get secret "$secret_name" >/dev/null 2>&1; then
    echo "🗑️  Deleting existing secret '$secret_name'"
    kubectl delete secret "$secret_name"
  fi

  # Generate private key
  openssl genrsa -out "$key_file" 2048

  # Use faketime if the cert should already be expired
  if [ "$days" -lt 0 ]; then
    echo "⏱️  Using faketime to backdate cert (expired)"
    faketime "2 days ago" openssl req -x509 -new -nodes \
      -key "$key_file" \
      -subj "/CN=${secret_name}" \
      -days 1 \
      -out "$crt_file"
  else
    openssl req -x509 -new -nodes \
      -key "$key_file" \
      -subj "/CN=${secret_name}" \
      -days "$days" \
      -out "$crt_file"
  fi

  # Create Kubernetes TLS secret
  kubectl create secret tls "$secret_name" \
    --cert="$crt_file" \
    --key="$key_file"

  i=$((i + 2))
done

echo "✅ All TLS secrets created successfully."
