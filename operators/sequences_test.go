package operators

import (
	"testing"
)

var (
	b = Terminal(`b`, []byte("b"))
	c = Terminal(`c`, []byte("c"))
)

func TestConcat(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		rule := Concat(`abc`, a, b, c)
		nodes := rule([]byte("abc"))
		if len(nodes) != 1 {
			t.Error("expected one node")
			return
		}

		node := nodes[0]
		if string(node.Value) != "abc" {
			t.Error("invalid value")
		}

		if len(node.Children) != 3 {
			t.Error("expected three children")
		}

		if string(node.Children[0].Value) != "a" {
			t.Error("invalid value")
		}

		if string(node.Children[1].Value) != "b" {
			t.Error("invalid value")
		}

		if string(node.Children[2].Value) != "c" {
			t.Error("invalid value")
		}
	})

	t.Run("Complex", func(t *testing.T) {
		rule := Concat(``,
			a,
			Repeat0Inf(`*( a / b )`,
				Alts(`a / b`,
					a,
					b,
				),
			),
			a,
		)

		if nodes := rule([]byte("aa")); len(nodes) != 1 {
			// "aa"
			t.Errorf("expected one node, got %d", len(nodes))
		}

		if nodes := rule([]byte("aaa")); len(nodes) != 1 {
			// "aa" and "aaa"
			t.Errorf("expected two nodes, got %d", len(nodes))
		}

		if nodes := rule([]byte("aaba")); len(nodes) != 1 {
			// "aa" and "aaba"
			t.Errorf("expected two nodes, got %d", len(nodes))
		}
	})
}

func TestAlts(t *testing.T) {
	rule := Alts(`a / b`, a, b)
	for _, s := range []string{
		"a",
		"b",
		"abc",
	} {
		t.Run("", func(t *testing.T) {
			if len(rule([]byte(s))) != 1 {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if rule([]byte("c")) != nil {
		t.Errorf("value found for \"c\"")
	}

	t.Run("Complex", func(t *testing.T) {
		rule := Repeat0Inf(`*( a / b )`,
			Alts(`a / b`,
				a,
				b,
			),
		)

		if nodes := rule([]byte("aa")); len(nodes) != 3 {
			// "", "a" and "aa"
			t.Errorf("expected three node, got %d", len(nodes))
		}

		if nodes := rule([]byte("aaa")); len(nodes) != 4 {
			// "", "a", "aa" and "aaa"
			t.Errorf("expected four nodes, got %d", len(nodes))
		}

		if nodes := rule([]byte("aaba")); len(nodes) != 5 {
			// "", "a", "aa", "aab" and "aaba"
			t.Errorf("expected five nodes, got %d", len(nodes))
		}
	})

	t.Run("Complex 2", func(t *testing.T) {
		rule := Concat("(*a / *b) a",
			Alts("*a / *b",
				Repeat0Inf("*a", a),
				Repeat0Inf("*b", b),
			),
			a,
		)

		if nodes := rule([]byte("aa")); len(nodes) != 1 {
			t.Errorf("expected one node, got %d", len(nodes))
		} else {
			if len(nodes[0].Children) != 2 {
				t.Errorf("expected two subnodes, %d", len(nodes[0].Children))
			}
			for _, n := range nodes[0].Children {
				if len(n.Value) != 1 {
					t.Errorf("expected length 1, got %d", len(n.Value))
				}
			}
		}
	})
}
