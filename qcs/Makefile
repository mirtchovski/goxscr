include $(GOROOT)/src/Make.inc

TARG=\
	qcs\

all: $(SHARED:%=%.$O) $(TARG)

$(TARG): %: %.$O $(SHARED:%=%.$O)
	$(LD) -o $@ $<

%.$O: %.go Makefile
	$(GC) -o $@ $<


clean:
	rm -f *.[$(OS)] $(TARG) $(CLEANFILES)
