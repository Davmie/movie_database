.PHONY: test run genData clean

genData:
	python3.11 scripts/genData.py

run:
	go run cmd/main.go

test:
	go clean -testcache
	cd internal && go test $$(go list ./... | grep -v /mocks) -cover

clean:
	rm -rf allure-results