#!/bin/bash

# Start containerd in the background
containerd > /var/log/containerd.log 2>&1 &
CONTAINERD_PID=$!

# Wait for containerd to be ready
timeout=30
while [ $timeout -gt 0 ]; do
  if crictl info > /dev/null 2>&1; then
    echo "containerd is ready."
    break
  fi
  echo "Waiting for containerd... ($timeout seconds remaining)"
  sleep 1
  ((timeout--))
done

if [ $timeout -eq 0 ]; then
  echo "Error: containerd failed to start within 30 seconds."
  echo "--- containerd logs ---"
  cat /var/log/containerd.log
  exit 1
fi

# Execute the test command
"$@"
EXIT_CODE=$?

# Cleanup
echo "Tests finished. Cleaning up containerd (PID $CONTAINERD_PID)..."
kill $CONTAINERD_PID
wait $CONTAINERD_PID 2>/dev/null

exit $EXIT_CODE
