# disable debug package creation because of a bug when producing debuginfo
# packages: http://fedoraproject.org/wiki/PackagingDrafts/Go#Debuginfo
%global debug_package   %{nil}

Name:           sshproxy
Version:        0.1.0
Release:        1%{?dist}
Summary:        SSH proxy
License:        CeCILL-B
Source:         %{name}-%{version}.tar.xz
BuildArch:      %{ix86} x86_64 %{arm}

BuildRequires:  golang >= 1.3
BuildRequires:  golang(github.com/BurntSushi/toml)
BuildRequires:  golang(github.com/docker/docker/pkg/term)
BuildRequires:  golang(github.com/kr/pty)
BuildRequires:  golang(github.com/op/go-logging)
BuildRequires:  asciidoc
Summary:        SSH proxy

%description
%{summary}

This package provides an SSH proxy which can be used on a gateway to
automatically connect a remote user to a defined internal host.

%prep
%setup -q

%build
# set up temporary build gopath, and put our directory there
mkdir -p ./_build/src
ln -s $(pwd) ./_build/src/sshproxy

export GOPATH=$(pwd)/_build:%{gopath}
make

%install
make install DESTDIR=%{buildroot} prefix=%{_prefix} mandir=%{_mandir}
install -d %{buildroot}%{_sysconfdir}
install -p -m 0644 config/sshproxy.cfg %{buildroot}%{_sysconfdir}

%files
%doc Licence_CeCILL-B_V1-en.txt Licence_CeCILL-B_V1-fr.txt
%config(noreplace) %{_sysconfdir}/sshproxy.cfg
%{_sbindir}/sshproxy
%{_bindir}/sshproxy-replay
%{_mandir}/man5/sshproxy.cfg.5*
%{_mandir}/man8/sshproxy.8*
%{_mandir}/man8/sshproxy-replay.8*

%changelog
* Thu Feb 12 2015 Arnaud Guignard <arnaud.guignard@cea.fr> - 0.1.0-1
- sshproxy 0.1.0
