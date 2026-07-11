# FH6 不同用途与不同等级升级配件选择优先级方案

版本：v1.0  
适用项目：FH6 Data Out + AI Agent 调校分析助手 / 调校基线生成器  
文档用途：用于在程序中根据车辆用途、PI 等级、驱动形式和车辆短板，生成升级配件推荐优先级。

---

## 1. 文档目标

本文件用于整理《极限竞速：地平线 6》（FH6）中不同用途、不同等级车辆的升级配件选择优先级。

核心目标不是生成“唯一最强改装”，而是为程序提供一套可执行的升级推荐规则：

```text
输入：
- 车辆 PI 等级
- 用途：公路、街头、Time Attack、直线、拉力、越野、漂移
- 驱动形式：FWD / RWD / AWD
- 当前车辆短板：推头、甩尾、起步打滑、刹车弱、跳跃不稳等

输出：
- P0 必选升级件
- P1 高优先级升级件
- P2 视情况升级件
- P3 最后补充升级件
- 不推荐升级件
```

该文档可作为：

```text
1. 调校基线生成器的升级件推荐规则；
2. AI Agent 生成改装建议的知识库；
3. Data Out 遥测分析后的下一步升级建议来源；
4. 前端“升级建议”页面的数据基础。
```

---

## 2. 优先级标记

| 标记 | 含义 | 程序行为 |
|---|---|---|
| P0 必选 | 该用途的核心升级，通常优先安装 | 默认推荐 |
| P1 高优先级 | 大多数车辆适用，PI 效率较高 | 默认推荐，但可按 PI 预算筛选 |
| P2 视情况 | 取决于车辆、路线、驱动和当前短板 | 需要理由或条件触发 |
| P3 最后补 | PI 有余量或车辆短板明显时才安装 | 低优先级 |
| 不推荐 | 通常浪费 PI 或破坏用途特性 | 默认不推荐，除非用户手动选择 |

---

## 3. FH6 PI 等级口径

FH6 的等级口径建议按以下方式处理：

| 等级 | PI 范围 | 升级策略 |
|---|---:|---|
| D | 100–400 | 少改，保持轻量，不建议大件转换 |
| C | 401–500 | 基础轻量、基础动力、差速器优先 |
| B | 501–600 | 低预算构筑，优先低 PI 高收益件 |
| A | 601–700 | 均衡等级，适合完整基础操控件 + 适度动力 |
| S1 | 701–800 | 轮胎、差速、刹车、空力开始重要 |
| S2 | 801–900 | 高速稳定、强抓地、刹车、空力优先 |
| R | 901–998 | 赛道 / 极限构筑，完整底盘、空力、热熔胎优先 |
| X | 999 | 娱乐 / 极速 / 非均衡构筑，不作为常规推荐基线 |

---

## 4. 通用升级原则

### 4.1 先用途，后配件

升级件选择必须先明确用途：

| 用途 | 核心目标 |
|---|---|
| 公路 / 街头 | 抓地、刹车、弯中稳定、出弯牵引 |
| Time Attack | 极限抓地、刹车稳定、高速空力、圈速 |
| 直线 / 极速 | 动力重量比、换挡效率、起步抓地、低风阻 |
| 拉力 / Dirt | AWD、拉力胎、拉力悬挂、吸震、可控滑动 |
| 越野 / Cross Country | 越野胎、高车身、长行程、落地稳定 |
| 漂移 | RWD、漂移悬挂、差速锁止、齿比、动力响应 |

### 4.2 不同用途的“不要优先做”

| 用途 | 不要优先做 |
|---|---|
| 公路 / 街头 | 不要盲目堆马力，不要忽略轮胎和刹车 |
| Time Attack | 不要为了极速牺牲空力和刹车 |
| 直线 / 极速 | 不要装过多高下压力空力件 |
| 拉力 | 不要装 Race Slick / Race Suspension |
| 越野 | 不要低车高、硬悬挂、光头胎 |
| 漂移 | 不要装过强抓地胎导致漂不起来 |

### 4.3 PI 预算原则

升级件受 PI 上限约束，因此应遵循：

```text
低等级：优先低 PI 高收益件
中等级：先底盘抓地，再补动力
高等级：完整底盘 + 轮胎 + 刹车 + 空力，再补动力
```

