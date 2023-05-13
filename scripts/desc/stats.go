package main

import (
	"dbutils/pkg/desc"
	"dbutils/pkg/file"
	"dbutils/pkg/stat"
	"encoding/json"
	"log"
	"os"
	"strings"
)

var hackEnStatDescContentEntries = [][2]string{
	{`#|60 "Gain {0} Vaal Soul Per Second during effect" per_minute_to_per_second 1`,
		`60 "Gain {0} Vaal Soul Per Second during effect" per_minute_to_per_second 1`},
	{`1|# "[DNT] Area contains {0} additional Common Chest Marker"`,
		`1 "[DNT] Area contains {0} additional Common Chest Marker"`},
}

func hackEnStatDescContent(content string) string {
	for _, entry := range hackEnStatDescContentEntries {
		if strings.Contains(content, entry[0]) {
			content = strings.ReplaceAll(content, entry[0], entry[1])
		} else {
			log.Printf("hack missed: %v", entry[0])
		}
	}
	return content
}

var hackZhStatDescContentEntries = [][2]string{
	{`#|60 "生效期间每秒获得 {0} 个瓦尔之灵" per_minute_to_per_second 1`,
		`60 "生效期间每秒获得 {0} 个瓦尔之灵" per_minute_to_per_second 1`},
	{`#|-1 "能量护盾全满状态下防止{0:+d}%的被压制法术伤害" reminderstring ReminderTextSuppression`,
		`#|-1 "能量护盾全满状态下防止{0:+d}%的被压制法术伤害的总量" reminderstring ReminderTextSuppression`},
	{`#|-1 "枯萎技能会使干扰持续时间延长 {0}%" negate 1`,
		`#|-1 "枯萎技能会使干扰持续时间缩短 {0}%" negate 1`},
	{`#|-1 "【寒霜爆】技能会使减益效果的持续时间延长 {0}%" negate 1`,
		`#|-1 "【寒霜爆】技能会使减益效果的持续时间缩短 {0}%" negate 1`},
	{`#|-1 "每 10 秒获得 {0}% 的元素伤害增益，持续 4 秒" negate 1`,
		`#|-1 "每 10 秒获得 {0}% 的元素伤害减益，持续 4 秒" negate 1`},
	{`#|-1 "若腐化，则全域暴击率提高 {0}%" negate 1 reminderstring ReminderTextIfItemCorrupted`,
		`#|-1 "若腐化，则全域暴击率降低 {0}%" negate 1 reminderstring ReminderTextIfItemCorrupted`},
	{`1|# "[DNT]该区域会额外出现{0}个普通宝箱标记"`,
		`1 "[DNT]该区域会额外出现{0}个普通宝箱标记"`},
	{"\t1\r\n\t\t# \"图腾放置速度加快 {0}%\"\r\n",
		"\t2\r\n\t\t1|# \"图腾放置速度加快 {0}%\"\r\n#|-1 \"图腾放置速度减慢 {0}%\"\r\n"},
	{"\t\t1|# \"若你近期内有击败敌人，则效果区域扩大 {0}%，最多 50%\" reminderstring ReminderTextRecently\r\n\t\t#|-1 \"若你近期内有击败敌人，则效果区域缩小 {0}%\" negate 1  reminderstring ReminderTextRecently\r\n",
		"\t\t1|# \"若你近期内有击败敌人，则效果区域扩大 {0}%，最多 50%\" reminderstring ReminderTextRecently\r\n\t\t#|-1 \"若你近期内有击败敌人，则效果区域缩小 {0}%，最多 50%\" negate 1  reminderstring ReminderTextRecently\r\n"},
	{"【毒雨】可以额外发射 1 个箭矢", "【毒雨】可以额外发射 {0} 个箭矢"},
	{`1|# "如果诅咒持续时间已经过去 25%，\n则你诅咒的敌人的移动速度被减缓 25%"`, `1|# "如果诅咒持续时间已经过去 25%，\n则你诅咒的敌人的移动速度被减缓 {0}%"`},
}

func hackZhStatDescContent(content string) string {
	for _, entry := range hackZhStatDescContentEntries {
		if strings.Contains(content, entry[0]) {
			content = strings.ReplaceAll(content, entry[0], entry[1])
		} else {
			log.Printf("hack missed: %v", entry[0])
		}
	}
	return content
}

