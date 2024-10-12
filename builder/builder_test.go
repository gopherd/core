package builder

import "testing"

func TestInfo_String(t *testing.T) {
	runtimeVersion = func() string {
		return "go.test"
	}
	tests := []struct {
		name string
		info buildInfo
		want string
	}{
		{
			name: "empty",
			info: buildInfo{},
			want: " () built at  by go.test",
		},
		{
			name: "full",
			info: buildInfo{
				Name:     "app",
				Version:  "v1.0.0",
				Branch:   "main",
				Commit:   "abcdefg",
				DateTime: "2021-01-01T00:00:00Z",
			},
			want: "app v1.0.0(main: abcdefg) built at 2021-01-01T00:00:00Z by go.test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.info.String(); got != tt.want {
				t.Errorf("buildInfo.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	defaultAppName := appName()
	tests := []struct {
		name string
		want buildInfo
	}{
		{
			name: "empty",
			want: buildInfo{
				Name:     "",
				Version:  "",
				Branch:   "",
				Commit:   "",
				DateTime: "",
			},
		},
		{
			name: "full",
			want: buildInfo{
				Name:     "app",
				Version:  "v1.0.0",
				Branch:   "main",
				Commit:   "abcdefg",
				DateTime: "2021-01-01T00:00:00Z",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name = tt.want.Name
			version = tt.want.Version
			branch = tt.want.Branch
			commit = tt.want.Commit
			datetime = tt.want.DateTime
			if name == "" {
				tt.want.Name = defaultAppName
			}
			if got := Info(); got != tt.want {
				t.Errorf("Info() = %v, want %v", got, tt.want)
			}
		})
	}
}
