package config

import "testing"

func TestCargarUsaDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x")
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
