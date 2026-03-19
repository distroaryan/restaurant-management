.PHONY: test-integration test-e2e test-integration-sum test-e2e-sum

test-integration:
	go test -v ./internal/tests/integration/...

test-e2e:
	go test -v ./internal/tests/e2e/...

testsum-integration:
	gotestsum --format testname -- -v ./internal/tests/integration/...

testsum-e2e:
	gotestsum --format testname -- -v ./internal/tests/e2e/...

load-test:
	k6 run scripts/load_test.js