---

## 5. 公路 / 街头 / Time Attack 升级优先级

### 5.1 用途目标

```text
目标：
- 弯中抓地
- 入弯刹车稳定
- 出弯牵引
- 高速弯稳定
- 可控推头/甩尾平衡
```

---

### 5.2 D / C / B 级公路

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Race Differential | 解锁差速器调校，低 PI 高收益 |
| P0 | Anti-Roll Bars | 解锁防倾杆调校，修正推头/甩尾 |
| P1 | Street / Sport Tire | 视 PI 预算选择，不一定直接上最高胎 |
| P1 | Tire Width | 驱动轮优先，RWD 优先后胎 |
| P1 | Weight Reduction | 降低车重，提升刹车、转向、加速 |
| P1 | Exhaust / Intake | 轻量加动力，PI 效率较好 |
| P2 | Brakes | 重车、街头赛、刹车点多时装 |
| P2 | Sport / Race Suspension | 需要调弹簧和车高时装 |
| P3 | Transmission | 原厂齿比差或需要终传比时装 |
| 不推荐 | Engine Swap | 低等级 PI 空间小，容易浪费预算 |
| 不推荐 | AWD Swap | 除非原车严重打滑，否则低等级不优先 |

---

### 5.3 A700 公路

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Sport / Semi-Slick Tire | A700 抓地核心，按车重和路线选择 |
| P0 | Race Differential | 必装，影响出弯和入弯稳定 |
| P0 | Anti-Roll Bars | 必装，中弯平衡核心 |
| P1 | Race / Sport Suspension | 建议解锁弹簧、车高、阻尼 |
| P1 | Tire Width | 前后胎宽独立优化 |
| P1 | Weight Reduction | 通常比盲目加动力更稳 |
| P1 | Brakes | 街头、山路、重车优先 |
| P2 | Transmission | 齿比不合理、起步或高速异常时装 |
| P2 | Clutch / Driveline | 低成本改善换挡和响应 |
| P2 | Aero | 高速弯多、Time Attack 可装；低速街头可不装 |
| P2 | Mild Power Upgrades | 排气、进气、点火、燃油、活塞 |
| P3 | Turbo / Supercharger | 只在动力明显不足时装 |
| P3 | Engine Swap | 原厂发动机 PI 效率差时才考虑 |

---

### 5.4 S1 公路

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Semi-Slick / Race Slick | 速度提升后，抓地优先 |
| P0 | Race Differential | RWD / AWD 必装 |
| P0 | Race Suspension | 支撑高速弯和压弯 |
| P0 | Anti-Roll Bars | 中高速平衡核心 |
| P0 | Brakes | S1 车速高，刹车稳定性重要 |
| P1 | Aero | 高速弯、山路、Time Attack 建议装 |
| P1 | Tire Width | 前后都应重点处理 |
| P1 | Weight Reduction | 刹车和变向收益明显 |
| P1 | Transmission | 高马力车需要调齿比 |
| P2 | AWD Conversion | 高马力抓不住时考虑，但可能推头 |
| P2 | Forced Induction | 需要补中高速动力时装 |
| P3 | Engine Swap | 看动力重量比，不要破坏重量分布 |

---

### 5.5 S2 / R 公路与 Time Attack

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Race Slick | 不要先堆马力 |
| P0 | Aero 前后空力 | 高速稳定和弯中速度核心 |
| P0 | Race Brakes | 高速刹车必须稳定 |
| P0 | Race Suspension / ARB / Diff | 完整底盘调校必需 |
| P1 | Tire Width | 提高机械抓地 |
| P1 | Wheel Spacers | 提高稳定性，视车辆可用性 |
| P1 | Weight Reduction | 高速刹车和变向收益大 |
| P1 | Transmission / Clutch / Driveline | 高速齿比必须重设 |
| P2 | AWD Conversion | S2 高马力可考虑，但要处理推头 |
| P2 | Turbo / Supercharger | 抓地和刹车完成后补动力 |
| P2 | Engine Swap | 只在原厂动力平台不足时考虑 |
| P3 | Roll Cage | 老车、高速不稳或车身过软时考虑 |

