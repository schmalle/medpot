BIN := medpot

compile: dependencies
	bash scripts/compile_medpot.sh

dependencies:
	bash scripts/dependencies.sh

install:
	mkdir -p /etc/medpot/
	mkdir -p /var/log/medpot
	cp ./template/* /etc/medpot/
	touch /var/log/medpot/medpot.log
	cp ./scripts/packet.txt /etc/medpot/
	cp $(BIN) /usr/bin/

uninstall:
	rm /usr/bin/$(BIN)
	rm -r /etc/medpot

clean:
	rm $(BIN)
