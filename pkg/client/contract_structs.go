package client

type ViewStateResult struct {
	Values []StateItem   `json:"values"`
	Proof  TrieProofPath `json:"proof"`
}

type StateItem struct {
	Key   string        `json:"key"`
	Value string        `json:"value"`
	Proof TrieProofPath `json:"proof"`
}

// TrieProofPath is a set of serialized TrieNodes that are encoded in base64. Represent proof of inclusion of some TrieNode in the MerkleTrie.
type TrieProofPath = []string
