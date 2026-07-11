const SHARE_FUNCTION_NAME = "shareTune";

function buildSharePath(shareId) {
  return `/pages/share-detail/index?id=${encodeURIComponent(shareId)}`;
}

async function createTuneShare(snapshot) {
  ensureCloudReady();
  const response = await callShareFunction({
    type: "create",
    payload: clonePlain(snapshot),
  });
  const body = response.result || {};
  if (!body.success) {
    throw normalizeShareError({ message: body.message || "分享链接生成失败" });
  }
  return body.data || {};
}

async function getTuneShare(shareId) {
  ensureCloudReady();
  const response = await callShareFunction({
    type: "get",
    shareId,
  });
  const body = response.result || {};
  if (!body.success) {
    throw normalizeShareError({ message: body.message || "分享内容读取失败" });
  }
  return body.data || {};
}

async function callShareFunction(data) {
  try {
    return await wx.cloud.callFunction({
      name: SHARE_FUNCTION_NAME,
      data,
    });
  } catch (err) {
    throw normalizeShareError(err);
  }
}

function normalizeShareError(err) {
  const message = err && err.errMsg ? err.errMsg : String(err && err.message ? err.message : err || "");
  if (message.indexOf("FUNCTION_NOT_FOUND") !== -1 || message.indexOf("-501000") !== -1) {
    return new Error("分享云函数 shareTune 未部署。请在微信开发者工具中上传并部署 cloudfunctions/shareTune。");
  }
  if (message.indexOf("COLLECTION_NOT_EXIST") !== -1 || message.indexOf("collection") !== -1) {
    return new Error("分享数据集合 tune_shares 未创建。请在云开发数据库中创建该集合后重试。");
  }
  return new Error(message || "分享服务暂不可用，请稍后再试。");
}

function shareTitle(summary) {
  if (summary && summary.displayName) {
    return `FH6 ${summary.displayName}`;
  }
  if (!summary || (!summary.carClass && !summary.pi)) {
    return "FH6车辆调校小助手";
  }
  const useCase = summary && summary.useCaseLabel ? summary.useCaseLabel : "车辆";
  const carClass = summary && summary.carClass ? summary.carClass : "";
  const pi = summary && summary.pi ? summary.pi : "";
  const drivetrain = summary && summary.drivetrain ? ` ${summary.drivetrain}` : "";
  return `FH6 ${useCase}调校 ${carClass}${pi}${drivetrain}`;
}

function ensureCloudReady() {
  const app = getApp();
  if (!app.globalData.env) {
    throw new Error("云环境未配置，无法生成或读取分享链接。");
  }
}

function clonePlain(value) {
  return JSON.parse(JSON.stringify(value || {}));
}

module.exports = {
  buildSharePath,
  createTuneShare,
  getTuneShare,
  shareTitle,
};
