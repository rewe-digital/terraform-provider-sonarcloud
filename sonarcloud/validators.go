package sonarcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Copied from https://www.terraform.io/plugin/framework/validation
type stringLengthBetweenValidator struct {
	Min int
	Max int
}

func stringLengthBetween(min int, max int) *stringLengthBetweenValidator {
	return &stringLengthBetweenValidator{Min: min, Max: max}
}

func (v stringLengthBetweenValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("string length must be between %d and %d", v.Min, v.Max)
}

func (v stringLengthBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("string length must be between `%d` and `%d`", v.Min, v.Max)
}

// Validate checks if the length of the string attribute is between Min and Max
func (v stringLengthBetweenValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if str.Unknown || str.Null {
		return
	}

	strLen := len(str.Value)

	if strLen < v.Min || strLen > v.Max {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Length",
			fmt.Sprintf("String length must be between %d and %d, got: %d.", v.Min, v.Max, strLen),
		)

		return
	}
}

type allowedOptionsValidator struct {
	Options []string
}

func allowedOptions(options ...string) *allowedOptionsValidator {
	return &allowedOptionsValidator{Options: options}
}

func (v allowedOptionsValidator) Description(_ context.Context) string {
	return fmt.Sprintf("option must be one of %v", v.Options)
}

func (v allowedOptionsValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("option must be one of `%v`", v.Options)
}

func (v allowedOptionsValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var set types.Set
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &set)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if set.Unknown || set.Null {
		return
	}

	var values []string
	diags = set.ElementsAs(ctx, &values, true)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	options := make(map[string]struct{})
	for _, option := range v.Options {
		options[option] = struct{}{}
	}

	for _, val := range values {
		if _, ok := options[val]; !ok {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Invalid Set Element",
				fmt.Sprintf("Element must be one of %v, got: %s.", v.Options, val),
			)

			return
		}
	}
}

type allowedSetOptionsValidator struct {
	Options []string
}

func allowedSetOptions(options ...string) *allowedOptionsValidator {
	return &allowedOptionsValidator{Options: options}
}

func (v allowedSetOptionsValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value in set must be one of %v", v.Options)
}

func (v allowedSetOptionsValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value in set must be one of `%v`", v.Options)
}

func (v allowedSetOptionsValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var set types.Set
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &set)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if set.Unknown || set.Null {
		return
	}

	var lastElem string
	valid := false
	for _, elem := range set.Elems {
		lastElem = elem.String()
		valid = false
		for _, option := range v.Options {
			if option == elem.String() {
				valid = true
				break
			}
		}
		if !valid {
			break
		}
	}

	if !valid {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Value in Set",
			fmt.Sprintf("String must be one of %v, got: %s.", v.Options, lastElem),
		)

		return
	}
}
