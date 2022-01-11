package views

import (
	"fmt"
	"github.com/hashicorp/terraform/internal/command/format"
	//"github.com/hashicorp/terraform/internal/command/jsonplan"
	//"github.com/hashicorp/terraform/internal/command/jsonstate"
	"github.com/hashicorp/terraform/internal/states/statefile"
	"github.com/hashicorp/terraform/internal/tfdiags"

	"github.com/hashicorp/terraform/internal/command/arguments"
	"github.com/hashicorp/terraform/internal/plans"
	"github.com/hashicorp/terraform/internal/terraform"
)

// FIXME: this is a temporary partial definition of the view for the show
// command, in place to allow access to the plan renderer which is now in the
// views package.
type Show interface {
	Plan(plan *plans.Plan, schemas *terraform.Schemas)

	State(stateFile *statefile.File, schemas *terraform.Schemas)

	// Diagnostics renders early diagnostics, resulting from argument parsing.
	Diagnostics(diags tfdiags.Diagnostics)
}

// FIXME: the show view should support both human and JSON types. This code is
// currently only used to render the plan in human-readable UI, so does not yet
// support JSON.
func NewShow(vt arguments.ViewType, view *View) Show {
	switch vt {
	case arguments.ViewJSON:
		return &ShowJSON{view: NewJSONView(view)}
	case arguments.ViewHuman:
		return &ShowHuman{view: view}
	default:
		panic(fmt.Sprintf("unknown view type %v", vt))
	}
}

type ShowHuman struct {
	view *View
}

var _ Show = (*ShowHuman)(nil)

func (v *ShowHuman) Plan(plan *plans.Plan, schemas *terraform.Schemas) {
	renderPlan(plan, schemas, v.view)
}

func (v *ShowHuman) State(stateFile *statefile.File, schemas *terraform.Schemas) {
	if stateFile == nil {
		v.view.streams.Println("No state.")
	}

	v.view.streams.Println(format.State(&format.StateOpts{
		State:   stateFile.State,
		Color:   v.view.colorize,
		Schemas: schemas,
	}))
}

func (v *ShowHuman) Diagnostics(diags tfdiags.Diagnostics) {
	v.view.Diagnostics(diags)
}

// The ShowJSON implementation renders show results as a JSON object.
// This object includes top-level fields summarizing the result, and an array
// of JSON diagnostic objects.
type ShowJSON struct {
	view *JSONView
}

var _ Show = (*ShowJSON)(nil)

func (v *ShowJSON) Plan(plan *plans.Plan, schemas *terraform.Schemas) {
	// JSON logic for plans from command
	//jsonPlan, err := jsonplan.Marshal(config, plan, stateFile, schemas)
	//
	//if err != nil {
	//	c.Ui.Error(fmt.Sprintf("Failed to marshal plan to json: %s", err))
	//	return 1
	//}
	//c.Ui.Output(string(jsonPlan))
	//return 0
}

func (v *ShowJSON) State(stateFile *statefile.File, schemas *terraform.Schemas) {
	// JSON logic for state files from command
	// At this point, it is possible that there is neither state nor a plan.
	// That's ok, we'll just return an empty object.
	//jsonState, err := jsonstate.Marshal(stateFile, schemas)
	//if err != nil {
	//	c.Ui.Error(fmt.Sprintf("Failed to marshal state to json: %s", err))
	//	return 1
	//}
	//c.Ui.Output(string(jsonState))
}

// Diagnostics should only be called if show cannot be executed.
// In this case, we choose to render human-readable diagnostic output,
// primarily for backwards compatibility.
func (v *ShowJSON) Diagnostics(diags tfdiags.Diagnostics) {
	v.view.Diagnostics(diags)
}
