package bdd

import "github.com/octohelm/x/testing/snapshot"

func MatchSnapshot(build func(s *snapshot.Snapshot), snapshotName string) Checker {
	return asChecker(snapshot.Match(snapshotName), Build(build))
}

func Snapshot(name string) *snapshot.Snapshot {
	return snapshot.Load(name)
}
