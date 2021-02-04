#!/bin/bash

set -x

export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get --yes upgrade

apt-get install --yes make gcc

# Install docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get install --yes docker-ce docker-ce-cli containerd.io
usermod -aG docker vagrant

# Install k8s
snap install microk8s --classic --channel=1.19/stable
microk8s.status --wait-ready
ufw allow in on cbr0
ufw allow out on cbr0
ufw default allow routed

microk8s enable registry
microk8s enable storage
microk8s enable dns
microk8s enable ingress

microk8s.kubectl config view --raw > /tailing-sidecar/.kube-config

snap alias microk8s.kubectl kubectl

usermod -a -G microk8s vagrant

# Install go
wget "https://golang.org/dl/go1.15.6.linux-amd64.tar.gz"
tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz
rm go1.15.6.linux-amd64.tar.gz


# Install kubebuilder
curl -L "https://go.kubebuilder.io/dl/2.3.1/linux/amd64" | tar -xz -C /tmp/
mv /tmp/kubebuilder_2.3.1_linux_amd64 /usr/local/kubebuilder

echo "export PATH=$PATH:/usr/local/kubebuilder/bin:/usr/local/go/bin" >> /home/vagrant/.bashrc

# Install kustomize
curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash
mv kustomize /usr/local/bin/
