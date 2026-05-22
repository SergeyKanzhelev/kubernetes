FROM ubuntu:24.04

# Install basic dependencies
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    iptables \
    ca-certificates \
    iproute2 \
    socat \
    util-linux \
    kmod \
    && rm -rf /var/lib/apt/lists/*

# Install etcd
ENV ETCD_VERSION="v3.5.15"
RUN wget https://github.com/etcd-io/etcd/releases/download/${ETCD_VERSION}/etcd-${ETCD_VERSION}-linux-amd64.tar.gz \
    && tar xzvf etcd-${ETCD_VERSION}-linux-amd64.tar.gz \
    && mv etcd-${ETCD_VERSION}-linux-amd64/etcd* /usr/local/bin/ \
    && rm -rf etcd-${ETCD_VERSION}-linux-amd64*

# Install containerd
ENV CONTAINERD_VERSION="2.0.0"
RUN wget https://github.com/containerd/containerd/releases/download/v${CONTAINERD_VERSION}/containerd-${CONTAINERD_VERSION}-linux-amd64.tar.gz \
    && tar Cxzvf /usr/local containerd-${CONTAINERD_VERSION}-linux-amd64.tar.gz \
    && rm containerd-${CONTAINERD_VERSION}-linux-amd64.tar.gz

# Install runc
ENV RUNC_VERSION="v1.2.0"
RUN wget https://github.com/opencontainers/runc/releases/download/${RUNC_VERSION}/runc.amd64 \
    && install -m 755 runc.amd64 /usr/local/sbin/runc \
    && rm runc.amd64

# Install crictl
ENV CRICTL_VERSION="v1.31.1"
RUN wget https://github.com/kubernetes-sigs/cri-tools/releases/download/${CRICTL_VERSION}/crictl-${CRICTL_VERSION}-linux-amd64.tar.gz \
    && tar Cxzvf /usr/local/bin crictl-${CRICTL_VERSION}-linux-amd64.tar.gz \
    && rm crictl-${CRICTL_VERSION}-linux-amd64.tar.gz

# Configure containerd
RUN mkdir -p /etc/containerd \
    && containerd config default > /etc/containerd/config.toml \
    && sed -i 's/disabled_plugins = \[ "cri" \]/disabled_plugins = \[ \]/g' /etc/containerd/config.toml \
    && sed -i 's/SystemdCgroup = false/SystemdCgroup = true/g' /etc/containerd/config.toml

# Install CNI plugins
ENV CNI_VERSION="v1.6.0"
RUN mkdir -p /opt/cni/bin \
    && wget https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-linux-amd64-${CNI_VERSION}.tgz \
    && tar Cxzvf /opt/cni/bin cni-plugins-linux-amd64-${CNI_VERSION}.tgz \
    && rm cni-plugins-linux-amd64-${CNI_VERSION}.tgz

# Configure CNI
RUN mkdir -p /etc/cni/net.d
COPY <<EOF /etc/cni/net.d/10-containerd-net.conflist
{
 "cniVersion": "1.0.0",
 "name": "containerd-net",
 "plugins": [
   {
     "type": "bridge",
     "bridge": "cni0",
     "isGateway": true,
     "ipMasq": true,
     "promiscMode": true,
     "ipam": {
       "type": "host-local",
       "ranges": [
         [{
           "subnet": "10.88.0.0/16"
         }],
         [{
           "subnet": "2001:db8:4860::/64"
         }]
       ],
       "routes": [
         { "dst": "0.0.0.0/0" },
         { "dst": "::/0" }
       ]
     }
   },
   {
     "type": "portmap",
     "capabilities": {"portMappings": true},
     "externalSetMarkChain": "KUBE-MARK-MASQ"
   }
 ]
}
EOF

COPY <<EOF /etc/cni/net.d/99-loopback.conf
{
  "cniVersion": "1.0.0",
  "name": "lo",
  "type": "loopback"
}
EOF

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
