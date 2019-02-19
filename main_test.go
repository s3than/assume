package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-ini/ini"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_returnProfileName(t *testing.T) {
	type args struct {
		sect *ini.Section
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := returnProfileName(tt.args.sect); got != tt.want {
				t.Errorf("returnProfileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_remainingTime(t *testing.T) {
	type args struct {
		sect *ini.Section
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := remainingTime(tt.args.sect); got != tt.want {
				t.Errorf("remainingTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_expired(t *testing.T) {
	type args struct {
		sect *ini.Section
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expired(tt.args.sect); got != tt.want {
				t.Errorf("expired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fmtDuration(t *testing.T) {
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name  string
		args  args
		want  time.Duration
		want1 time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := fmtDuration(tt.args.d)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fmtDuration() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("fmtDuration() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
