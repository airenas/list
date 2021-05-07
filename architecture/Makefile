generate: zoom.multiple.png

zoom.multiple.png: zoom.multiple.puml 
	cat $< | docker run --rm -i think/plantuml -tpng > $@_
	mv $@_ $@
