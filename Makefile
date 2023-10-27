builder := go
builddir := bin
exe := $(builddir)/wgmanager
#config := .env.sample
#systemd := piproxyclient.service
#install := install.sh

all: $(exe) #$(builddir)/$(config) $(install) $(builddir)/$(systemd)

#$(builddir)/$(config): $(config)
#		cp $(config) $(builddir)/$(config)

#$(builddir)/$(install): $(install)
#		cp $(install) $(builddir)/$(install)

#$(builddir)/$(systemd): $(systemd)
#		cp $(systemd) $(builddir)/$(systemd)

$(exe): main.go go.mod go.sum models router utils
		$(builder) build -o $(exe) $<

.PHONY = clean

clean: 
		rm -r $(builddir)
