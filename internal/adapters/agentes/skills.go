package agentes

import (
	"embed"
	"fmt"
	"strings"

	skillassets "github.com/Unikyri/vivi-perfilamiento-leads/skills"
)

// skillsFS is provided by the compatible-root asset package in skills/.
// Keeping this value as embed.FS preserves the original package-local shape.
var skillsFS embed.FS = skillassets.FS

// mapaSkills: qué skills carga cada agente (doc 05 §5).
var mapaSkills = map[string][]string{
	"asesora": {
		"tono-colsubsidio", "normalizacion-de-declarados", "siguiente-mejor-pregunta",
		"explicacion-financiera-humana", "presentacion-de-proyectos",
	},
	"investigadora": {"dominio-caja"},
	"nutricionista": {"tono-colsubsidio", "redaccion-con-dignidad", "explicacion-financiera-humana"},
	"documentadora": {"ficha-comercial"},
}

// CargarSkills concatena el cuerpo de las skills de un agente para inyectarlo al prompt.
func CargarSkills(agente string) (string, error) {
	var b strings.Builder
	for _, nombre := range mapaSkills[agente] {
		datos, err := skillsFS.ReadFile(nombre + "/SKILL.md")
		if err != nil {
			return "", fmt.Errorf("leyendo skill %s: %w", nombre, err)
		}
		b.WriteString("\n\n=== SKILL: " + nombre + " ===\n")
		b.WriteString(quitarFrontmatter(string(datos)))
	}
	return b.String(), nil
}

// quitarFrontmatter elimina el bloque YAML: es metadata, no instrucciones.
func quitarFrontmatter(s string) string {
	if !strings.HasPrefix(s, "---") {
		return s
	}
	if i := strings.Index(s[3:], "---"); i >= 0 {
		return strings.TrimSpace(s[i+6:])
	}
	return s
}
