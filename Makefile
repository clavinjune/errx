# Copyright 2022 clavinjune/errutil
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

include tools.mk

check:
	@go run $(licenser) verify
	@go run $(linter) run

clean:
	@rm -rf result.json coverage.out

fmt:
	@gofmt -w -s .
	@go run $(importer) -w .
	@go vet ./...
	@go mod tidy
	@go run $(licenser) apply -r "clavinjune/errutil" 2> /dev/null

test:
	@go test -count=1 -v -json -coverprofile=coverage.out -covermode=count ./... > result.json

test/coverage: test
	@go tool cover -html=coverage.out
