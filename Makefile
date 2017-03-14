all:
	cd alsa-bindings && make clean && make
	#c-for-go -out ./ -fancy alsa-bindings.yml
	mv alsa-bindings/libalsa.so ./
	cp libalsa.so alsa/

