include $(GOROOT)/src/Make.inc

TARG=deco\
	moire\
	palette\
	popsquares\
	rorschach\
	wander\

SHARED=xscr

all: $(SHARED:%=%.$O) $(TARG)

$(TARG): %: %.$O
	$(LD) -o $@ $<

%.$O: %.go Makefile
	$(GC) -o $@ $<


clean:
	rm -f *.[$(OS)] $(TARG) $(CLEANFILES)
