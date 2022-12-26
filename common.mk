.PHONY: mod fmt imports lint vet

mod:
	go mod vendor
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint

fmt:
	find . -type f -name '*.go' -not -path "./vendor/*" | xargs -n 1 gofmt -d -e | tee gofmt.txt
	test ! -s gofmt.txt
	rm -rf gofmt.txt

imports:
	find . -type f -name '*.go' -not -path "./vendor/*" | xargs -n 1 goimports -d -e | tee goimports.txt
	test ! -s goimports.txt
	rm -rf goimports.txt

lint:
	find . -type f -name '*.go' -not -path "./vendor/*" | xargs -n 1 golint | tee golint.txt
	test ! -s golint.txt
	rm -rf golint.txt
