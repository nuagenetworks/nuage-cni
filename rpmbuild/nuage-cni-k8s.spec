%define nuage_cni_k8s_binary   nuage-cni-k8s
%define nuage_cni_service  nuage-cni.service
%define nuage_cni_service_file  scripts/k8s/nuage-cni.service
%define nuage_cni_daemon  nuage-cni
%define nuage_cni_yaml_file   nuage-cni.yaml
%define nuage_cni_netconf_file   cninetconf/k8s/nuage-net.conf
%define nuage_cni_netconf   nuage-net.conf
%define vsp_k8s_datadir /usr/share/vsp-k8s
%define vsp_k8s_yaml vsp-k8s.yaml
%define vsp_k8s_yaml_path %{vsp_k8s_datadir}/%{vsp_k8s_yaml}
%define nuage_vrs_platform_script /usr/share/openvswitch/scripts/vrs-platform-lib 
%undefine _missing_build_ids_terminate_build

Name: nuage-cni-k8s
Version: 0.0
Release: 1%{?dist}
Summary: Nuage CNI Plugin
Group: System Environments/Daemons
License: ALU EULA and ASL 2.0
Source0: nuage-cni-k8s-%{version}.tar.gz

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
install --directory $RPM_BUILD_ROOT%{vsp_k8s_datadir}

install -m 755 %{nuage_cni_k8s_binary} $RPM_BUILD_ROOT/usr/bin
install -m 755 %{nuage_cni_service_file} $RPM_BUILD_ROOT/etc/systemd/system/%{nuage_cni_service}
install -m 644 %{nuage_cni_yaml_file} $RPM_BUILD_ROOT/etc/default/%{nuage_cni_yaml_file}
install -m 644 %{nuage_cni_netconf_file} $RPM_BUILD_ROOT/etc/cni/net.d/%{nuage_cni_netconf}
install -m 644 %{vsp_k8s_yaml}.template  $RPM_BUILD_ROOT%{vsp_k8s_yaml_path}

%post

if [ "$1" = "1" ]; then # first time install only.
test -e %{nuage_vrs_platform_script} || exit 0
. %{nuage_vrs_platform_script}
add_platform k8s 
fi

if [ "$1" = "2" ]; then
    mv $RPM_BUILD_ROOT/etc/default/%{nuage_cni_yaml_file}.orig $RPM_BUILD_ROOT/etc/default/%{nuage_cni_yaml_file}
    cp $RPM_BUILD_ROOT%{vsp_k8s_yaml_path}.orig $RPM_BUILD_ROOT%{vsp_k8s_yaml_path}
fi
systemctl enable %{nuage_cni_daemon}
systemctl start %{nuage_cni_daemon}

%preun
if [ "$1" = "0" ]; then
    systemctl stop %{nuage_cni_daemon}
    systemctl disable %{nuage_cni_daemon}
    test -e %{nuage_vrs_platform_script} || exit 0
    . %{nuage_vrs_platform_script}
    remove_platform k8s 
fi

%postun
if [ "$1" = "0" ]; then
   rm -rf $RPM_BUILD_ROOT%{vsp_k8s_datadir}
fi

%clean
rm -rf $RPM_BUILD_ROOT

%files

/usr/bin/%{nuage_cni_k8s_binary}
/etc/systemd/system/%{nuage_cni_service}
/etc/default/%{nuage_cni_yaml_file}
/etc/cni/net.d/%{nuage_cni_netconf}
%{vsp_k8s_datadir}
%attr(644, root, nobody) /etc/default/%{nuage_cni_yaml_file}
%attr(644, root, nobody) %{vsp_k8s_yaml_path}
%doc

%changelog
