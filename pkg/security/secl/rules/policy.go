// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package rules

// Policy represents a policy file which is composed of a list of rules and macros
type Policy struct {
	Name    string
	Source  string
	Version string
	Rules   []*RuleDefinition
	Macros  []*MacroDefinition
}
