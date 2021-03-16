FROM fedora:33

RUN dnf update -y && dnf install --setopt=install_weak_deps=False -y  \
	libguestfs-tools-c  \
	&& dnf clean all \
	&& rm -rf /var/cache/dnf/
