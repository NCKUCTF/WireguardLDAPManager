builder := go
builddir := bin
exe := wgldapmanager
instdir := /usr/local/bin
#config := .env.sample
#systemd := piproxyclient.service
#install := install.sh

all: $(builddir)/$(exe) #$(builddir)/$(config) $(install) $(builddir)/$(systemd)

#$(builddir)/$(config): $(config)
#		cp $(config) $(builddir)/$(config)

#$(builddir)/$(install): $(install)
#		cp $(install) $(builddir)/$(install)

#$(builddir)/$(systemd): $(systemd)
#		cp $(systemd) $(builddir)/$(systemd)

$(builddir)/$(exe): main.go go.mod go.sum models router utils
		$(builder) build -o $(builddir)/$(exe) $<

install: $(instdir)/$(exe)

$(instdir)/$(exe): $(builddir)/$(exe)
		cp $(builddir)/$(exe) $(instdir)/$(exe)
		chown root:root $(instdir)/$(exe)
		chmod 4755 $(instdir)/$(exe)

uninstall:
		rm -rf $(instdir)/$(exe)
clean: 
		rm -rf $(builddir)
