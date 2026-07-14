package main

import "testing"

func TestVerifyProof(t *testing.T) {
	tests := []struct {
		name     string
		tx       string
		root     string
		siblings []string
		want     bool
	}{
		{
			name:     "valid proof for a",
			tx:       "a",
			root:     "7ab33be6039bf251a90a1f1667f77860dcd939fc87e7cf61cfcbb49c812370d4",
			siblings: []string{"3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d", "3474e6da9272c8c6282c1fcd66ae507af47c831286b0ac7bd718f7d37738a2f7"},
			want:     true,
		},
		{
			name:     "valid proof for b",
			tx:       "b",
			root:     "7ab33be6039bf251a90a1f1667f77860dcd939fc87e7cf61cfcbb49c812370d4",
			siblings: []string{"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb", "3474e6da9272c8c6282c1fcd66ae507af47c831286b0ac7bd718f7d37738a2f7"},
			want:     true,
		},
		{
			name:     "valid proof for c",
			tx:       "c",
			root:     "7ab33be6039bf251a90a1f1667f77860dcd939fc87e7cf61cfcbb49c812370d4",
			siblings: []string{"18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4", "ab19ec537f09499b26f0f62eed7aefad46ab9f498e06a7328ce8e8ef90da6d86"},
			want:     true,
		},
		{
			name:     "invalid proof for wrong sibling",
			tx:       "a",
			root:     "7ab33be6039bf251a90a1f1667f77860dcd939fc87e7cf61cfcbb49c812370d4",
			siblings: []string{"2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6", "3474e6da9272c8c6282c1fcd66ae507af47c831286b0ac7bd718f7d37738a2f7"},
			want:     false,
		},
		{
			name:     "invalid proof for wrong tx",
			tx:       "x",
			root:     "7ab33be6039bf251a90a1f1667f77860dcd939fc87e7cf61cfcbb49c812370d4",
			siblings: []string{"3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d", "3474e6da9272c8c6282c1fcd66ae507af47c831286b0ac7bd718f7d37738a2f7"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := verifyProof(tt.tx, tt.root, tt.siblings); got != tt.want {
				t.Fatalf("verifyProof(%q, %q, %v) = %v, want %v", tt.tx, tt.root, tt.siblings, got, tt.want)
			}
		})
	}
}
