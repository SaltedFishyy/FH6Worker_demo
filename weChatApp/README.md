# FH6车辆调校小助手

这是微信小程序车辆调校 MVP。前端负责参数输入、广告软限制和结果展示；后端使用 CloudBase 云函数 `calculateTune` 计算 Road / Drift 静态调校基线。

## 配置

在 `miniprogram/config.js` 中填写：

- `cloudEnvId`：CloudBase 环境 ID。
- `bannerAdUnitId`：顶部 Banner 广告位 ID，留空则隐藏。
- `rewardedAdUnitId`：激励视频广告位 ID，留空则开发模式不拦截调校。

## 部署云函数

在微信开发者工具中上传并部署 `cloudfunctions/calculateTune`，或使用 `uploadCloudFunction.sh` 指向当前项目和环境部署。

## 数据策略

小程序和云函数不保存车辆输入与调校结果。广告次数只存储在小程序本地缓存中，属于可绕过的软限制。
