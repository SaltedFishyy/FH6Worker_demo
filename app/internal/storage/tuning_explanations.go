package storage

var tuneAdjustmentExplanations = []TuneAdjustmentExplanation{
	{Category: "轮胎", Item: "胎压", Detail: "前侧", Description: "胎压调整的目的，在于让轮胎获得最大的地面接触面积，同时避免轮胎过热或丧失响应性。较低胎压可提升接触面积，但太低会降低响应速度并导致轮胎过热；较高胎压可提升响应速度，但会降低胎温，可能导致车辆瞬间失去抓地力。"},
	{Category: "轮胎", Item: "胎压", Detail: "后侧", Description: "胎压调整的目的，在于让轮胎获得最大的地面接触面积，同时避免轮胎过热或丧失响应性。较低胎压可提升接触面积，但太低会降低响应速度并导致轮胎过热；较高胎压可提升响应速度，但会降低胎温，可能导致车辆瞬间失去抓地力。"},
	{Category: "齿轮", Item: "前进档", Detail: "最终传动", Description: "最终传动比值越高，传动比越短，可提升加速但会降低极速，或导致最终齿轮成为瓶颈。降低最终传动比值可获得更高极速，但会降低加速性能。"},
	{Category: "齿轮", Item: "前进档", Detail: "1档", Description: "单独传动比会影响加速和极速。较高传动比可提升加速，较低传动比可获得更高极速，应根据赛道需要和引擎马力、扭矩特性设置。"},
	{Category: "齿轮", Item: "前进档", Detail: "2档", Description: "单独传动比会影响加速和极速。较高传动比可提升加速，较低传动比可获得更高极速，应根据赛道需要和引擎马力、扭矩特性设置。"},
	{Category: "齿轮", Item: "前进档", Detail: "3档", Description: "单独传动比会影响加速和极速。较高传动比可提升加速，较低传动比可获得更高极速，应根据赛道需要和引擎马力、扭矩特性设置。"},
	{Category: "齿轮", Item: "前进档", Detail: "4档", Description: "单独传动比会影响加速和极速。较高传动比可提升加速，较低传动比可获得更高极速，应根据赛道需要和引擎马力、扭矩特性设置。"},
	{Category: "齿轮", Item: "前进档", Detail: "5档", Description: "单独传动比会影响加速和极速。较高传动比可提升加速，较低传动比可获得更高极速，应根据赛道需要和引擎马力、扭矩特性设置。"},
	{Category: "齿轮", Item: "前进档", Detail: "6档", Description: "单独传动比会影响加速和极速。较高传动比可提升加速，较低传动比可获得更高极速，应根据赛道需要和引擎马力、扭矩特性设置。"},
	{Category: "齿轮", Item: "前进档", Detail: "7档", Description: "单独传动比会影响加速和极速。较高传动比可提升加速，较低传动比可获得更高极速，应根据赛道需要和引擎马力、扭矩特性设置。"},
	{Category: "齿轮", Item: "前进档", Detail: "8档", Description: "单独传动比会影响加速和极速。较高传动比可提升加速，较低传动比可获得更高极速，应根据赛道需要和引擎马力、扭矩特性设置。"},
	{Category: "轮胎定位", Item: "外倾角", Detail: "前侧", Description: "负值外倾角可提升过弯时轮胎与地面的接触面积，但会降低直线行驶时的接触面积；负值过大时，会降低直线加速和减速性能。正值外倾角会降低过弯接触面积并降低稳定性。"},
	{Category: "轮胎定位", Item: "外倾角", Detail: "后侧", Description: "负值外倾角可提升过弯时轮胎与地面的接触面积，但会降低直线行驶时的接触面积；负值过大时，会降低直线加速和减速性能。正值外倾角会降低过弯接触面积并降低稳定性。"},
	{Category: "轮胎定位", Item: "束角", Detail: "前侧", Description: "束角可改变转向响应。内束提高稳定性但降低转向反应速度；外束提高转向反应速度，但会降低稳定性。"},
	{Category: "轮胎定位", Item: "束角", Detail: "后侧", Description: "束角可改变转向响应。内束提高稳定性但降低转向反应速度；外束提高转向反应速度，但会降低稳定性。"},
	{Category: "轮胎定位", Item: "前轮后倾角", Detail: "角度", Description: "正值后倾角可提升直线稳定性，并在转弯时提供适当的负值外倾角，让车辆在直线加速和刹车时保持轮胎接触状态，同时保留弯中抓地。"},
	{Category: "防倾杆", Item: "防倾杆", Detail: "前侧", Description: "防倾杆用于控制车体运动，并在稳定过弯时平衡转向不足与转向过度。缓解转向不足可调软前防倾杆、调硬后防倾杆；缓解转向过度可调软后防倾杆、调硬前防倾杆。"},
	{Category: "防倾杆", Item: "防倾杆", Detail: "后侧", Description: "防倾杆用于控制车体运动，并在稳定过弯时平衡转向不足与转向过度。缓解转向不足可调软前防倾杆、调硬后防倾杆；缓解转向过度可调软后防倾杆、调硬前防倾杆。"},
	{Category: "弹簧", Item: "弹簧", Detail: "前侧", Description: "弹簧控制加速、刹车、过弯时的车重转移。前弹簧较软可增加前轮抓地并缓解转向不足，但太软可能重刹托底；前弹簧较硬可缓解转向过度，但太硬可能导致转向不足。"},
	{Category: "弹簧", Item: "弹簧", Detail: "后侧", Description: "弹簧控制加速、刹车、过弯时的车重转移。后弹簧较软可增加后轮抓地并缓解转向过度；后弹簧较硬会增加过度转向。"},
	{Category: "弹簧", Item: "车身高度", Detail: "前侧", Description: "车身高度决定离地间隙和重心。降低车身高度可降低重心并提升过弯性能，但太低可能托底并导致瞬间失控。应在避免托底的情况下尽可能降低车身高度。"},
	{Category: "弹簧", Item: "车身高度", Detail: "后侧", Description: "车身高度决定离地间隙和重心。降低车身高度可降低重心并提升过弯性能，但太低可能触底；适当提高后车身高度，也可在急加速时帮助控制重量转移。"},
	{Category: "阻尼", Item: "回弹硬度", Detail: "前侧", Description: "回弹阻尼控制悬挂伸展速率。增加前回弹会提高过渡性转向不足；降低前回弹会提高过渡性转向过度。"},
	{Category: "阻尼", Item: "回弹硬度", Detail: "后侧", Description: "回弹阻尼控制悬挂伸展速率。增加后回弹会提高过渡性转向不足；降低后回弹会提高过渡性转向过度。"},
	{Category: "阻尼", Item: "压缩硬度", Detail: "前侧", Description: "压缩阻尼控制悬挂压缩速率。增加前压缩会提高过渡性转向不足，过高则会让车辆在不平路面不稳定；降低前压缩会提高过渡性转向过度。"},
	{Category: "阻尼", Item: "压缩硬度", Detail: "后侧", Description: "压缩阻尼控制悬挂压缩速率。增加后压缩会提高过渡性转向不足，过高则会让车辆在不平路面不稳定；降低后压缩会提高过渡性转向过度。"},
	{Category: "空气动力学设置", Item: "下压力", Detail: "前侧", Description: "提升下压力可让车辆与路面保持更好接触、更快热胎，并改善高速操控性，但会增加风阻。空力平衡越低，车辆更偏转向不足；越高，则更偏转向过度。"},
	{Category: "空气动力学设置", Item: "下压力", Detail: "后侧", Description: "提升下压力可让车辆与路面保持更好接触、更快热胎，并改善高速操控性，但会增加风阻。空力平衡越低，车辆更偏转向不足；越高，则更偏转向过度。"},
	{Category: "刹车", Item: "制动力", Detail: "平衡", Description: "刹车平衡影响制动力分配、刹车距离和刹车时的转向平衡。向后调整会提升刹车时转向过度；向前调整会提升转向不足和稳定性，但也可能导致重刹转向不足。"},
	{Category: "刹车", Item: "制动力", Detail: "压力", Description: "刹车油压影响踏板行程与制动力关系。降低压力可提升产生显著制动力所需踏板行程，但太低会导致减速不足；提升压力可缩短制动力建立行程，但太高容易抱死。"},
	{Category: "差速器", Item: "后侧", Detail: "加速", Description: "加速差速器设置控制加速时差速器锁定速度。提高后轮加速差速会加剧后驱或四驱车辆的转向过度；降低加速设置会减慢锁定速度，改善部分动力甩尾。"},
	{Category: "差速器", Item: "后侧", Detail: "减速", Description: "减速差速器设置控制减速时差速器锁定速度。提高减速设置会提升锁定速度，但过快会影响操控；提高后轮减速设置可缓解松开油门时的转向过度。"},
}

