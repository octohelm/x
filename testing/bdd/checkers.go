package bdd

import (
	"fmt"
	"iter"
	"slices"

	"github.com/octohelm/x/cmp"
	"github.com/octohelm/x/testing/snapshot"
	testingv2 "github.com/octohelm/x/testing/v2"
)

func SliceHaveLen[Slice ~[]E, E any](expect int, actual Slice) Checker {
	return testingv2.Helper(testingv2.Expect(actual,
		testingv2.Helper(testingv2.Be(func(slices Slice) error {
			if len(slices) != expect {
				return fmt.Errorf("expected %d slices, got %d", expect, len(slices))
			}
			return nil
		})),
	))
}

func ErrorAs[V error](expect *V, err error) Checker {
	return testingv2.Helper(testingv2.Expect(err,
		testingv2.Helper(testingv2.ErrorAs(expect)),
	))
}

func ErrorIs(expect error, err error) Checker {
	return testingv2.Helper(testingv2.Expect(err,
		testingv2.Helper(testingv2.ErrorIs(expect))),
	)
}

func HasError(err error) Checker {
	return testingv2.Helper(testingv2.Expect(err,
		testingv2.Helper(testingv2.Be(func(v error) error {
			if err == nil {
				return fmt.Errorf("expected error, got nil")
			}
			return nil
		})),
	))
}

func NoError(err error) Checker {
	return testingv2.Helper(testingv2.Expect(err,
		testingv2.Helper(testingv2.Be(func(v error) error {
			if err != nil {
				return fmt.Errorf("expected no error, got %v", err)
			}
			return nil
		})),
	))
}

func Zero[V any](actual V) Checker {
	return testingv2.Helper(testingv2.Expect(actual,
		testingv2.Helper(testingv2.Be(cmp.Zero[V]())),
	))
}

func Nil[V any](actual V) Checker {
	return testingv2.Helper(testingv2.Expect(actual,
		testingv2.Helper(testingv2.Be(cmp.Nil[V]())),
	))
}

func True(actual bool) Checker {
	return testingv2.Helper(testingv2.Expect(actual,
		testingv2.Helper(testingv2.Be(cmp.True())),
	))
}

func False(actual bool) Checker {
	return testingv2.Helper(testingv2.Expect(actual,
		testingv2.Helper(testingv2.Be(cmp.False())),
	))
}

func Equal[V any](expect V, actual V) Checker {
	return testingv2.Helper(testingv2.Expect(actual,
		testingv2.Helper(testingv2.Equal(expect)),
	))
}

func NotEqual[V any](expect V, actual V) Checker {
	return testingv2.Helper(testingv2.Expect(actual,
		testingv2.Helper(testingv2.NotEqual(expect)),
	))
}

func EqualSeq[V any](expect iter.Seq[V], actual iter.Seq[V]) Checker {
	return testingv2.Helper(testingv2.Expect(
		slices.AppendSeq(make([]V, 0), actual),
		testingv2.Helper(testingv2.Equal(slices.AppendSeq(make([]V, 0), expect))),
	))
}

func MatchSnapshot(build func(s *snapshot.Snapshot), snapshotName string) Checker {
	c := &snapshot.Context{}
	s, err := c.Load()
	if err != nil {
		panic(err)
	}

	build(s)

	return testingv2.Helper(testingv2.Expect(
		s.Files(),
		testingv2.Helper(
			testingv2.MatchSnapshot(snapshotName),
		),
	))
}
