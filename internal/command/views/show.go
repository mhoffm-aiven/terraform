package views

import (
	"fmt"
	"github.com/hashicorp/terraform/internal/command/format"
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
	//case arguments.ViewJSON:
	//	return &ShowJSON{view: view}
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

//// The ShowJSON implementation renders show results as a JSON object.
//// This object includes top-level fields summarizing the result, and an array
//// of JSON diagnostic objects.
//type ShowJSON struct {
//	view *View
//}
//
//var _ Show = (*ShowJSON)(nil)
//
//func (v *ShowJSON) Plan(plan *plans.Plan, schemas *terraform.Schemas) {
//	// FormatVersion represents the version of the json format and will be
//	// incremented for any change to this format that requires changes to a
//	// consuming parser.
//	const FormatVersion = "1.0"
//
//	type Output struct {
//		FormatVersion string `json:"format_version"`
//
//		// We include some summary information that is actually redundant
//		// with the detailed diagnostics, but avoids the need for callers
//		// to re-implement our logic for deciding these.
//		Valid        bool                    `json:"valid"`
//		ErrorCount   int                     `json:"error_count"`
//		WarningCount int                     `json:"warning_count"`
//		Diagnostics  []*viewsjson.Diagnostic `json:"diagnostics"`
//	}
//
//	output := Output{
//		FormatVersion: FormatVersion,
//		Valid:         true, // until proven otherwise
//	}
//	configSources := v.view.configSources()
//	for _, diag := range diags {
//		output.Diagnostics = append(output.Diagnostics, viewsjson.NewDiagnostic(diag, configSources))
//
//		switch diag.Severity() {
//		case tfdiags.Error:
//			output.ErrorCount++
//			output.Valid = false
//		case tfdiags.Warning:
//			output.WarningCount++
//		}
//	}
//	if output.Diagnostics == nil {
//		// Make sure this always appears as an array in our output, since
//		// this is easier to consume for dynamically-typed languages.
//		output.Diagnostics = []*viewsjson.Diagnostic{}
//	}
//
//	j, err := json.MarshalIndent(&output, "", "  ")
//	if err != nil {
//		// Should never happen because we fully-control the input here
//		panic(err)
//	}
//	v.view.streams.Println(string(j))
//
//	if diags.HasErrors() {
//	}
//}
//
//func (v *ShowJSON) State(stateFile *statefile.File, schemas *terraform.Schemas) {
//	// FormatVersion represents the version of the json format and will be
//	// incremented for any change to this format that requires changes to a
//	// consuming parser.
//	const FormatVersion = "1.0"
//
//	type Output struct {
//		FormatVersion string `json:"format_version"`
//
//		// We include some summary information that is actually redundant
//		// with the detailed diagnostics, but avoids the need for callers
//		// to re-implement our logic for deciding these.
//		Valid        bool                    `json:"valid"`
//		ErrorCount   int                     `json:"error_count"`
//		WarningCount int                     `json:"warning_count"`
//		Diagnostics  []*viewsjson.Diagnostic `json:"diagnostics"`
//	}
//
//	output := Output{
//		FormatVersion: FormatVersion,
//		Valid:         true, // until proven otherwise
//	}
//	configSources := v.view.configSources()
//	for _, diag := range diags {
//		output.Diagnostics = append(output.Diagnostics, viewsjson.NewDiagnostic(diag, configSources))
//
//		switch diag.Severity() {
//		case tfdiags.Error:
//			output.ErrorCount++
//			output.Valid = false
//		case tfdiags.Warning:
//			output.WarningCount++
//		}
//	}
//	if output.Diagnostics == nil {
//		// Make sure this always appears as an array in our output, since
//		// this is easier to consume for dynamically-typed languages.
//		output.Diagnostics = []*viewsjson.Diagnostic{}
//	}
//
//	j, err := json.MarshalIndent(&output, "", "  ")
//	if err != nil {
//		// Should never happen because we fully-control the input here
//		panic(err)
//	}
//	v.view.streams.Println(string(j))
//
//	if diags.HasErrors() {
//		return 1
//	}
//	return 0
//}
//
//// Diagnostics should only be called if show cannot be executed.
//// In this case, we choose to render human-readable diagnostic output,
//// primarily for backwards compatibility.
//func (v *ShowJSON) Diagnostics(diags tfdiags.Diagnostics) {
//	v.view.Diagnostics(diags)
//}
