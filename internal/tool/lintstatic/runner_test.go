package lintstatic

import "testing"

func TestNormalizeTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "empty args use all packages",
			args: nil,
			want: []string{allPackagesPattern},
		},
		{
			name: "keeps provided targets",
			args: []string{"./internal/..."},
			want: []string{"./internal/..."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := normalizeTargets(tt.args)
			if len(got) != len(tt.want) {
				t.Fatalf("normalizeTargets() length = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("normalizeTargets()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
