const runtimeConfig = require("./config");

App({
  globalData: {
    config: runtimeConfig,
    env: runtimeConfig.cloudEnvId,
  },

  onLaunch() {
    if (!wx.cloud) {
      console.error("请使用 2.2.3 或以上的基础库以使用云能力");
      return;
    }

    if (!runtimeConfig.cloudEnvId) {
      console.warn("CloudBase 环境 ID 未配置，车辆调校云函数调用会被阻止。");
      return;
    }

    wx.cloud.init({
      env: runtimeConfig.cloudEnvId,
      traceUser: true,
    });
  },
});
