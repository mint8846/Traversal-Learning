#!/bin/bash

set -e

NFS_EXPORT_PATH=${NFS_EXPORT_PATH}
# No IP restrictions applied by default since this is for simple testing
NFS_CLIENT_SUBNET=${NFS_CLIENT_SUBNET:-*}

# Dynamically generate exports file
cat > /etc/exports << EOF
${NFS_EXPORT_PATH} ${NFS_CLIENT_SUBNET}(rw,sync,fsid=0,no_root_squash,no_subtree_check,insecure)
EOF

# Set NFS directory permissions
mkdir -p ${NFS_EXPORT_PATH}
chmod 755 ${NFS_EXPORT_PATH}

# Start NFS service
mount -t nfsd nfsd /proc/fs/nfsd 2>/dev/null || echo "nfsd already mounted"
mount -t rpc_pipefs rpc_pipefs /var/lib/nfs/rpc_pipefs 2>/dev/null || echo "rpc_pipefs already mounted"

service rpcbind start
service nfs-kernel-server start

# Apply exports
exportfs -rav

echo "NFS Server started successfully (${NFS_EXPORT_PATH})"
sleep 2  # Wait NFS server

# Starting UDC process
exec /app/udc