package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const schemaWeekURL = "https://api.anilibria.app/api/v1/anime/schedule/week"

type scheduleItem struct {
	Release release `json:"release"`
}

type release struct {
	Name struct {
		Main    string `json:"main"`
		English string `json:"english"`
	} `json:"name"`

	PublishDay *struct {
		Value       int    `json:"value"`
		Description string `json:"description"`
	} `json:"publish_day"`

	IsOngoing bool `json:"is_ongoing"`

	EpisodesTotal *int `json:"episodes_total"`

	NextReleaseEpisodeNumber *int `json:"next_release_episode_number"`

	LatestEpisode *struct {
		Ordinal int `json:"ordinal"`
	} `json:"latest_episode"`
}

func BuildSchemaText(ctx context.Context, loc *time.Location) (string, error) {
	if loc == nil {
		loc = time.Local
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, schemaWeekURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "tg-bot/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("anilibria status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var items []scheduleItem
	if err := json.Unmarshal(raw, &items); err != nil {
		var wrapped struct {
			Data []scheduleItem `json:"data"`
		}
		if err2 := json.Unmarshal(raw, &wrapped); err2 != nil {
			return "", err
		}
		items = wrapped.Data
	}

	// dayValue -> titles(set)
	byDay := make(map[int]map[string]struct{}, 7)
	dayName := make(map[int]string, 7)


for _, it := range items {
		r := it.Release

		if r.PublishDay == nil {
			continue
		}

		d := r.PublishDay.Value
		if d < 1 || d > 7 {
			continue
		}

		title := strings.TrimSpace(r.Name.Main)
		if title == "" {
			title = strings.TrimSpace(r.Name.English)
		}
		if title == "" {
			continue
		}

		// формируем “красивую” строку с серией
		title = formatTitleWithEpisode(r, title)

		if byDay[d] == nil {
			byDay[d] = map[string]struct{}{}
		}
		byDay[d][title] = struct{}{}

		if r.PublishDay.Description != "" {
			dayName[d] = r.PublishDay.Description
		}
	}

	now := time.Now().In(loc)
	weekStart := startOfISOWeek(now)
	weekEnd := weekStart.AddDate(0, 0, 6)

	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("🗓 Расписание на неделю %s–%s\n",
		weekStart.Format("02.01.2006"), weekEnd.Format("02.01.2006")))
	b.WriteString(fmt.Sprintf("📍 Сегодня: %s · %s\n\n",
		russianWeekdayName(isoDayValue(now)), now.Format("02.01.2006")))

	for d := 1; d <= 7; d++ {
		dayDate := weekStart.AddDate(0, 0, d-1) // 00:00
		diff := int(dayDate.Sub(todayStart).Hours() / 24)

		name := dayName[d]
		if name == "" {
			name = russianWeekdayName(d)
		}

		badge := dayBadge(diff)
		b.WriteString(fmt.Sprintf("━━━━━━━━━━━━━━\n📌 %s · %s  %s\n",
			name, dayDate.Format("02.01.2006"), badge))

		var titles []string
		for t := range byDay[d] {
			titles = append(titles, t)
		}
		sort.Strings(titles)

		if len(titles) == 0 {
			b.WriteString("— нет релизов\n")
			continue
		}

		for _, t := range titles {
			b.WriteString("• " + t + "\n")
		}
	}

	return strings.TrimSpace(b.String()), nil
}
