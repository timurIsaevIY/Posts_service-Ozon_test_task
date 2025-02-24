ifeq ($(OS),Windows_NT)
    SLEEP = powershell -noprofile -command "Start-Sleep -s"
else
    SLEEP = sleep
endif

.PHONY: test-env-up unit-test-cover wait-db clean

test-env-up:
	docker-compose -f docker-compose.test.yml up -d
	@$(MAKE) wait-db

wait-db:
	@echo "Waiting for database to be ready..."
ifeq ($(OS),Windows_NT)
	powershell -noprofile -command "Start-Sleep -s 5"
else
	sleep 5
endif

# Запуск только юнит-тестов с покрытием
unit-test-cover:
	go list ./... | findstr /V /C:"/tests" /C:"/mocks" /C:"/graphql"> unit_test_packages.txt
	echo mode: set > coverage.out
	for /F "tokens=*" %%i in (unit_test_packages.txt) do ( \
		go test -coverprofile=coverage.tmp -covermode=set -v -short %%i && \
		powershell -Command "Get-Content coverage.tmp | Select-Object -Skip 1 | Set-Content coverage_clean.tmp" && \
		type coverage_clean.tmp >> coverage.out \
	) || exit /b 1
	del coverage.tmp coverage_clean.tmp
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	start coverage.html

# Очистка временных файлов покрытия
clean:
	del coverage.out coverage.html unit_test_packages.txt 2>nul

