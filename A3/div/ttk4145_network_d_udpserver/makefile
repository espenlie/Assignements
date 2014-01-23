# name of the executable
EXNAME = udp_d_server

# output file, comedi, warnings, debugSymbol, (-cov, -D)
DMDFLAGS = -of$(EXNAME) -w -g
DSOURCES = $(shell find -name "*.d")

all: $(EXNAME)

$(EXNAME) : $(DSOURCES)
	dmd $^ $(DMDFLAGS) && rm -f $(EXNAME)".o"

c:
	rm -f $(shell find -name "*.o") $(EXNAME)
	
clean: c

build: c all

run: build
	./$(EXNAME)

.PHONY: all r c
