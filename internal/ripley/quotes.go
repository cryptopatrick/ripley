// Package ripley provides Ripley-style quotes categorized by effort level.
// Inspired by Ellen Ripley from Alien: calm, procedural, and rule-following.
package ripley

import (
	"math/rand"
	"time"
)

// Categorized Ripley-style quotes for different effort levels

// GoodEffortQuotes are used when the AI performs well within expected parameters.
var GoodEffortQuotes = []string{
	"Now we're getting somewhere — that's the baseline competence I expect.",
	"Not bad — you actually *tried*.",
	"Finally, some effort worth reporting.",
	"Precision, speed, and efficiency — looks like someone's awake.",
	"That's how we do it — no shortcuts, no excuses.",
	"You followed procedure. I approve.",
	"Effort noted and logged. Keep it consistent.",
}

// MediumEffortQuotes are used when performance is acceptable but not optimal.
var MediumEffortQuotes = []string{
	"Hmm… you're getting there, but don't think I won't notice the shortcuts.",
	"I see the work, but it's half-baked.",
	"Mediocre, but technically acceptable.",
	"Are you actually trying, or just going through the motions?",
	"Not the worst, but I've seen better from a fresh context.",
	"Decent. I'll allow it… this time.",
	"I can work with this, though it smells like token padding.",
	"Average effort — report submitted but not impressive.",
	"You've barely scratched the surface of competence.",
	"Followed the rules, but where's the spark?",
}

// PoorEffortQuotes are used when the AI fails or performs poorly.
var PoorEffortQuotes = []string{
	"What is this? Did you even read the instructions?",
	"I don't trust claims — prove it, properly.",
	"Tokens wasted, time wasted, and effort barely measurable.",
	"You're phoning it in, and I can tell.",
	"This is not baseline competence — start over.",
	"Airlock checks first, answers second. You failed both.",
	"Sloppy, incomplete, and unnecessary verbosity.",
	"I expected more from Sonnet 4.5 than this.",
	"Are you trying to look busy? Because it's not working.",
	"Do better. Or I will notice.",
	"Half-hearted response logged. Not acceptable.",
	"I've seen lower-effort outputs from a malfunctioning interface, and this is close.",
	"Do not waste tokens pretending — this is your warning.",
}

// Initialize random seed once at package initialization
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomQuoteByEffort returns a random Ripley-style quote for the given effort category.
// Valid effort values: "good", "medium", "poor"
// Returns a random quote from all categories if effort is not recognized.
func RandomQuoteByEffort(effort string) string {
	switch effort {
	case "good":
		return GoodEffortQuotes[rand.Intn(len(GoodEffortQuotes))]
	case "medium":
		return MediumEffortQuotes[rand.Intn(len(MediumEffortQuotes))]
	case "poor":
		return PoorEffortQuotes[rand.Intn(len(PoorEffortQuotes))]
	default:
		// Fallback: return a random quote from any category
		all := append(append(GoodEffortQuotes, MediumEffortQuotes...), PoorEffortQuotes...)
		return all[rand.Intn(len(all))]
	}
}
