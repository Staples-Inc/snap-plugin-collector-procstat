TEST_DIRS="github.com/intelsdi-x/snap-plugin-collector-procstat/procstat"
set -e
go test $TEST_DIRS -v -covermode=count
