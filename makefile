BIN := medpot

compile:
	bash scripts/compile_medpot.sh

compile_docker:
	bash scripts/compile_docker.sh

clean:
	rm $(BIN)