func ListTuneAdjustmentExplanations() []TuneAdjustmentExplanation {
	out := make([]TuneAdjustmentExplanation, len(tuneAdjustmentExplanations))
	copy(out, tuneAdjustmentExplanations)
	return out
}

func TuneAdjustmentExplanationsForAction(actionItem string, gear int) []TuneAdjustmentExplanation {
	keys := explanationKeysForAction(actionItem, gear)
	out := make([]TuneAdjustmentExplanation, 0, len(keys))
	for _, key := range keys {
		for _, item := range tuneAdjustmentExplanations {
			if explanationKey(item.Category, item.Item, item.Detail) == key {
				out = append(out, item)
				break
			}
		}
	}
	return out
}

func explanationKeysForAction(actionItem string, gear int) []string {
	switch actionItem {
	case "gear_1":
		return []string{explanationKey("齿轮", "前进档", "1档")}
	case "current_gear":
		if gear >= 1 && gear <= 10 {
			return []string{explanationKey("齿轮", "前进档", string(rune('0'+gear))+"档")}
		}
	case "final_drive":
		return []string{explanationKey("齿轮", "前进档", "最终传动")}
	case "brake_balance":
		return []string{explanationKey("刹车", "制动力", "平衡")}
	case "brake_pressure":
		return []string{explanationKey("刹车", "制动力", "压力")}
	case "rear_diff_accel", "drive_diff_accel":
		return []string{explanationKey("差速器", "后侧", "加速")}
	case "rear_diff_decel":
		return []string{explanationKey("差速器", "后侧", "减速")}
	case "drive_tire_pressure", "tire_pressure":
		return []string{explanationKey("轮胎", "胎压", "前侧"), explanationKey("轮胎", "胎压", "后侧")}
	case "front_arb":
		return []string{explanationKey("防倾杆", "防倾杆", "前侧")}
	case "rear_arb":
		return []string{explanationKey("防倾杆", "防倾杆", "后侧")}
	case "front_rebound":
		return []string{explanationKey("阻尼", "回弹硬度", "前侧")}
	case "rear_rebound":
		return []string{explanationKey("阻尼", "回弹硬度", "后侧")}
	case "front_camber":
		return []string{explanationKey("轮胎定位", "外倾角", "前侧")}
	case "front_and_rear_aero":
		return []string{explanationKey("空气动力学设置", "下压力", "前侧"), explanationKey("空气动力学设置", "下压力", "后侧")}
	case "ride_height":
		return []string{explanationKey("弹簧", "车身高度", "前侧"), explanationKey("弹簧", "车身高度", "后侧")}
	case "spring_rate":
		return []string{explanationKey("弹簧", "弹簧", "前侧"), explanationKey("弹簧", "弹簧", "后侧")}
	case "bump":
		return []string{explanationKey("阻尼", "压缩硬度", "前侧"), explanationKey("阻尼", "压缩硬度", "后侧")}
	}
	return nil
}

func explanationKey(category string, item string, detail string) string {
	return category + "/" + item + "/" + detail
}
