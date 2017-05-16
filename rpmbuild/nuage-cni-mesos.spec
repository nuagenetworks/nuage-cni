%define nuage_cni_mesos_binary   nuage-cni-mesos
%define nuage_cni_service  nuage-cni.service
%define nuage_cni_service_file  scripts/mesos/nuage-cni.service
%define nuage_cni_daemon  nuage-cni
%define nuage_cni_yaml_file   nuage-cni.yaml
%define nuage_cni_netconf_file   cninetconf/mesos/nuage-net.conf
%define nuage_cni_netconf   nuage-net.conf
%undefine _missing_build_ids_terminate_build

Name: nuage-cni-mesos
Version: 0.0
Release: 1%{?dist}
Summary: Nuage CNI Plugin
Group: System Environments/Daemons
License: ALU EULA and ASL 2.0
Source0: nuage-cni-mesos-%{version}.tar.gz

BuildRequires:  %{?go_compiler:compiler(go_compiler)}%{!?go_compiler:golang}

%description
%{summary}

%prep
%setup -q

%build

%pre
if [ "$1" = "2" ]; then
    cp $RPM_BUILD_ROOT/etc/default/%{nuage_cni_yaml_file} $RPM_BUILD_ROOT/etc/default/%{nuage_cni_yaml_file}.orig
fi

%install
install --directory $RPM_BUILD_ROOT/usr/bin
install --directory $RPM_BUILD_ROOT/etc/default
install --directory $RPM_BUILD_ROOT/etc/systemd/system
install --directory $RPM_BUILD_ROOT/etc/cni/net.d

install -m 755 %{nuage_cni_mesos_binary} $RPM_BUILD_ROOT/usr/bin
install -m 755 %{nuage_cni_service_file} $RPM_BUILD_ROOT/etc/systemd/system/%{nuage_cni_service}
install -m 644 %{nuage_cni_yaml_file} $RPM_BUILD_ROOT/etc/default/%{nuage_cni_yaml_file}
install -m 644 %{nuage_cni_netconf_file} $RPM_BUILD_ROOT/etc/cni/net.d/%{nuage_cni_netconf}

%post
if [ "$1" = "2" ]; then
    mv $RPM_BUILD_ROOT/etc/default/%{nuage_cni_yaml_file}.orig $RPM_BUILD_ROOT/etc/default/%{nuage_cni_yaml_file}
fi
systemctl enable %{nuage_cni_daemon}
systemctl start %{nuage_cni_daemon}

%preun
if [ "$1" = "0" ]; then
    systemctl stop %{nuage_cni_daemon}
    systemctl disable %{nuage_cni_daemon}
fi

%clean
rm -rf $RPM_BUILD_ROOT

%files

/usr/bin/%{nuage_cni_mesos_binary}
/etc/systemd/system/%{nuage_cni_service}
/etc/default/%{nuage_cni_yaml_file}
/etc/cni/net.d/%{nuage_cni_netconf}
%attr(644, root, nobody) /etc/default/%{nuage_cni_yaml_file}
%doc

%changelog
