package storage

func ListUpgradeUnlockRules() []UpgradeUnlockRule {
	rules := []UpgradeUnlockRule{
		{Category: "刹车", UpgradeName: "赛车版刹车", Unlocks: "解锁刹车调校"},
		{Category: "弹簧与阻尼器", UpgradeName: "赛车版弹簧与阻尼器", Unlocks: "解锁弹簧、阻尼器与四轮定位调校"},
		{Category: "弹簧与阻尼器", UpgradeName: "拉力弹簧与阻尼器", Unlocks: "解锁弹簧、阻尼器与四轮定位调校"},
		{Category: "弹簧与阻尼器", UpgradeName: "漂移弹簧与阻尼器", Unlocks: "解锁弹簧、阻尼器与四轮定位调校"},
		{Category: "前防倾杆", UpgradeName: "赛车版前防倾杆", Unlocks: "解锁前防倾杆硬度调校"},
		{Category: "后防倾杆", UpgradeName: "赛车版后防倾杆", Unlocks: "解锁后防倾杆硬度调校"},
		{Category: "变速箱", UpgradeName: "跑车版变速箱", Unlocks: "解锁最终传动比调校"},
		{Category: "变速箱", UpgradeName: "赛车版变速箱", Unlocks: "解锁全传动比调校"},
		{Category: "差速器", UpgradeName: "跑车版差速器 1.5向", Unlocks: "解锁加速差速器调校"},
		{Category: "差速器", UpgradeName: "赛车版差速器 双向", Unlocks: "解锁全差速器调校"},
		{Category: "差速器", UpgradeName: "拉力赛差速器 双向", Unlocks: "解锁全差速器调校"},
		{Category: "差速器", UpgradeName: "漂移差速器 双向", Unlocks: "解锁全差速器调校"},
		{Category: "差速器", UpgradeName: "越野差速器 双向", Unlocks: "解锁全差速器调校"},
	}
	return rules
}
