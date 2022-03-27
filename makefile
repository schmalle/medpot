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
	cp $(BIN) /usr/bin/

clean:
	rm $(BIN)
