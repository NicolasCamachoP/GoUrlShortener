package main

import (
	"fmt"
	"testing"
)

func Test_genId(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test_genId - Empty Ok",
			args: args{url: ""},
			want: "47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU=",
		},
		{
			name: "Test_genId - Ok",
			args: args{url: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
			want: "zTcvuFFIcA-ogJXjSS0_n1vrQ-VV5f8m2V9aatw2-OY=",
		},
		{
			name: "Test_genId - Ok Url",
			args: args{url: "https://google.com"},
			want: "BQRvJsg-jIiz3asuq2PQ0WIkrB5WRTX8dc3O7kegk40=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genId(tt.args.url); got != tt.want {
				t.Errorf("genId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortener_GetUrl(t *testing.T) {
	type fields struct {
		dbHandler IRepository
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "TestShortener_GetUrl - Internal error",
			fields:  fields{&MockFaultyDbHandler{}},
			wantErr: true,
		},
		{
			name:    "TestShortener_GetUrl - Not found",
			fields:  fields{&MockDbHandler{}},
			want:    "",
			wantErr: false,
		},
		{
			name:    "TestShortener_GetUrl - Ok",
			fields:  fields{&MockDbHandler{}},
			args:    args{id: "BQRvJsg-jIiz3asuq2PQ0WIkrB5WRTX8dc3O7kegk40="},
			want:    "https://google.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				dbHandler: tt.fields.dbHandler,
			}
			got, err := s.GetUrl(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Shortener.GetUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Shortener.GetUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortener_SaveUrl(t *testing.T) {
	type fields struct {
		dbHandler IRepository
	}
	type args struct {
		targetUrl string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "TestShortener_SaveUrl - Ok new",
			fields:  fields{dbHandler: &MockDbHandler{}},
			args:    args{"https://youtube.com"},
			want:    "FGeTGg6McbgYdz8bP7nVAxuYB1PIHJwTXITBjzhAfV0=",
			wantErr: false,
		},
		{
			name:    "TestShortener_SaveUrl - Ok already saved",
			fields:  fields{dbHandler: &MockDbHandler{}},
			args:    args{"https://google.com"},
			want:    "BQRvJsg-jIiz3asuq2PQ0WIkrB5WRTX8dc3O7kegk40=",
			wantErr: false,
		},
		{
			name:    "TestShortener_SaveUrl - Internal error",
			fields:  fields{dbHandler: &MockFaultyDbHandler{}},
			args:    args{"https://google.com"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				dbHandler: tt.fields.dbHandler,
			}
			got, err := s.SaveUrl(tt.args.targetUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Shortener.SaveUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Shortener.SaveUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortener_ShutDown(t *testing.T) {
	type fields struct {
		dbHandler IRepository
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "TestShortener_ShutDown - Empty handler",
			wantErr: true,
		},
		{
			name:    "TestShortener_ShutDown - Ok",
			fields:  fields{&MockDbHandler{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				dbHandler: tt.fields.dbHandler,
			}
			if err := s.ShutDown(); (err != nil) != tt.wantErr {
				t.Errorf("Shortener.ShutDown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//--------------------- Mocks ---------------------//
type MockDbHandler struct {
	IRepository
}

func (mdh *MockDbHandler) ShutDown() error { return nil }

func (mdh *MockDbHandler) GetItem(id string) (interface{}, error) {
	if id == "BQRvJsg-jIiz3asuq2PQ0WIkrB5WRTX8dc3O7kegk40=" {
		return "https://google.com", nil
	}
	return nil, nil
}

func (mdh *MockDbHandler) Exists(id string) bool {
	return id == "BQRvJsg-jIiz3asuq2PQ0WIkrB5WRTX8dc3O7kegk40="
}

func (mdh *MockDbHandler) SetItem(id string, value interface{}) error { return nil }

type MockFaultyDbHandler struct {
	IRepository
}

func (f *MockFaultyDbHandler) GetItem(id string) (interface{}, error) {
	return nil, fmt.Errorf("mock error")
}

func (f *MockFaultyDbHandler) SetItem(id string, value interface{}) error {
	return fmt.Errorf("mock error")
}
func (f *MockFaultyDbHandler) Exists(id string) bool {
	return false
}
