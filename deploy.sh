#!/bin/bash
set -euo pipefail

# === Configuration ===
ip=""                                # Target server IP address
ssh_key="/Users/name/runtime/secret" # Local path to SSH private key directory
ip_test="10.0.0.8"                   # Test environment IP
ip_release="10.0.8.6"                # Production environment IP
remote_admin_user="ubuntu"           # SSH user (with key access)
remote_deploy_user="work"            # Service execution user (target owner)

project_name="go-schedule"        # Application name
root_dir=$(pwd)                   # Local project root directory
cmd_dir="$root_dir/cmd"           # Go source directory
config_dir=""                     # Configuration directory path
tar_file="${project_name}.tar.gz" # Deployment archive filename

times=$(date +"%Y%m%d.%H%M%S")             # Unique timestamp for versioning
remote_pwd="/home/$remote_deploy_user"     # Remote base directory for deploy user
remote_tmp="$remote_pwd/$times"            # Remote temporary staging directory
remote_backup="$remote_pwd/backup"         # Remote backup directory
remote_project="$remote_pwd/$project_name" # Remote final production path

# --- Environment Setup ---
if [ "$1" == "test" ]; then
  ip=$ip_test
  ssh_key="$ssh_key/dev_main"
  config_dir="$cmd_dir/test"
elif [ "$1" == "release" ]; then
  ip=$ip_release
  ssh_key="$ssh_key/srd_qcloud_ubuntu"
  config_dir="$cmd_dir/release"
else
  echo "✗ Error: Deployment target not specified. Use 'test' or 'release'." 1>&2
  exit 1
fi

# --- Deployment start ---
echo "Deploying $project_name..."
echo "Environment : $1"
echo "Server      : $remote_deploy_user @ $ip"
echo "=========================="

# 1. Build Artifact
echo "[1/4] Building Go binary (linux/amd64)..."
cd "$cmd_dir"
GOOS=linux GOARCH=amd64 go build -o main main.go
echo "|---- Build artifact 'main' created successfully."

# 2. Create Deployment Package
echo "[2/4] Creating deployment archive ($tar_file)..."
cd "$root_dir"
# Package binary and config directory, excluding macOS metadata.
COPYFILE_DISABLE=1 tar --no-xattrs --exclude=".DS_Store" --exclude=".__MACOSX" -czf "$tar_file" -C "$cmd_dir" main "$(basename "$config_dir")"
rm -f "$cmd_dir/main" # Clean up local build artifact
echo "|---- Archive created and local binary cleaned up."

# 3. Upload Artifact and Prepare Remote
echo "[3/4] Transferring artifact to remote server..."

# Create temporary upload directory (/tmp/$times) and target directories owned by deploy user.
ssh -i "$ssh_key" "$remote_admin_user@$ip" "\
  mkdir -p /tmp/$times && \
  sudo -u $remote_deploy_user bash -c 'mkdir -p $remote_tmp $remote_project $remote_backup'"

# Securely copy the tarball to the remote temporary location.
scp -i "$ssh_key" "$tar_file" "$remote_admin_user@$ip:/tmp/$times/"
rm -f "$tar_file" # Clean up local tar file

# Move files from /tmp to $remote_tmp and ensure deploy user owns them.
ssh -i "$ssh_key" "$remote_admin_user@$ip" "\
  sudo mv /tmp/$times/* $remote_tmp && \
  sudo chown -R $remote_deploy_user:$remote_deploy_user $remote_tmp && \
  rm -rf /tmp/$times"
echo "|---- Transfer successful. Files staged and permission set."

# 4. Remote Deployment and Service Restart
echo "[4/4] Executing remote deployment script (as $remote_deploy_user)..."
# Use 'Here Document' for multi-line remote execution via SSH.
ssh -i "$ssh_key" -T "$remote_admin_user@$ip" "sudo -u $remote_deploy_user bash -s" <<EOF

# 4.1 Backup current production version
if [ -d "$remote_project" ]; then
  mv "$remote_project" "$remote_backup/$project_name.$times"
  mkdir -p "$remote_project"
  echo "|---- Current version backed up to $remote_backup/$project_name.$times"
fi

# 4.2 Extract artifact within staging area
tar -xzf "$remote_tmp/$tar_file" -C "$remote_tmp"
rm -f "$remote_tmp/$tar_file"
echo "|---- Deployment package extracted."

# 4.3 Atomic replacement of production directory
# Move content from staging to production, ignoring errors if source is empty.
mv -f "$remote_tmp"/* "$remote_project"/ 2>/dev/null || true
rm -rf "$remote_tmp" # Clean up staging directory
echo "|---- Production path updated ($remote_project)."

# 4.4 Restart application service
supervisorctl restart $project_name
echo "|---- Service restart initiated via supervisor."
EOF

# --- Deployment Complete ---
echo "=========================="
echo "Deployed successfully."
echo "Path        : $remote_project"
echo "Time        : $(date '+%Y-%m-%d %H:%M:%S')"
