package enum_test

import (
	"testing"

	"github.com/gopherd/core/enum"
)

func TestRegistry(t *testing.T) {
	var r enum.Registry
	if r.Lookup("Color") != nil {
		t.Errorf("LookupDescriptor failed: Color found")
	}
	if err := r.Register(&enum.Descriptor{
		Name:        "Color",
		Description: "Color enum",
		Members: []enum.MemberDescriptor{
			{Name: "Red", Value: 0, Description: "Red color"},
			{Name: "Green", Value: 1, Description: "Green color"},
			{Name: "Blue", Value: 2, Description: "Blue color"},
		},
	}); err != nil {
		t.Errorf("RegisterDescriptor failed: %v", err)
	}
	if err := r.Register(&enum.Descriptor{
		Name:        "Shape",
		Description: "Shape enum",
		Members: []enum.MemberDescriptor{
			{Name: "Circle", Value: 0, Description: "Circle shape"},
			{Name: "Square", Value: 1, Description: "Square shape"},
			{Name: "Triangle", Value: 2, Description: "Triangle shape"},
		},
	}); err != nil {
		t.Errorf("RegisterDescriptor failed: %v", err)
	}
	if err := r.Register(&enum.Descriptor{
		Name:        "Color",
		Description: "Color enum",
		Members: []enum.MemberDescriptor{
			{Name: "Red", Value: 0, Description: "Red color"},
			{Name: "Green", Value: 1, Description: "Green color"},
			{Name: "Blue", Value: 2, Description: "Blue color"},
		},
	}); err == nil {
		t.Errorf("RegisterDescriptor failed: expected error, got nil")
	}

	if d := r.Lookup("Color"); d == nil {
		t.Errorf("LookupDescriptor failed: Color not found")
	} else {
		if d.Name != "Color" {
			t.Errorf("LookupDescriptor failed: expected Color, got %s", d.Name)
		}
		if d.Description != "Color enum" {
			t.Errorf("LookupDescriptor failed: expected Color enum, got %s", d.Description)
		}
		if len(d.Members) != 3 {
			t.Errorf("LookupDescriptor failed: expected 3 members, got %d", len(d.Members))
		}
	}

	if d := r.Lookup("Shape"); d == nil {
		t.Errorf("LookupDescriptor failed: Shape not found")
	} else {
		if d.Name != "Shape" {
			t.Errorf("LookupDescriptor failed: expected Shape, got %s", d.Name)
		}
		if d.Description != "Shape enum" {
			t.Errorf("LookupDescriptor failed: expected Shape enum, got %s", d.Description)
		}
		if len(d.Members) != 3 {
			t.Errorf("LookupDescriptor failed: expected 3 members, got %d", len(d.Members))
		}
	}

	if d := r.Lookup("Size"); d != nil {
		t.Errorf("LookupDescriptor failed: Size found")
	}
}
