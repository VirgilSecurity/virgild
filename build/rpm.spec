Name:  virgild
Version: %{version}
Release: 1%{?dist}
Summary: An Open Source public key management system
Group:	Applications/System
License: GPL
URL:	https://github.com/VirgilSecurity/virgild	
Source: %{name}

Requires(pre): /usr/sbin/groupadd /usr/sbin/useradd
%define debug_package %{nil}

%description
VirgilD is a public key management system which is fully compatible with Virgil cloud. VirgilD can work like a caching service of public keys. It is the core of distributed key management system. It operates in a chainable manner allowing to build decentralized trust models at any scale.

%install
pwd
ls -l
mkdir -p %{buildroot}/%{_bindir}
mkdir -p %{buildroot}/etc/systemd/system/multi-user.target.wants
mkdir -p %{buildroot}/var/lib/virgild
mkdir -p %{buildroot}/var/log
touch %{buildroot}/var/log/virgild.log

install -m 0755 %{_sourcedir}/virgild %{buildroot}/%{_bindir}/
install -m 0644 %{_sourcedir}/build/virgild.service %{buildroot}/etc/systemd/system/multi-user.target.wants/virgild.service
install -m 0644 %{_sourcedir}/build/virgild.conf %{buildroot}/etc/virgild.conf

%pre
if ! id -g %{name} > /dev/null 2>&1; then
  groupadd -r %{name}
fi
if ! id -u %{name} > /dev/null 2>&1; then
  useradd -g %{name} -G %{name}  -d %{_datadir}/%{name} -r -s /sbin/nologin %{name}
fi

%post
systemctl daemon-reload

%files
%defattr(-,root,root)
%attr(755,%{name},%{name}) %{_bindir}/virgild
%attr(644,%{name},%{name}) /etc/systemd/system/multi-user.target.wants/virgild.service
%attr(644,%{name},%{name}) /etc/virgild.conf
%attr(700,%{name},%{name}) /var/lib/virgild
%attr(644,%{name},%{name}) /var/log/virgild.log


