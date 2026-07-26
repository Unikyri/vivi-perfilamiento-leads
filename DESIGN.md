# DESIGN.md — Sistema de Diseño Vivi Colsubsidio

---
name: Vivi — Advisor Panel (#panel-asesor)
description: The calm lead workspace — a rounded-card system for the Colsubsidio housing advisor's dashboard, redesigned from a pinned reference.
colors:
  colsubsidio-azul: "#003DA6"
  colsubsidio-azul-oscuro: "#002D7A"
  colsubsidio-azul-claro: "#0056B3"
  colsubsidio-amarillo: "#FFC700"
  colsubsidio-amarillo-suave: "#FFF3C4"
  ok: "#2E8540"
  ok-ink: "#256B34"
  alerta: "#E8A200"
  alerta-ink: "#8A5A00"
  neutro: "#9E9E9E"
  gris-ink: "#5B5B5B"
  error: "#C0392B"
  papel: "#FFFFFF"
  papel-suave: "#F5F4EF"
  fondo: "#FAFAF8"
  zona-verde-bg: "#EAF6EC"
  zona-ambar-bg: "#FDF3DF"
  zona-gris-bg: "#F1F1EF"
  texto-principal: "#1F2937"
  texto-secundario: "#4B5563"
  borde-suave: "#E5E7EB"
  regla-fuerte: "#C9CDD3"
typography:
  headline:
    fontFamily: "'Inter', system-ui, -apple-system, sans-serif"
    fontSize: "1.5rem"
    fontWeight: 800
    lineHeight: 1.2
    letterSpacing: "normal"
  title:
    fontFamily: "'Inter', system-ui, -apple-system, sans-serif"
    fontSize: "0.85rem"
    fontWeight: 800
    lineHeight: 1.3
    letterSpacing: "0.04em"
  body:
    fontFamily: "'Inter', system-ui, -apple-system, sans-serif"
    fontSize: "0.85rem"
    fontWeight: 400
    lineHeight: 1.5
    letterSpacing: "normal"
  label:
    fontFamily: "'Inter', system-ui, -apple-system, sans-serif"
    fontSize: "0.7rem"
    fontWeight: 700
    lineHeight: 1.3
    letterSpacing: "0.04em"
  cifra:
    fontFamily: "'Inter', system-ui, -apple-system, sans-serif"
    fontSize: "0.85rem"
    fontWeight: 600
    lineHeight: 1.3
    letterSpacing: "normal"
rounded:
  chip: "8px"
  card: "14px"
  pill: "999px"
spacing:
  xs: "0.3rem"
  sm: "0.6rem"
  md: "1rem"
  lg: "1.5rem"
  xl: "2rem"
components:
  hoja-informe:
    backgroundColor: "{colors.papel}"
    textColor: "{colors.texto-principal}"
    rounded: "{rounded.card}"
    padding: "1.5rem 1.75rem"
  fila-lead:
    backgroundColor: "{colors.papel-suave}"
    rounded: "{rounded.card}"
    padding: "0.85rem 1rem"
  banner-siguiente-paso:
    backgroundColor: "{colors.colsubsidio-azul}"
    textColor: "{colors.blanco}"
    rounded: "{rounded.card}"
    padding: "1rem 1.25rem"
  btn-copiar-resumen:
    backgroundColor: "{colors.colsubsidio-amarillo}"
    textColor: "{colors.colsubsidio-azul-oscuro}"
    rounded: "{rounded.chip}"
    padding: "0.5rem 1rem"
  btn-copiar-resumen-hover:
    backgroundColor: "#E6B200"
  sello-verificado:
    backgroundColor: "{colors.zona-verde-bg}"
    textColor: "{colors.ok-ink}"
    rounded: "{rounded.pill}"
    padding: "0.14rem 0.55rem"
  sello-declarado:
    backgroundColor: "#EAF0FB"
    textColor: "{colors.colsubsidio-azul}"
    rounded: "{rounded.pill}"
    padding: "0.14rem 0.55rem"
  sello-inferido:
    backgroundColor: "{colors.zona-ambar-bg}"
    textColor: "{colors.alerta-ink}"
    rounded: "{rounded.pill}"
    padding: "0.14rem 0.55rem"
  tab:
    backgroundColor: "transparent"
    textColor: "{colors.texto-secundario}"
    padding: "0.75rem 1.25rem"
  tab-active:
    backgroundColor: "{colors.papel-suave}"
    textColor: "{colors.colsubsidio-azul}"
---

# Design System: Vivi — Advisor Panel

## Overview

**Creative North Star: "The Lead Workspace"**

The advisor panel is a calm, rounded-card workspace the advisor scans mid-call — not an audited paper report. This redesign replaces the previous "Auditor's Ledger" register (flat paper, square corners, IBM Plex Mono figures) with soft white cards floating on a pale panel ground, following a reference screenshot the user pinned directly: a clean lead-queue/ficha dashboard with avatar-initial cards, pill badges, and a single humanist sans typeface throughout. The Colsubsidio blue/amarillo brand and logo are unchanged — this redesign is about card language and typography, not corporate identity.

This is still an Operate-mode surface: legibility and scan speed win over expression. The card language serves that goal by grouping each lead's identity, status, and action into one visually contained unit the eye can jump between, instead of parsing table-like rows.

**Key Characteristics:**
- Rounded white cards (`--radio-card`, 14px) on a soft off-white/blue-gray ground — no flat paper-report register, no square corners.
- Single typeface (`Inter`) for everything, including figures — hierarchy comes from weight/size, not a mono/sans split. `.cifra` keeps `tabular-nums` for column alignment, but no longer switches font family.
- Semaphore status shown twice, redundantly but subtly: a small colored dot inline in a pill label, and the pill's own tint — never a thick accent border (that reads as generic AI-slop card decoration, and was deliberately removed after the mechanical detector flagged it).
- Avatar-initial circles give each lead a visual anchor, matching the pinned reference.
- Provenance stamps and the "one amarillo per view" rule carry over unchanged from the previous system — both were sound and not tied to the paper-report metaphor.
- Color stays restrained: blue for identity/primary actions, amarillo for the single next-step, the semaphore trio for status only. Not "drenched" — cards are white/off-white, color is spent on accents and pills, never a card's whole background.

## Colors

Unchanged from the previous system — this redesign did not touch the Colsubsidio brand or semaphore palette, only the shapes and typography that carry it.

### Primary
- **Colsubsidio Azul** (`#003DA6`): top bar background, section titles (`.seccion-titulo`), masthead names, active tab text, banner and gauge fills. The dominant color of the system.
- **Colsubsidio Azul Oscuro** (`#002D7A`): pressed/deep variant — demo badge text, `btn-copiar-resumen` text on amarillo.
- **Colsubsidio Azul Claro** (`#0056B3`): the confidence-gauge fill, a lighter step of the same ink.

### Secondary
- **Colsubsidio Amarillo** (`#FFC700`): reserved for exactly one active action per view — the "Copiar resumen" button. Not used for badges or highlights elsewhere.
- **Amarillo Suave** (`#FFF3C4`): the disabled state of that same button only.

### Neutral
- **Papel** (`#FFFFFF`): card background (`.hoja-informe`, `.barra-track`, avatar-adjacent surfaces).
- **Papel Suave** (`#F5F4EF`): the tinted surface for lead-queue cards (`.fila-lead`), active-tab background, gauge tracks, and chart-column cards (`.columna-grafico`) — a warm off-white that separates a card from the pure-white sheet behind it.
- **Fondo** (`#FAFAF8`): the outer panel background behind all cards.
- **Texto Principal / Secundario / Borde Suave / Regla Fuerte**: unchanged — primary text, secondary/caption text, soft dividers, and the stronger masthead/table rule.

### Semaphore fill/border (verde, ámbar, gris)
Unchanged: **Ok/Verde** (`#2E8540`, zone bg `#EAF6EC`), **Alerta/Ámbar** (`#E8A200`, zone bg `#FDF3DF`), **Neutro/Gris** (`#9E9E9E`, zone bg `#F1F1EF`), **Error** (`#C0392B`, system errors only).

### Named Rules
**The Ink-Not-Fill Rule** (unchanged): semaphore text always uses the darkened `-ink` counterpart; the base token is reserved for fills, borders, and the status dot.

**The One Amarillo Rule** (unchanged): amarillo appears on exactly one element per view, the next-step action inside `.banner-siguiente-paso`.

## Typography

**Font:** `Inter` (weights 400/500/600/700/800, loaded via Google Fonts in `index.html`), with `system-ui, -apple-system, sans-serif` fallback — applied globally (`html, body`), not scoped to one panel.

**Character:** One family carries the whole surface, matching the pinned reference and the Operate-mode guidance that workhorse system-adjacent sans faces serve dashboards well. Hierarchy comes entirely from weight and size, never a font-family switch. `.cifra` keeps `font-variant-numeric: tabular-nums` so scores/amounts/dates still align in a column, but no longer swaps to a mono family — the ledger-instrument read was retired along with the paper-report metaphor.

### Hierarchy
- **Headline** (800, 1.4rem–1.9rem): masthead names and the cola's lead-count figure.
- **Title** (800, 0.85rem, uppercase, 0.04em tracking): `.seccion-titulo`, blue section headers.
- **Subtitle** (700, 0.8rem, uppercase): `.subseccion-titulo`, gray instead of blue.
- **Body** (400–600, 0.85rem–0.98rem): card prose, resumen lines, observation lists.
- **Label** (600–700, 0.6rem–0.75rem, uppercase): field labels, gauge labels, table headers.
- **Cifra** (600–800, tabular-nums): applied to every score/amount/date/id, same rule as before, just no font-family switch.

## Layout

Unchanged shell: `.split` (`#panel-chat` | `#panel-asesor`) stacks vertically below 768px, chat first then advisor. New addition: **collapsible panels** — a small chevron button (`.btn-colapsar`) sits on each side of the seam between the two panels; clicking one collapses that panel to zero width (`.panel-colapsado`) and the sibling panel grows to fill the freed space via `:has()` (`.split:has(#panel-asesor.panel-colapsado) #panel-chat`). Hidden below 768px, where panels already stack full-width.

Every tab (Cola, Ficha, Gerencia) still renders inside a `.hoja-informe` card, but it is now one rounded card among several rather than a single flat sheet: lead-queue rows (`.fila-lead`), chart columns (`.columna-grafico`), and the verdict zone (`.zona-veredicto`) are each their own rounded card, not a ledger row or a bordered zone inside one continuous sheet.

At ≤768px, `.fila-lead` collapses to a 3-column/2-row grid (`avatar` / `cuerpo` / `prioridad`+`accion`), `.lead-ruta` hides, and the collapse-toggle buttons hide.

## Elevation & Depth

Cards use a soft double shadow (`--sombra-card`) for ambient lift, and lead-queue rows get a stronger shadow (`--sombra-fila`) on hover/focus. This is a genuine shift from the old "one sheet, one shadow" rule: every card-like unit (queue row, chart column, verdict zone) now carries its own soft shadow, because the workspace is explicitly a collection of cards, not a single report page with bordered sections.

## Shapes

Rounded is the default: `--radio-card` (14px) on the sheet, queue rows, chart columns, and the verdict/banner zones; `--radio-chip` (8px) on tab corners and the warning-banner label; `--radio-pill` (999px) on avatars, status pills, provenance stamps, the afiliado badge, and buttons (`btn-ver-chat`, `select-lead-demo`-style chrome). Square corners were the previous system's load-bearing trait; this redesign inverts that on purpose — the pinned reference is a soft, rounded dashboard, and the user asked for that register while keeping brand colors intact.

**The No-Side-Accent Rule.** Lead-queue cards do not carry a colored left border for semaphore state — an earlier draft did, and the mechanical design detector flagged it as the canonical AI-generated-UI tell (`side-tab` / thick left accent on a card). Status is communicated instead by a small dot + colored pill label (`.fila-lead-zona`) and, for affiliation, a separate muted/blue pill (`.fila-lead-afiliado`) — both legible without relying on a border cliché.

## Components

### Buttons
- **Primary (`.btn-copiar-resumen`):** amarillo background, azul-oscuro text, 8px radius — unchanged as the one amarillo element per view.
- **Row action (`.btn-ver-chat`):** WhatsApp-register green, now pill-shaped (`--radio-pill`) to match the rounded card language; still an intentional exception to the panel's blue/amarillo palette because it signals the WhatsApp hand-off destination.
- **Tabs (`.tabs button`):** transparent background, rounded top corners, azul bottom-border + `papel-suave` fill when active.
- **Collapse toggle (`.btn-colapsar`):** small chevron button on the panel seam, flips via `scaleX` when its panel is collapsed; lives in `base.css` as shell chrome, not part of the card system proper.

### Chips / Pills (signature component)
- **Provenance stamps (`.sello-fuente`):** now pill-shaped (`--radio-pill`) instead of the old near-square "stamp." Border-style distinction is preserved — double border (verificado), solid (declarado), dashed (inferido) — so confidence reads without the label.
- **Zona pill (`.fila-lead-zona`):** a small colored dot + uppercase label, no background/border — the semaphore signal for queue cards.
- **Afiliado pill (`.fila-lead-afiliado`):** new component. Shows only "Afiliado" / "No afiliado" on queue cards, sourced from `LeadEnCola.afiliado` — deliberately does **not** show a segment/category, because `GET /api/leads` does not return one; category only exists on the full `Ficha` object and stays in `ficha.ts`'s identity line ("Afiliado · Categoría X").
- **Zona badge on the ficha verdict (`.zona-veredicto`):** still a full tinted, bordered region — now with rounded corners — since a whole-card verdict band is not the same "side accent" pattern the No-Side-Accent Rule targets.

### Cards / Containers
- **Corner Style:** rounded (`--radio-card`, 14px).
- **Background:** `--papel` for the outer sheet, `--papel-suave` for queue rows and chart-column cards.
- **Shadow Strategy:** every card gets `--sombra-card`; queue rows add `--sombra-fila` on hover/focus.
- **Avatar:** `--radio-pill` circle, initial letter, `--colsubsidio-azul-claro` fill.

### Inputs / Fields
- **Style:** unchanged from the previous system — borderless selects in a bordered/pill wrapper (`.selector-proyecto-wrap`, now pill-shaped).
- **Focus:** unchanged, fixed accessibility commitment — `outline: 2px solid var(--colsubsidio-azul)` on `:focus-visible`.

### Instrument Gauges
`.medidor-cupo` and `.medidor-confianza` keep their calibrated-scale/tick-mark structure, now with rounded, overflow-hidden tracks instead of square ones. `.barra-fill` (Gerencia) same treatment. All fills still animate via `transform: scaleX()`, never `width`.

## Do's and Don'ts

### Do:
- **Do** use `.cifra` (`tabular-nums`) for every score, amount, date, id, and percentage — the alignment rule survives even though the font-family switch does not.
- **Do** use the `-ink` semaphore variant for any text on a semaphore-colored surface.
- **Do** give every card-like unit (`hoja-informe`, `fila-lead`, `columna-grafico`, `zona-veredicto`, `banner-siguiente-paso`) rounded corners and its own soft shadow.
- **Do** treat the Colsubsidio blue/amarillo/logo brand system as fixed — this redesign changed shape and type, never brand color or the logo.
- **Do** keep amarillo to exactly one active-action element per view.
- **Do** show semaphore state as a dot + pill label, never a border accent.

### Don't:
- **Don't** use `.panel-chat` / WhatsApp visual patterns as a reference for the advisor panel — still out of scope; `.btn-ver-chat`'s WhatsApp-green remains the one confirmed exception.
- **Don't** add a colored left/side border to a card for status — this is the specific anti-pattern this redesign removed; use the dot+pill instead.
- **Don't** show lead segment/category (`Básico`/`Joven`/`Medio`, etc.) on queue cards — `GET /api/leads` doesn't return it; only the full `Ficha` view has `categoria`.
- **Don't** reintroduce a second (mono) font family for figures — one Inter family carries the whole surface now.
- **Don't** square off card corners — rounded is load-bearing for this register, the inverse of the previous system.
