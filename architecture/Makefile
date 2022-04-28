generate: zoom.multiple.png multiple.channels.png

%.png: %.puml 
	cat $< | docker run --rm -i think/plantuml -tpng > $@_
	mv $@_ $@
