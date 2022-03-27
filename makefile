BIN := medpot

compile:
	bash scripts/compile_medpot.sh

install:
	mkdir -p /etc/medpot/
	mkdir -p /var/log/medpot
	cp ./template/* /etc/medpot/
	touch /var/log/medpot/medpot.log
	bash scripts/dependencies.sh
	cp medpot /usr/bin/

clean:
	rm $(BIN)
