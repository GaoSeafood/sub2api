package service

import "strings"

var claude55ThinkingVariantBases = map[string]struct{}{
	"claude-opus-5-5":   {},
	"claude-sonnet-5-5": {},
}

// resolveClaude55ThinkingVariantForLevel lets a base Claude 5.5 model select
// its Antigravity upstream variant from the request's reasoning preference.
// A request that already names a tier remains pinned to that tier.
func resolveClaude55ThinkingVariantForLevel(account *Account, requestedModel, preferred string) (string, bool) {
	if account == nil {
		return "", false
	}
	model := strings.TrimSpace(strings.TrimPrefix(requestedModel, "models/"))
	if _, ok := claude55ThinkingVariantBases[model]; !ok {
		return "", false
	}

	mapping := account.GetModelMapping()
	if len(mapping) == 0 {
		return "", false
	}
	baseMapped, baseMatched := resolveRequestedModelInMapping(mapping, model)
	baseMapped = strings.TrimSpace(baseMapped)
	if baseMatched && baseMapped != model && baseMapped != model+"-high" {
		// Preserve explicit aliases and pinned low/medium mappings. The default
		// base -> high mapping can be overridden by the request's effort.
		return "", false
	}

	switch strings.ToLower(strings.TrimSpace(preferred)) {
	case "low", "minimal", "none":
		preferred = "low"
	case "medium":
		preferred = "medium"
	default:
		preferred = "high"
	}
	order := []string{preferred}
	for _, level := range []string{"high", "medium", "low"} {
		if level != preferred {
			order = append(order, level)
		}
	}
	for _, level := range order {
		candidate := model + "-" + level
		var (
			mapped  string
			matched bool
		)
		if baseMatched {
			mapped, matched = mapping[candidate]
		} else {
			mapped, matched = resolveRequestedModelInMapping(mapping, candidate)
		}
		if matched && strings.TrimSpace(mapped) != "" {
			return strings.TrimSpace(mapped), true
		}
	}
	return "", false
}