---

## 6. 直线 / 极速 / Drag 升级优先级

### 6.1 用途目标

```text
目标：
- 起步抓地
- 低挡不空转
- 换挡效率
- 中高速持续加速
- 最高速
- 低风阻
```

---

### 6.2 D / C / B 直线

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Exhaust / Intake | 低 PI 加动力 |
| P0 | Weight Reduction | 提升全段加速 |
| P1 | Clutch / Driveline | 低成本提升传动效率 |
| P1 | Transmission | 需要调终传比时装 |
| P1 | Tire Width | 驱动轮优先 |
| P2 | Tire Compound | 不一定上最高胎，防止吃光 PI |
| P2 | Ignition / Fuel / Pistons | 用来补 PI |
| P3 | Turbo / Supercharger | 低等级谨慎，容易打滑 |
| 不推荐 | Aero | 直线低风阻优先 |
| 不推荐 | Heavy Roll Cage | 增重可能抵消收益 |

---

### 6.3 A700 直线

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Transmission | 1挡 / 终传比 / 最高挡必须可调 |
| P0 | Clutch / Driveline | 减少传动损失 |
| P0 | Exhaust / Intake | 高 PI 效率动力件 |
| P1 | Ignition / Fuel / Pistons | 用于补动力 |
| P1 | Tire Width | 起步抓地优先 |
| P1 | Weight Reduction | 提升全段加速 |
| P1 | Forced Induction | 直线常用，但要看打滑 |
| P2 | AWD Conversion | 起步强，但吃 PI |
| P2 | Drag Tire / High-Grip Tire | 起步抓不住时装 |
| P3 | Engine Swap | 看动力重量比，不只看最大马力 |
| 不推荐 | 高下压力 Aero | 会拖慢极速 |

---

### 6.4 S1 / S2 / R 直线与极速

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | High Power Engine / Forced Induction | 高速路和直线核心 |
| P0 | Race Transmission | 必须精调每挡 |
| P0 | Clutch / Driveline | 高马力降低传动损失 |
| P0 | Tire Width / Drag Tire | 起步和低挡牵引 |
| P1 | AWD Conversion | 高马力直线常见，但吃 PI 和重量 |
| P1 | Weight Reduction | 0–300 提速明显 |
| P2 | Brakes | 直线不一定，但测速后需要稳定 |
| P2 | Roll Cage | 高速发飘时考虑 |
| P3 | Aero | 除非高速不稳，否则尽量低阻 |

---

## 7. 拉力 / Dirt 升级优先级

### 7.1 用途目标

```text
目标：
- 土路抓地
- 吸收颠簸
- 跳跃落地稳定
- 可控滑动
- 出弯牵引
```

---

### 7.2 D / C / B 拉力

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Rally Tires | 路面适配优先 |
| P0 | Rally Suspension | 吸震和车高核心 |
| P0 | Race / Rally Differential | 差速器调校核心 |
| P1 | AWD Conversion | 拉力大多数优先 AWD |
| P1 | Weight Reduction | 降低惯性和跳跃负担 |
| P1 | Clutch / Driveline | 响应和传动效率 |
| P2 | Anti-Roll Bars | 可调即可，不要过硬 |
| P2 | Exhaust / Intake | 温和补动力 |
| P3 | Transmission | 原厂齿比不适合土路时装 |
| 不推荐 | Race Slick | 土路不适合 |
| 不推荐 | Race Suspension | 土路吸震差 |

---

### 7.3 A700 拉力

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | AWD Conversion | 拉力核心 |
| P0 | Rally Tires | 土路抓地核心 |
| P0 | Rally Suspension | 吸震和车高 |
| P0 | Race Differential | 前后 / 中央差速必须可调 |
| P1 | Weight Reduction | 对跳跃和转向帮助大 |
| P1 | Transmission | 土路常用 2–5 挡，需要可调 |
| P1 | Clutch / Driveline | 换挡和响应 |
| P2 | Tire Width | 前后都加，但不要牺牲过多 PI |
| P2 | Mild Power | 排气、进气、点火、燃油、活塞 |
| P2 | Turbo / Supercharger | 需要中速加速时装 |
| P3 | Brakes | 下坡 / 高速混合路段再加强 |
| 不推荐 | Race Suspension / Slick Tire | 会破坏土路表现 |

