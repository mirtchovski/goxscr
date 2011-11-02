include $(GOROOT)/src/Make.inc

TARG=deco\
	moire\
	palette\
	popsquares\
	qc\
	rorschach\
	spirograph\
	wander\

SHARED=xscr

all: $(SHARED:%=%.$O) $(TARG)

$(TARG): %: %.$O $(SHARED:%=%.$O)
	$(LD) -o $@ $<

%.$O: %.go Makefile
	$(GC) -o $@ $<


clean:
	rm -f *.[$(OS)] $(TARG) $(CLEANFILES)
