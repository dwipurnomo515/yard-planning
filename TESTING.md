🧪 Testing Guide
Running Tests
Run All Tests
bash
make test
Or manually:

bash
go test -v -race -cover ./...
Run Tests with Coverage Report
bash
make test-coverage
This will:

Run all tests with race detection
Generate coverage report
Create HTML coverage report (coverage.html)
Run Specific Tests
bash

# Test specific package

go test -v ./internal/repository

# Test specific function

go test -v -run TestYardRepository_GetByCode ./internal/repository

# Test with verbose output

go test -v -race ./...
Linting
Run Linter
bash
make lint
Or manually:

bash
golangci-lint run --config .golangci.yml
Auto-fix Issues
bash
make lint-fix
Install Linter
bash
make install-lint
Or manually:

bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
Test Coverage Goals
Repository Layer: > 80%
Service Layer: > 70%
Handler Layer: > 60%
Writing New Tests
Repository Test Example
go
func TestYourRepository_YourMethod(t \*testing.T) {
db, mock, err := sqlmock.New()
assert.NoError(t, err)
defer db.Close()

    repo := NewYourRepository(db)

    // Setup mock expectations
    rows := sqlmock.NewRows([]string{"id", "name"}).
        AddRow(1, "Test")

    mock.ExpectQuery("SELECT").
        WithArgs(1).
        WillReturnRows(rows)

    // Execute test
    result, err := repo.YourMethod(1)

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "Test", result.Name)

    assert.NoError(t, mock.ExpectationsWereMet())

}
Service Test Example
go
func TestYourService_YourMethod(t \*testing.T) {
// Setup
service := &YourService{}

    // Test cases
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test",
            want:    "expected",
            wantErr: false,
        },
        {
            name:    "invalid input",
            input:   "",
            want:    "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := service.YourMethod(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }

}
Continuous Integration
Add to .github/workflows/test.yml:

yaml
name: Tests

on:
push:
branches: [ main ]
pull_request:
branches: [ main ]

jobs:
test:
runs-on: ubuntu-latest
steps: - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run tests
        run: go test -v -race -cover ./...

      - name: Run linter
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