---

### 7.4 S1 / S2 拉力

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Rally Tires | 高马力土路仍然依赖合适轮胎 |
| P0 | Rally Suspension | 吸震、跳跃、落地稳定 |
| P0 | AWD / Race Diff | 高马力土路必需 |
| P0 | Transmission | 齿比决定土路牵引 |
| P1 | Weight Reduction | 控制高马力车姿态 |
| P1 | Tire Width | 抓地和稳定性 |
| P1 | Brakes | 高速土路需要 |
| P2 | Forced Induction / Engine Swap | 只在能踩住油时补 |
| P2 | Aero | 高速 Dirt 路线可少量使用 |
| 不推荐 | 过度硬底盘 | 过坑弹飞 |

---

## 8. 越野 / Cross Country 升级优先级

### 8.1 用途目标

```text
目标：
- 草地、沙地、河床通过性
- 高车身
- 长行程
- 跳跃落地稳定
- 低中速扭矩
```

---

### 8.2 D / C / B 越野

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Off-road Tires | 越野核心 |
| P0 | Off-road Suspension | 高车身、长行程 |
| P0 | AWD Conversion | 越野基本必需 |
| P1 | Differential | 前后 / 中央差速可调 |
| P1 | Weight Reduction | 减少落地冲击 |
| P1 | Clutch / Driveline | 响应和传动效率 |
| P2 | Tire Width | 越野稳定性 |
| P2 | Torque Power Upgrades | 低中速扭矩比极速更重要 |
| P3 | Transmission | 爬坡 / 出弯需要短齿比时装 |
| 不推荐 | Race Suspension | 过坑不稳 |
| 不推荐 | Race Slick | 路面错误 |

---

### 8.3 A700 越野

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Off-road Tires | 必选 |
| P0 | Off-road Suspension | 必选 |
| P0 | AWD | 牵引和通过性核心 |
| P0 | Race / Off-road Differential | 差速可调 |
| P1 | Tire Width | 大车、皮卡、SUV 优先 |
| P1 | Weight Reduction | 落地和变向提升明显 |
| P1 | Transmission | 越野齿比偏短 |
| P2 | Brakes | 下坡、跳跃后入弯需要 |
| P2 | Turbo / Supercharger | 提升低中速扭矩 |
| P3 | Engine Swap | 看动力重量比和车头重量 |
| 不推荐 | Aero | 多数越野路线收益低 |

---

### 8.4 S1 / S2 越野

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Off-road Tire | 不可省 |
| P0 | Off-road Suspension | 高车身和长行程 |
| P0 | AWD / Diff / Transmission | 高马力越野必须 |
| P1 | Weight Reduction | 控制跳跃和落地 |
| P1 | Brakes | 高速越野需要 |
| P1 | Torque Engine / Forced Induction | 中低速再加速 |
| P2 | Tire Width / Track Width | 抗翻滚、稳定 |
| P3 | Aero | 只有高速越野路线可少量考虑 |
| 不推荐 | 低车高 / 硬悬挂 | 过坑和落地容易失控 |

---

## 9. 漂移升级优先级

### 9.1 用途目标

```text
目标：
- 后轮可控滑动
- 足够转向角
- 油门响应快
- 常用挡位齿比合适
- 差速器高锁止
```

---

### 9.2 D / C / B 漂移

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | RWD 保持 / RWD Conversion | 漂移核心 |
| P0 | Drift Suspension | 转向角和姿态 |
| P0 | Race Differential | 后差速高锁止 |
| P1 | Transmission | 常用 3 / 4 挡必须可调 |
| P1 | Clutch / Driveline | 响应和换挡 |
| P1 | Exhaust / Intake | 温和补动力 |
| P2 | Drift Tire / Sport Tire | 看车是否太滑或太黏 |
| P2 | Flywheel | 提高转速响应 |
| P3 | Weight Reduction | 适度，太轻可能难控 |
| 不推荐 | Race Slick | 漂不起来 |

---

