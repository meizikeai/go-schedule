#!/bin/bash
set -euo pipefail

# === Configuration ===
ip_test="10.0.0.6"                   # Test environment IP
ip_release="10.0.0.8"    # Production environment IP
target_ips=""                         # Target server IP address
ssh_key=""                            # Path to the local SSH private key
remote_admin_user="ubuntu"            # SSH user (with key access)
remote_deploy_user="work"             # Service execution user (target owner)
ssh_key_base="${SSH_KEY_BASE_PATH:-}" # Path to the local SSH private key

old_name="go-schedule"
project_name="$old_name"   # Application name
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
  target_ips=$ip_test
  ssh_key="$ssh_key_base/dev_main"
  config_dir="$cmd_dir/test"
elif [ "$1" == "release" ]; then
  target_ips=$ip_release
  ssh_key="$ssh_key_base/srd_qcloud_ubuntu"
  config_dir="$cmd_dir/release"
else
  echo "✗ Error: Deployment target not specified. Please use 'test' or 'release'." 1>&2
  exit 1
fi
if [ ! -f "$ssh_key" ]; then
  echo "✗ Error: SSH Key not found at $ssh_key" 1>&2
  exit 1
fi

# --- Deployment start ---
echo "Deploying $project_name ..."
echo "Environment : $1 environment"
echo "=========================="

# 1. Build Artifact
echo "[1/4] Building Go binary (linux/amd64)..."
cd $cmd_dir
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $config_dir/main main.go
[ -f $config_dir/main ] || {
  echo "Build failed"
  exit 1
}
echo "|---- Build artifact 'main' generated successfully."

# 2. Create Deployment Package
echo "[2/4] Creating $tar_file..."
cd $cmd_dir
COPYFILE_DISABLE=1 tar --no-xattrs --exclude=".DS_Store" -zcvf $tar_file $1
rm -f $config_dir/main
echo "|---- Deployment archive created and temporary binary cleaned up."

echo "[3/4] Transferring artifact to remote server..."
for ip in $target_ips; do
  echo "|---- $remote_deploy_user @ $ip"

  # 3. Upload Artifact and Prepare Remote
  ssh -i $ssh_key $remote_admin_user@$ip "mkdir -p /tmp/$times"
  scp -i $ssh_key $tar_file $remote_admin_user@$ip:/tmp/$times/
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
    sudo systemctl restart $old_name
    echo "|---- Service restart triggered via Systemctl."
EOF
  echo "|---- Finished deployment on $ip"
done

# 5. Cleanup Temporary Files
rm -f $tar_file

# 6. Deployment Notification
if [ "$1" == "release" ]; then
  curl -X POST -H "Content-Type: application/json" -d '{
    "msg_type": "post",
    "content": {
      "post": {
        "zh_cn": {
          "title": "Deployment Notification",
          "content": [
            [{"tag":"text","text":"Project: '$project_name'"}],
            [{"tag":"text","text":"Environment: '$1'"}],
            [{"tag":"text","text":"Status: Success"}]
          ]
        }
      }
    }
  }' https://open.feishu.cn/open-apis/bot/v2/hook/00000000-0000-0000-0000-000000000000
  echo
fi

# --- Deployment Complete ---
echo "=========================="
echo "Deployment completed successfully."
echo "Path        : $remote_project"
echo "Version     : $times"
echo "Time        : $(date '+%Y-%m-%d %H:%M:%S')"
