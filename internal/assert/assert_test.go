// assert_test.go
package assert

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestNever(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestNeverHelper")
	cmd.Env = append(os.Environ(), "TEST_NEVER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Expected exit status 1, got err: %v", err)
}

func TestNeverHelper(t *testing.T) {
	if os.Getenv("TEST_NEVER") != "1" {
		return
	}
	Never("This should always panic")
}

func TestNotNil_NilValue(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestNotNilNilValueHelper")
	cmd.Env = append(os.Environ(), "TEST_NOTNIL_NIL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Expected exit status 1, got err: %v", err)
}

func TestNotNilNilValueHelper(t *testing.T) {
	if os.Getenv("TEST_NOTNIL_NIL") != "1" {
		return
	}
	var item interface{} = nil
	NotNil(item, "Item is nil")
}

func TestNotNil_NonNilValue(t *testing.T) {
	var item interface{} = "not nil"
	NotNil(item, "Item is not nil")
}

func TestAssert_TrueCondition(t *testing.T) {
	Assert(true, "Condition is true")
}

func TestAssert_FalseCondition(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestAssertFalseConditionHelper")
	cmd.Env = append(os.Environ(), "TEST_ASSERT_FALSE=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Expected exit status 1, got err: %v", err)
}

func TestAssertFalseConditionHelper(t *testing.T) {
	if os.Getenv("TEST_ASSERT_FALSE") != "1" {
		return
	}
	Assert(false, "Condition is false")
}

func TestNoError_NoError(t *testing.T) {
	var err error = nil
	NoError(err, "There should be no error")
}

func TestNoError_WithError(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestNoErrorWithErrorHelper")
	cmd.Env = append(os.Environ(), "TEST_NOERROR_ERROR=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Expected exit status 1, got err: %v", err)
}

func TestNoErrorWithErrorHelper(t *testing.T) {
	if os.Getenv("TEST_NOERROR_ERROR") != "1" {
		return
	}
	var err error = fmt.Errorf("An error occurred")
	NoError(err, "There is an error")
}