### 9.3 A700 漂移

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Drift Suspension | 转向角核心 |
| P0 | Race Diff | 后差速高锁止 |
| P0 | Transmission | 齿比是漂移手感核心 |
| P1 | Drift Tire / Sport Tire | 新手可 Drift，高手按车感选 |
| P1 | Power Upgrades | 500–700hp 区间好控 |
| P1 | Clutch / Driveline / Flywheel | 提高响应 |
| P2 | Weight Reduction | 适度即可 |
| P2 | ARB | 调整横摆响应 |
| P3 | Engine Swap / Forced Induction | 需要更大角度或高速漂移时 |
| 不推荐 | AWD 公路化抓地构筑 | 失去漂移特性 |

---

### 9.4 S1 / S2 漂移

| 优先级 | 升级件 | 说明 |
|---|---|---|
| P0 | Drift Suspension | 必选 |
| P0 | Race Diff | 必选 |
| P0 | Race Transmission | 必选 |
| P0 | High Power Engine / Turbo / Supercharger | 高速漂移和刷分 |
| P1 | Drift Tire / Sport Tire | 视抓地和滑动控制 |
| P1 | Clutch / Driveline / Flywheel | 保持高转响应 |
| P1 | ARB | 控制角度和切换 |
| P2 | Weight Reduction | 视车重和稳定性 |
| P2 | Aero | 漂移一般不优先，除非高速稳定需要 |
| P3 | AWD Drift | 特定玩法，不做默认推荐 |

---

## 10. 通用升级件评分矩阵

该矩阵可用于程序内部 `upgradeScore` 计算。

分值说明：

```text
5 = 强推荐
4 = 高优先级
3 = 视情况
2 = 低优先级
1 = 少数情况
0 = 中性
负数 = 通常不推荐
```

| 升级件类别 | 直线 / Speed | 公路 / Grip | 漂移 | 拉力 / Dirt | 越野 |
|---|---:|---:|---:|---:|---:|
| Intake / Exhaust | 5 | 4 | 4 | 4 | 4 |
| Ignition / Fuel / Pistons | 5 | 4 | 4 | 4 | 4 |
| Camshaft | 3 | 2 | 2 | 1 | 1 |
| Turbo / Supercharger | 5 | 3 | 4 | 3 | 4 |
| Intercooler / Oil Cooler | 4 | 2 | 2 | 1 | 1 |
| Flywheel | 3 | 2 | 5 | 3 | 2 |
| Brakes | 2 | 5 | 2 | 3 | 3 |
| Race Suspension | 2 | 5 | 0 | -5 | -5 |
| Rally Suspension | -5 | -5 | -5 | 5 | 3 |
| Off-road Suspension | -5 | -5 | -5 | 3 | 5 |
| Anti-Roll Bars | 3 | 5 | 5 | 3 | 2 |
| Roll Cage | 2 | 3 | 1 | 2 | 2 |
| Weight Reduction | 5 | 5 | 4 | 5 | 5 |
| Transmission | 5 | 3 | 5 | 3 | 3 |
| Clutch / Driveline | 5 | 4 | 5 | 4 | 4 |
| Differential | 4 | 5 | 5 | 5 | 5 |
| Stock-to-Slick Tire | 4 | 5 | -5 | -5 | -5 |
| Drift Tire | 0 | -3 | 5 | -5 | -5 |
| Rally Tire | -3 | -3 | -5 | 5 | 3 |
| Off-road Tire | -5 | -5 | -5 | 3 | 5 |
| Tire Width | 4 | 5 | 2 | 4 | 4 |
| Rim Weight Optimization | 2 | 3 | 2 | 2 | 2 |
| Wheel Spacers | 2 | 3 | 2 | 2 | 3 |
| Aero | 1 | 4 | 1 | 1 | 0 |
| AWD Conversion | 3 | 3 | -4 | 5 | 5 |
| Forced Induction Conversion | 5 | 3 | 4 | 4 | 4 |
| Engine Swap | 5 | 4 | 3 | 3 | 3 |

---

## 11. 程序化推荐规则

### 11.1 输入结构建议

