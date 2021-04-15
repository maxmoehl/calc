package calc

import (
	"testing"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    float64
		wantErr bool
	}{
		{
			name:    "simple addition",
			arg:     "2+3",
			want:    5,
			wantErr: false,
		},
		{
			name:    "simple subtraction",
			arg:     "2-3",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "simple multiplication",
			arg:     "2*3",
			want:    6,
			wantErr: false,
		},
		{
			name:    "simple division",
			arg:     "3/2",
			want:    1.5,
			wantErr: false,
		},
		{
			name:    "combining multiple operations",
			arg:     "2+3*2",
			want:    8,
			wantErr: false,
		},
		{
			name:    "simple parentheses calculation",
			arg:     "2*(3+4)",
			want:    14,
			wantErr: false,
		},
		{
			name:    "test square root macro",
			arg:     "1+sqrt{4}",
			want:    3,
			wantErr: false,
		},
		{
			name:    "test power macro",
			arg:     "pow{2,4}",
			want:    16,
			wantErr: false,
		},
		{
			name:    "test empty expression",
			arg:     "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "test complex expression",
			arg:     "sqrt{4}*2+(3*pow{1+1,sqrt{4}})",
			want:    16,
			wantErr: false,
		},
		{
			name:    "test parentheses only",
			arg:     "(3)",
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Eval(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("run() got = %v, want %v", got, tt.want)
			}
		})
	}
}
