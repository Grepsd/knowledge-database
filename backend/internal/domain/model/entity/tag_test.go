package entity

import (
	uuid "github.com/satori/go.uuid"
	"reflect"
	"strings"
	"testing"
)

func TestCreateTag(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"name too long",
			args{strings.Repeat("i", NameMaxLength+1)},
			true,
		},
		{
			"valid",
			args{strings.Repeat("i", NameMaxLength)},
			false,
		},
		{
			"name too short",
			args{""},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTag(NewTagID(uuid.NewV4()), NewTagName(tt.args.name))
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewTag(t *testing.T) {
	type args struct {
		ID   uuid.UUID
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *Tag
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTag(NewTagID(tt.args.ID), NewTagName(tt.args.name))
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTag_Name(t *testing.T) {
	tagName := NewTagName("valid I")
	tag, err := NewTag(NewTagID(uuid.NewV4()), tagName)
	if err != nil {
		t.Error("could not create new tag")
	}
	t.Run("valid", func(t *testing.T) {
		if tag.Name() != tagName {
			t.Errorf("expected %s got %s", tagName.String(), tag.Name().String())
		}
	})
}

func TestTag_Strings(t *testing.T) {
	name := NewTagName("valid I")
	id := NewTagID(uuid.NewV4())
	t.Run("valid name.String", func(t *testing.T) {
		if name.String() != name.value {
			t.Errorf("expected %s got %s", name.String(), name.value)
		}
	})
	t.Run("valid id.String", func(t *testing.T) {
		if id.String() != id.value.String() {
			t.Errorf("expected %s got %s", id.String(), id.value.String())
		}
	})
}