```ts
export type UpgradeUseCase =
  | "road"
  | "street"
  | "timeAttack"
  | "drag"
  | "speed"
  | "rally"
  | "offroad"
  | "drift";

export type UpgradeClass =
  | "D"
  | "C"
  | "B"
  | "A"
  | "S1"
  | "S2"
  | "R"
  | "X";

export type UpgradeRecommendationInput = {
  carName?: string;
  pi: number;
  className: UpgradeClass;
  drivetrain: "FWD" | "RWD" | "AWD";
  useCase: UpgradeUseCase;

  currentProblems?: Array<
    | "launch_wheelspin"
    | "launch_bog_down"
    | "understeer"
    | "oversteer"
    | "brake_lockup"
    | "weak_braking"
    | "high_speed_instability"
    | "suspension_bottom_out"
    | "bump_instability"
  >;

  piBudgetRemaining?: number;
};
```

### 11.2 输出结构建议

```ts
export type UpgradeRecommendationOutput = {
  useCase: UpgradeUseCase;
  className: UpgradeClass;
  summary: string;

  p0: UpgradeItemRecommendation[];
  p1: UpgradeItemRecommendation[];
  p2: UpgradeItemRecommendation[];
  p3: UpgradeItemRecommendation[];
  avoid: UpgradeItemRecommendation[];
};

export type UpgradeItemRecommendation = {
  name: string;
  category: string;
  reason: string;
  condition?: string;
  risk?: string;
};
```

---

## 12. AI Agent 升级建议规则

### 12.1 低等级 D / C / B

```text
当 PI 等级为 D / C / B：
1. 优先选择低 PI 高收益部件；
2. 不默认推荐发动机互换；
3. 不默认推荐 AWD 转换，除非用途为 Rally / Offroad；
4. 不默认推荐 Race Slick；
5. 优先选择差速器、防倾杆、轻量化、排气、进气、驱动轮胎宽；
6. 如果车辆仍然打滑，再考虑更高轮胎等级；
7. 如果车辆仍然动力不足，再补发动机小件。
```

### 12.2 A700

```text
当 PI 等级为 A：
1. 公路：Sport/Semi 胎、差速器、防倾杆、悬挂、减重优先；
2. 拉力：AWD、Rally 胎、Rally 悬挂、差速器优先；
3. 越野：AWD、Off-road 胎、Off-road 悬挂优先；
4. 直线：传动、动力、驱动轮抓地优先；
5. 漂移：RWD、Drift 悬挂、Race Diff、Transmission 优先。
```

### 12.3 S1 / S2 / R

```text
当 PI 等级为 S1 / S2 / R：
1. 先保证轮胎、刹车、差速、悬挂和高速稳定；
2. 公路 / Time Attack 必须考虑空力；
3. 高马力 RWD 若抓不住，可建议 AWD 转换；
4. 动力升级应排在抓地和刹车之后；
5. 如果车辆已经推头，不应继续盲目加 AWD 或前胎压力；
6. 如果高速侧滑，优先空力、轮胎和悬挂，再补动力。
```

---

## 13. Data Out 闭环触发升级建议

调校软件不应只根据静态用途推荐升级件，还应结合 Data Out 检测结果动态提高某些升级项优先级。

| Data Out 检测问题 | 提高优先级 |
|---|---|
| 起步烧胎 | 驱动轮胎宽、轮胎等级、差速器、齿比、AWD 转换 |
| 起步憋转 | 齿比、传动系统、动力小件 |
| 刹车弱 | Race Brakes、轮胎、减重 |
| 前轮抱死 | 轮胎、刹车平衡调校、前胎宽 |
| 入弯推头 | 前胎宽、前轮抓地、前 ARB、悬挂 |
| 出弯甩尾 | 后胎宽、后差速、齿比、后轮抓地 |
| 高速侧滑 | Aero、Race Tire、Race Suspension、Brakes |
| 过坑弹飞 | Rally / Off-road Suspension、车高、阻尼 |
| 落地触底 | Off-road Suspension、车高、弹簧、减重 |
| 胎温异常 | 胎压、轮胎等级、外倾角 |

---

## 14. 升级件冲突规则

### 14.1 不推荐组合

