go install ./...
~/snap-v0.14.0-beta/bin/snapctl plugin load $GOPATH/bin/snap-plugin-collector-procstat
~/snap-v0.14.0-beta/bin/snapctl plugin load ~/snap-v0.14.0-beta/plugin/snap-processor-passthru
~/snap-v0.14.0-beta/bin/snapctl plugin load ~/snap-v0.14.0-beta/plugin/snap-publisher-file

~/snap-v0.14.0-beta/bin/snapctl task create -t $GOPATH/src/github.com/intelsdi-x/snap-plugin-collector-procstat/examples/tasks/procstat.json
