package cmd

import (
	"os"
	"testing"
)

func TestUpgradeCmd_Run(t *testing.T) {
	type fields struct {
		AssumeVersion string
		SourceFile    string
		TargetFile    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"valid", fields{SourceFile: "v0.10.0.zip", TargetFile: "upgraded.zip"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := testFile(t, tt.fields.SourceFile)
			if err != nil {
				t.Fatal(err)
			}

			dir, err := os.MkdirTemp("", "catalystctl")
			if err != nil {
				t.Fatal(err)
			}
			err = os.Chdir(dir)
			if err != nil {
				t.Fatal(err)
			}

			c := &UpgradeCmd{
				AssumeVersion: tt.fields.AssumeVersion,
				SourceFile:    f,
				TargetFile:    tt.fields.TargetFile,
			}
			if err := c.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