func hackDescs(descs []*desc.Desc) {
	for _, d := range descs {
		// 雷鸣洗礼的`受到的冰霜伤害的 40% 转为火焰伤害`，应当为`受到的冰霜伤害的 40% 转为闪电伤害`
		// 与`受到的冰霜伤害的 {0}% 转为火焰伤害`冲突
		// 需要translator进行hack
		if d.Id == "cold_hit_and_dot_damage_%_taken_as_lightning" {
			if d.Texts[desc.LangZh][0].Template == "受到的冰霜伤害的 {0}% 转为火焰伤害" {
				d.Texts[desc.LangZh][0].Template = "受到的冰霜伤害的 {0}% 转为闪电伤害"
			} else {
				log.Printf("hack missed: %v", d.Id)
			}
			continue
		}

		// 血影的`每个狂怒球可使攻击速度减慢 4%`，应当为`每个狂怒球可使攻击和施法速度减慢 4%`
		// 与`每个狂怒球可使攻击速度加快 {0}%`,`每个狂怒球可使攻击速度减慢 {0}%`冲突
		// 需要translator进行hack
		if d.Id == "attack_and_cast_speed_+%_per_frenzy_charge" {
			if d.Texts[desc.LangZh][0].Template == "每个狂怒球可使攻击速度加快 {0}%" {
				d.Texts[desc.LangZh][0].Template = "每个狂怒球可使攻击和施法速度加快 {0}%"
			} else {
				log.Panicf("hack missed: %v", d.Id)
			}
			if d.Texts[desc.LangZh][1].Template == "每个狂怒球可使攻击速度减慢 {0}%" {
				d.Texts[desc.LangZh][1].Template = "每个狂怒球可使攻击和施法速度减慢 {0}%"
			} else {
				log.Panicf("hack missed: %v", d.Id)
			}
			continue
		}

		// 戴亚迪安的晨曦的`没有物理伤害`，应当为`不造成物理伤害`
		// 与武器上的`没有物理伤害`词缀产生冲突
		// 受影响物品：戴亚迪安的晨曦，异度天灾武器基底词缀
		// 需要translator进行hack
		if d.Id == "base_deal_no_physical_damage" {
			if d.Texts[desc.LangZh][0].Template == "没有物理伤害" {
				d.Texts[desc.LangZh][0].Template = "不造成物理伤害"
			}
		}
	}
}

var skipedDescIds = map[string]bool{
	// 中文相同英文不同，是地图词缀
	"map_projectile_speed_+%":                     true,
	"map_players_gain_soul_eater_on_rare_kill_ms": true,
	// 中文相同英文不同，是局部词缀
	"local_gem_experience_gain_+%": true,
	"local_accuracy_rating_+%":     true,
	// 中文相同英文不同，是技能说明
	"skill_range_+%": true,
	// 中文相同英文不同，是无效词缀
	"elemental_damage_taken_+%_during_flask_effect": true,
	"global_poison_on_hit":                          true,
	"bleed_on_melee_critical_strike":                true,
	// 中文相同英文不同，但是英文均为有效词缀
	"curse_on_hit_level_warlords_mark":                        true,
	"damage_taken_+%_if_you_have_taken_a_savage_hit_recently": true,
	"immune_to_bleeding":                                      true,
	"onslaught_buff_duration_on_kill_ms":                      true,
	// 中文相同英文不同，不知道正确词缀
	// 【断金之刃】的伤害提高，【断金之刃】的伤害降低
	"shattering_steel_damage_+%": true,
	"lancing_steel_damage_+%":    true,
}

func removeSkipedDesc(descs []*desc.Desc) []*desc.Desc {
	newDescs := make([]*desc.Desc, 0, len(descs))
	for _, d := range descs {
		if !skipedDescIds[d.Id] &&
			!strings.HasPrefix(d.Id, "map_") {
			newDescs = append(newDescs, d)
		}
	}

	return newDescs
}

func checkDuplicateZh(stats []*stat.Stat) {
	records := map[string]string{}

	for _, stat := range stats {
		if recordEn, ok := records[stat.Zh]; ok {
			if !strings.EqualFold(recordEn, stat.En) {
				log.Printf("warning diff en of: %v", stat.Zh)
				log.Print(recordEn)
				log.Print(stat.En)
			}
		} else {
			records[stat.Zh] = stat.En
		}
	}
}

func main() {
	enStatDescFile := "../../docs/desc/en_stat_descriptions.txt"
	zhStatDescFile := "../../docs/desc/zh_stat_descriptions.txt"

	enStatDescContent := file.ReadFileUTF16(enStatDescFile)
	zhStatDescContent := file.ReadFileUTF16(zhStatDescFile)

	enStatDescContent = hackEnStatDescContent(enStatDescContent)
	zhStatDescContent = hackZhStatDescContent(zhStatDescContent)

	descs := desc.Load(strings.Split(enStatDescContent, "\r\n"), strings.Split(zhStatDescContent, "\r\n"))
	descs = removeSkipedDesc(descs)
	hackDescs(descs)

	stats := desc.ToStats(descs)
	checkDuplicateZh(stats)

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("../../src/stats/desc.json", data, 0666)
}
