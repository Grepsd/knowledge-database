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