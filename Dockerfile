FROM registry.access.redhat.com/ubi8/ubi-minimal:8.4

WORKDIR /

COPY direct-csi /direct-csi
COPY CREDITS /licenses/CREDITS
COPY LICENSE /licenses/LICENSE

COPY centos.repo /etc/yum.repos.d/centos.repo

RUN \
    curl https://www.centos.org/keys/RPM-GPG-KEY-CentOS-Official -o /etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-Official && \
    microdnf install xfsprogs --nodocs && \
    microdnf clean all

ENTRYPOINT ["/direct-csi"]
