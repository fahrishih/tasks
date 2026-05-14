package domain

import (
	"errors"
	"testing"
)

func TestNewTask_ValidatesTitle(t *testing.T) {
	cases := []struct {
		name    string
		title   string
		wantErr error
	}{
		{"empty title", "", ErrInvalidTitle},
		{"whitespace-only title", "   ", ErrInvalidTitle},
		{"valid title", "Buy milk", nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewTask(tc.title, "")
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("got %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestTask_MarkComplete(t *testing.T) {
	task, _ := NewTask("Test", "")
	if task.Completed {
		t.Fatal("new task should not be completed")
	}
	task.MarkComplete()
	if !task.Completed {
		t.Fatal("MarkComplete should set Completed to true")
	}
}
