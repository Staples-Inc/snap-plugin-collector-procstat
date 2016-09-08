TEST_DIRS="github.com/Staples-Inc/snap-plugin-collector-procstat/procstat"
set -e
go test $TEST_DIRS -v -covermode=count
