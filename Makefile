run:
	@mkdir -p bin
	go run vyc.go main.vy bin/main.ll
	@echo
	clang bin/mir.ll -o bin/a.out -Wno-override-module
	@echo
	@bin/a.out

.PHONY: run