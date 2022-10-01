package main

import "testing"

func TestUrlRequest_IsValid(t *testing.T) {
	type fields struct {
		TargetUrl string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "TestUrlRequest_IsValid - Http ok",
			fields: fields{TargetUrl: "http://www.google.com"},
			want:   true,
		},
		{
			name:   "TestUrlRequest_IsValid - Https OK",
			fields: fields{TargetUrl: "https://www.google.com"},
			want:   true,
		},
		{
			name:   "TestUrlRequest_IsValid - Subdirectory OK",
			fields: fields{TargetUrl: "https://kubernetes.io/es/docs/concepts/overview/what-is-kubernetes/"},
			want:   true,
		},
		{
			name:   "TestUrlRequest_IsValid - Query params OK",
			fields: fields{TargetUrl: "https://www.youtube.com/watch?v=14WXUslfe3o"},
			want:   true,
		},
		{
			name:   "TestUrlRequest_IsValid - Anchor OK",
			fields: fields{TargetUrl: "https://github.com/magefree/mage#features"},
			want:   true,
		},
		{
			name:   "TestUrlRequest_IsValid - Not valid missing scheme",
			fields: fields{TargetUrl: "www.google.com"},
			want:   false,
		},
		{
			name:   "TestUrlRequest_IsValid - Not valid missing scheme",
			fields: fields{TargetUrl: "google.com"},
			want:   false,
		},
		{
			name:   "TestUrlRequest_IsValid - Not valid empty",
			fields: fields{TargetUrl: ""},
			want:   false,
		},
		{
			name:   "TestUrlRequest_IsValid - Not valid random string",
			fields: fields{TargetUrl: "foo"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &UrlRequest{
				TargetUrl: tt.fields.TargetUrl,
			}
			if got := ur.IsValid(); got != tt.want {
				t.Errorf("UrlRequest.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
