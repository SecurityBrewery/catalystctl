package cmd

import (
	"testing"
)

func TestValidateCmd_Run(t *testing.T) {
	type fields struct {
		SourceFile string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"valid", fields{SourceFile: "v0.10.0.zip"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := testFile(t, tt.fields.SourceFile)
			if err != nil {
				t.Fatal(err)
			}

			c := &ValidateCmd{SourceFile: f}
			if err := c.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

