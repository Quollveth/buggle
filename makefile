APPNAME := buggle
BINDIR := ./bin
SRCDIR := ./cmd

GOFLAGS :=

all: daemon tui

daemon:
#	@mkdir -p $(BINDIR)
#	go build $(GOFLAGS) -o $(BINDIR)/daemon/$(APPNAME) $(SRCDIR)/daemon
	echo "daemon doesnt exist yet lmao"

tui:
	@mkdir -p $(BINDIR)
	go build $(GOFLAGS) -o $(BINDIR)/tui/$(APPNAME) $(SRCDIR)/tui

clean:
	rm -rf $(BINDIR)

.PHONY: all clean
