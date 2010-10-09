include $(GOROOT)/src/Make.inc

TARG=deco\
	moire\
	popsquares\
	rorschach\

SHARED=xscr

all: $(SHARED:%=%.$O) $(TARG)

$(TARG): %: %.$O
	$(LD) -o $@ $<

%.$O: %.go Makefile
	$(GC) -o $@ $<


clean:
	rm -f *.[$(OS)] $(TARG) $(CLEANFILES)
