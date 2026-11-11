package model

import "testing"

func TestPerfumeID(t *testing.T) {
	cases := map[[2]string]string{
		{"Eau Sauvage", "Dior"}:           "eau-sauvage-dior",
		{"J'adore", "Dior"}:               "jadore-dior",
		{"Ô de Lancôme", "Lancôme"}:       "o-de-lancome-lancome",
		{"Babycat", "Yves Saint Laurent"}: "babycat-yves-saint-laurent",
		{"Eau des Merveilles", "Hermès"}:  "eau-des-merveilles-hermes",
		{"Vanilla Extasy", "Montale"}:     "vanilla-extasy-montale",
	}
	for in, want := range cases {
		if got := PerfumeID(in[0], in[1]); got != want {
			t.Errorf("PerfumeID(%q, %q) = %q, want %q", in[0], in[1], got, want)
		}
		if !ValidID(want) {
			t.Errorf("ValidID(%q) = false", want)
		}
	}
	for _, bad := range []string{"", "../etc", "a b", "UPPER", "-x", "x--y"} {
		if ValidID(bad) {
			t.Errorf("ValidID(%q) = true", bad)
		}
	}
}

func TestCleanName(t *testing.T) {
	if got := CleanName("Vanilla Extasy :"); got != "Vanilla Extasy" {
		t.Fatalf("got %q", got)
	}
}
