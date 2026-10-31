package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Target struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Filter string `yaml:"filter"`
	Sort   string `yaml:"sort"`
	Order  string `yaml:"order"`
}

var (
	validTargetTypes = []string{
		"issues",
		"pull-requests",
	}

	validTargetSorts = []string{
		"comments",
		"reactions",
		"reactions-+1",
		"reactions--1",
		"reactions-smile",
		"reactions-thinking_face",
		"reactions-heart",
		"reactions-tada",
		"interactions",
		"created",
		"updated",
	}

	validTargetOrders = []string{
		"asc",
		"desc",
	}
)

func (t Target) Validate() error {
	// validate name
	if !isValidName(t.Name) {
		return fmt.Errorf("invalid target.name %q, must be lowercase alphanum", t.Name)
	}

	// validate type
	if !slices.Contains(validTargetTypes, t.Type) {
		return fmt.Errorf("invalid target.type %q, must be one of: %s", t.Type, strings.Join(validTargetTypes, ", "))
	}

	// validate filter
	if strings.TrimSpace(t.Filter) == "" {
		return errors.New("missing target.filter")
	}

	// validate sort
	if t.Sort != "" && !slices.Contains(validTargetSorts, t.Sort) {
		return fmt.Errorf("invalid target.sort %q, must be one of: %s", t.Sort, strings.Join(validTargetSorts, ", "))
	}

	// validate order
	if t.Order != "" && !slices.Contains(validTargetOrders, t.Order) {
		return fmt.Errorf("invalid target.order %q, must be one of: %s", t.Order, strings.Join(validTargetOrders, ", "))
	}

	return nil
}

func isValidName(name string) bool {
	for _, ch := range name {
		if (ch < 'a' || ch > 'z') && (ch < '0' || ch > '9') {
			return false
		}
	}

	return true
}
