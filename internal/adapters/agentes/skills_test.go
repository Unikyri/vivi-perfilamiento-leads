package agentes

import (
	"strings"
	"testing"
)

var issue24SkillNames = []string{
	"tono-colsubsidio",
	"normalizacion-de-declarados",
	"siguiente-mejor-pregunta",
	"explicacion-financiera-humana",
	"dominio-caja",
	"redaccion-con-dignidad",
	"presentacion-de-proyectos",
	"ficha-comercial",
}

func TestTodasLasSkillsExisten(t *testing.T) {
	t.Parallel()
	for agente, nombres := range mapaSkills {
		for _, n := range nombres {
			if _, err := skillsFS.ReadFile(n + "/SKILL.md"); err != nil {
				t.Errorf("falta la skill %q del agente %q: %v", n, agente, err)
			}
		}
	}
}

func TestSkillsTienenFrontmatter(t *testing.T) {
	t.Parallel()
	for _, nombre := range issue24SkillNames {
		datos, err := skillsFS.ReadFile(nombre + "/SKILL.md")
		if err != nil {
			t.Fatalf("leyendo %s: %v", nombre, err)
		}
		contenido := strings.ReplaceAll(string(datos), "\r\n", "\n")
		if !strings.HasPrefix(contenido, "---\n") {
			t.Errorf("%s no empieza con frontmatter YAML", nombre)
		}
		frontmatterEnd := strings.Index(contenido[4:], "\n---\n")
		if frontmatterEnd < 0 {
			t.Errorf("%s no cierra el frontmatter YAML", nombre)
			continue
		}
		frontmatter := contenido[:frontmatterEnd+4]
		for _, field := range []string{"name:", "description:", "agente:", "fuente_de_verdad:"} {
			if !strings.Contains(frontmatter, "\n"+field) {
				t.Errorf("%s no declara %s", nombre, field)
			}
		}
		body := contenido[frontmatterEnd+9:]
		for _, section := range []string{"## Por qué existe", "## Instrucciones", "## Ejemplos", "## Criterios de aceptación"} {
			if !strings.Contains(body, section) {
				t.Errorf("%s no contiene %s", nombre, section)
			}
		}
	}
}

func TestCargarSkillsAsesora(t *testing.T) {
	s, err := CargarSkills("asesora")
	if err != nil {
		t.Fatal(err)
	}
	for _, nombre := range mapaSkills["asesora"] {
		if !strings.Contains(s, "=== SKILL: "+nombre+" ===") {
			t.Errorf("falta %s", nombre)
		}
	}
	if strings.Contains(s, "fuente_de_verdad:") || strings.Contains(s, "agente:") {
		t.Error("el frontmatter no debe inyectarse al prompt")
	}
	if !strings.Contains(s, "## Instrucciones") {
		t.Error("el cuerpo de la skill debe inyectarse")
	}
}

func TestCargarSkillsAgentesYAgenteDesconocido(t *testing.T) {
	for agente, nombres := range mapaSkills {
		contenido, err := CargarSkills(agente)
		if err != nil {
			t.Fatalf("%s: %v", agente, err)
		}
		for _, nombre := range nombres {
			if !strings.Contains(contenido, "=== SKILL: "+nombre+" ===") {
				t.Errorf("%s no carga %s", agente, nombre)
			}
		}
	}
	contenido, err := CargarSkills("desconocido")
	if err != nil {
		t.Fatal(err)
	}
	if contenido != "" {
		t.Fatalf("agente desconocido produjo contenido: %q", contenido)
	}
}
