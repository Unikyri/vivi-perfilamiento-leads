package config

import (
	"net/netip"
	"testing"
)

func TestCargarUsaDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x")
	t.Setenv("TRUSTED_PROXY_CIDRS", "")
	c, err := Cargar()
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if c.Puerto != "8080" {
		t.Errorf("Puerto = %q, esperado 8080", c.Puerto)
	}
	if c.TasaEA != 0.107 {
		t.Errorf("TasaEA = %v, esperado 0.107", c.TasaEA)
	}
	if c.LLMProvider != "gemini" {
		t.Errorf("LLMProvider = %q, esperado gemini", c.LLMProvider)
	}
}

func TestCargarFallaSinDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	if _, err := Cargar(); err == nil {
		t.Fatal("se esperaba error por DATABASE_URL vacía")
	}
}

func TestCargarDemoSeedRequiresExplicitTrue(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x")
	t.Setenv("DEMO_SEED", "false")
	c, err := Cargar()
	if err != nil || c.DemoSeed {
		t.Fatalf("config=%+v err=%v", c, err)
	}
	t.Setenv("DEMO_SEED", "true")
	c, err = Cargar()
	if err != nil || !c.DemoSeed {
		t.Fatalf("enabled config=%+v err=%v", c, err)
	}
}

func TestCargarTrustedProxyCIDRs(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x")
	t.Setenv("TRUSTED_PROXY_CIDRS", "127.0.0.0/8, ::1/128")
	c, err := Cargar()
	if err != nil || len(c.TrustedProxyCIDRs) != 2 {
		t.Fatalf("trusted proxies=%v err=%v", c.TrustedProxyCIDRs, err)
	}
	if !c.TrustedProxyCIDRs[0].Contains(netip.MustParseAddr("127.0.0.1")) || !c.TrustedProxyCIDRs[1].Contains(netip.MustParseAddr("::1")) {
		t.Fatalf("unexpected trusted proxies: %v", c.TrustedProxyCIDRs)
	}
}

func TestCargarRejectsMalformedTrustedProxyCIDRs(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x")
	t.Setenv("TRUSTED_PROXY_CIDRS", "127.0.0.1/999")
	if _, err := Cargar(); err == nil {
		t.Fatal("expected malformed trusted proxy CIDR error")
	}
}