| 组合 | 问题 |
|---|---|
| Race Suspension + Rally/Dirt | 吸震差，过坑弹飞 |
| Race Suspension + Cross Country | 车高和行程不足 |
| Race Slick + Off-road | 路面不匹配 |
| Race Slick + Rally | Dirt 路面抓地不足 |
| 高下压力 Aero + 极速 | 风阻拖慢 |
| Engine Swap + 低等级 PI | PI 空间被吃光，底盘没预算 |
| AWD Swap + 低等级公路车 | PI 成本高，可能变慢 |
| 最大马力 + 原厂胎 | 打滑，车速不上 |
| Heavy Roll Cage + 低动力轻车 | 增重影响加速 |

### 14.2 冲突解决优先级

```text
1. 先满足赛事用途
2. 再满足 PI 上限
3. 再满足轮胎和路面
4. 再满足驱动形式
5. 再根据 Data Out 短板修正
6. 最后补动力
```

---

## 15. 示例输出

### 15.1 A700 AWD 公路车

```json
{
  "useCase": "road",
  "class": "A",
  "pi": 700,
  "drivetrain": "AWD",
  "upgradePlan": {
    "P0": [
      "Race Differential",
      "Anti-Roll Bars",
      "Sport / Semi-Slick Tire",
      "Tire Width"
    ],
    "P1": [
      "Race Suspension",
      "Weight Reduction",
      "Brakes",
      "Clutch / Driveline"
    ],
    "P2": [
      "Transmission",
      "Aero",
      "Exhaust",
      "Intake",
      "Ignition / Fuel / Pistons"
    ],
    "P3": [
      "Turbo / Supercharger",
      "Engine Swap"
    ],
    "avoid": [
      "Off-road Suspension",
      "Rally Tire",
      "Unnecessary AWD Swap"
    ]
  }
}
```

### 15.2 A700 AWD 拉力车

```json
{
  "useCase": "rally",
  "class": "A",
  "pi": 700,
  "drivetrain": "AWD",
  "upgradePlan": {
    "P0": [
      "AWD Conversion",
      "Rally Tires",
      "Rally Suspension",
      "Race Differential"
    ],
    "P1": [
      "Weight Reduction",
      "Transmission",
      "Clutch / Driveline"
    ],
    "P2": [
      "Tire Width",
      "Mild Power Upgrades",
      "Turbo / Supercharger"
    ],
    "P3": [
      "Brakes",
      "Engine Swap"
    ],
    "avoid": [
      "Race Slick",
      "Race Suspension",
      "High Downforce Aero"
    ]
  }
}
```

### 15.3 S1 RWD 漂移车

```json
{
  "useCase": "drift",
  "class": "S1",
  "pi": 800,
  "drivetrain": "RWD",
  "upgradePlan": {
    "P0": [
      "Drift Suspension",
      "Race Differential",
      "Race Transmission"
    ],
    "P1": [
      "High Power Engine / Forced Induction",
      "Clutch / Driveline / Flywheel",
      "Drift Tire / Sport Tire"
    ],
    "P2": [
      "Anti-Roll Bars",
      "Weight Reduction"
    ],
    "P3": [
      "Aero",
      "Engine Swap"
    ],
    "avoid": [
      "Race Slick",
      "Off-road Tire",
      "Rally Suspension"
    ]
  }
}
```

---

## 16. 最终总结

按照 ForzaTune 的升级矩阵思路，升级配件选择不能只看“数值更大”，而要看：

```text
用途
PI 成本
路面
驱动形式
车辆短板
```

最简单的记法：

```text
公路 / 街头：
轮胎、胎宽、差速、防倾杆、悬挂、刹车、减重 > 动力

直线 / 极速：
动力、传动、减重、驱动轮抓地 > 操控件 > 空力

拉力：
AWD、Rally 胎、Rally 悬挂、差速、传动、减重 > 动力

越野：
AWD、Off-road 胎、Off-road 悬挂、高车身、差速、扭矩 > 极速

漂移：
RWD、Drift 悬挂、Race Diff、Transmission、动力响应 > 抓地胎
```

对于程序实现，建议：

```text
1. 先用 useCase + class + drivetrain 生成升级优先级；
2. 再用 PI 预算筛选；
3. 再用 Data Out 识别的车辆问题修正推荐；
4. 最后由 AI Agent 用自然语言解释为什么推荐这些升级件。
```
