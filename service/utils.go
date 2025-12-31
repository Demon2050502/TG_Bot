package service

import (
	"fmt"
	"time"
)

func isoDayValue(t time.Time) int {
	wd := int(t.Weekday())
	if wd == 0 {
		return 7
	}
	return wd
}


func splitByLimit(s string, limit int) []string {
	if limit <= 0 || len(s) <= limit {
		return []string{s}
	}
	var res []string
	for len(s) > limit {
		res = append(res, s[:limit])
		s = s[limit:]
	}
	if len(s) > 0 {
		res = append(res, s)
	}
	return res
}

func formatTitleWithEpisode(r release, title string) string {
	// Лучшее для schedule/week — next_release_episode_number (что выйдет)
	if r.NextReleaseEpisodeNumber != nil && *r.NextReleaseEpisodeNumber > 0 {
		if r.EpisodesTotal != nil && *r.EpisodesTotal > 0 {
			return fmt.Sprintf("%s — выйдет серия %d из %d", title, *r.NextReleaseEpisodeNumber, *r.EpisodesTotal)
		}
		return fmt.Sprintf("%s — выйдет серия %d", title, *r.NextReleaseEpisodeNumber)
	}

	// fallback: если вдруг есть latest_episode
	if r.LatestEpisode != nil && r.LatestEpisode.Ordinal > 0 {
		if r.EpisodesTotal != nil && *r.EpisodesTotal > 0 {
			return fmt.Sprintf("%s — серия %d из %d", title, r.LatestEpisode.Ordinal, *r.EpisodesTotal)
		}
		return fmt.Sprintf("%s — серия %d", title, r.LatestEpisode.Ordinal)
	}

	return title
}

func dayBadge(diff int) string {
	switch {
	case diff == 0:
		return "🟩 СЕГОДНЯ"
	case diff == 1:
		return "🟦 ЗАВТРА"
	case diff == -1:
		return "✅ ВЧЕРА"
	case diff < 0:
		return "✅ прошло " + pluralDays(-diff) + " назад"
	default:
		return "⏳ через " + pluralDays(diff)
	}
}

func pluralDays(n int) string {
	// 1 день, 2/3/4 дня, 5+ дней, 11-14 дней
	nMod100 := n % 100
	if nMod100 >= 11 && nMod100 <= 14 {
		return fmt.Sprintf("%d дней", n)
	}
	switch n % 10 {
	case 1:
		return fmt.Sprintf("%d день", n)
	case 2, 3, 4:
		return fmt.Sprintf("%d дня", n)
	default:
		return fmt.Sprintf("%d дней", n)
	}
}

func startOfISOWeek(t time.Time) time.Time {
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	dayStart := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return dayStart.AddDate(0, 0, -(wd - 1))
}

func russianWeekdayName(d int) string {
	switch d {
	case 1:
		return "Понедельник"
	case 2:
		return "Вторник"
	case 3:
		return "Среда"
	case 4:
		return "Четверг"
	case 5:
		return "Пятница"
	case 6:
		return "Суббота"
	case 7:
		return "Воскресенье"
	default:
		return "?"
	}
}
