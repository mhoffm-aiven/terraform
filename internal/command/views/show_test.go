package views

import (
	"testing"

	"github.com/hashicorp/terraform/internal/command/arguments"
	"github.com/hashicorp/terraform/internal/terminal"
)

// Ensure that the correct view type and in-automation settings propagate to the
// Operation view.
func TestShowHuman_operation(t *testing.T) {
	streams, done := terminal.StreamsForTesting(t)
	defer done(t)
	v := NewPlan(arguments.ViewHuman, NewView(streams).SetRunningInAutomation(true)).Operation()
	if hv, ok := v.(*OperationHuman); !ok {
		t.Fatalf("unexpected return type %t", v)
	} else if hv.inAutomation != true {
		t.Fatalf("unexpected inAutomation value on Operation view")
	}
}

// Verify that Hooks includes a UI hook
func TestShowHuman_hooks(t *testing.T) {
	streams, done := terminal.StreamsForTesting(t)
	defer done(t)
	v := NewPlan(arguments.ViewHuman, NewView(streams).SetRunningInAutomation((true)))
	hooks := v.Hooks()

	var uiHook *UiHook
	for _, hook := range hooks {
		if ch, ok := hook.(*UiHook); ok {
			uiHook = ch
		}
	}
	if uiHook == nil {
		t.Fatalf("expected Hooks to include a UiHook: %#v", hooks)
	}
}
