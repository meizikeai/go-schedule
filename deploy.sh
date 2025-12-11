#!/bin/bash
set -euo pipefail

# === Configuration ===
ip=""                                    # Target server IP address
ssh_key="/Users/qinxikun/runtime/secret" # Path to the local SSH private key
ip_test="10.0.0.8"                       # Test environment IP
ip_release="10.0.8.6"                    # Production environment IP
remote_admin_user="ubuntu"               # SSH user (with key access)
remote_deploy_user="work"                # Service execution user (target owner)

project_name="go-schedule"        # Application name
root_dir=$(pwd)                   # Local project root directory
cmd_dir="$root_dir/cmd"           # Go source directory
config_dir=""                     # Configuration directory path
tar_file="${project_name}.tar.gz" # Deployment archive filename

times=$(date +"%Y%m%d.%H%M%S")             # Unique timestamp for versioning
remote_pwd="/home/$remote_deploy_user"     # Remote base directory for the deployment user
remote_project="$remote_pwd/$project_name" # Remote final production path
remote_version="$remote_project/version"   # Remote backup directory

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
  echo "✗ Error: Deployment target not specified. Please use 'test' or 'release'." 1>&2
  exit 1
fi

# --- Deployment start ---
echo "Deploying $project_name..."
echo "Environment : $1 environment"
echo "Server      : $remote_deploy_user @ $ip"
echo "=========================="

# 1. Build Artifact
echo "[1/4] Building Go binary (linux/amd64)..."
cd $cmd_dir
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $config_dir/main main.go
echo "|---- Build artifact 'main' generated successfully."

# 2. Create Deployment Package
echo "[2/4] Creating $tar_file..."
cd $cmd_dir
COPYFILE_DISABLE=1 tar --no-xattrs --exclude=".DS_Store" -zcvf $tar_file $1
rm -f $config_dir/main
echo "|---- Deployment archive created and temporary binary cleaned up."

# 3. Upload Artifact and Prepare Remote
echo "[3/4] Transferring artifact to remote server..."
ssh -i $ssh_key $remote_admin_user@$ip "mkdir -p /tmp/$times"
scp -i $ssh_key $tar_file $remote_admin_user@$ip:/tmp/$times/
rm -f $tar_file
ssh -i $ssh_key $remote_admin_user@$ip "\
  sudo mv /tmp/$times/$tar_file $remote_pwd/$tar_file && \
  sudo chown -R $remote_deploy_user:$remote_deploy_user $remote_pwd/$tar_file && \
  rm -rf /tmp/$times"
echo "|---- Transfer completed. Files staged and permissions updated."

# 4. Remote Deployment and Service Restart
echo "[4/4] Executing remote deployment script (as $remote_deploy_user)..."
ssh -i $ssh_key -T $remote_admin_user@$ip sudo -u $remote_deploy_user bash -s <<EOF
mkdir -p $remote_project $remote_version/$times
tar -xzf $remote_pwd/$tar_file -C $remote_version/$times
rm -f $remote_pwd/$tar_file
cd $remote_version
mv $times/$1/main $times/
ls -dt */ | tail -n +4 | xargs rm -rf
cd $remote_project
ln -s -f version/$times/* .
echo "|---- Deployment package extracted successfully."
supervisorctl restart $project_name
echo "|---- Service restart triggered via Supervisor."
EOF

# --- Deployment Complete ---
echo "=========================="
echo "Deployment completed successfully."
echo "Path        : $remote_project"
echo "Version     : $times"
echo "Time        : $(date '+%Y-%m-%d %H:%M:%S')"
