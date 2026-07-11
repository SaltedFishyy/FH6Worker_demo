const { decorateResultGroups } = require("../../utils/tuneDescriptions");
const { getTuneShare, shareTitle } = require("../../utils/tuneShare");

Page({
  data: {
    shareId: "",
    loading: true,
    errorMessage: "",
    summary: {},
    badge: {
      carClass: "",
      pi: "",
      drivetrain: "",
    },
    resultGroups: [],
    warnings: [],
    nextTestPlan: [],
  },

  onLoad(options) {
    this.loadShare(options && options.id ? options.id : "");
  },

  async loadShare(shareId) {
    const id = String(shareId || "").trim();
    if (!id) {
      this.setData({
        loading: false,
        errorMessage: "分享链接缺少参数。",
      });
      return;
    }

    this.setData({
      shareId: id,
      loading: true,
      errorMessage: "",
    });
    try {
      const data = await getTuneShare(id);
      const summary = data.summary || {};
      const profile = data.result && data.result.profileDraft ? data.result.profileDraft : {};
      this.setData({
        loading: false,
        summary,
        badge: {
          carClass: profile.carClass || summary.carClass || "",
          pi: profile.pi || summary.pi || "",
          drivetrain: profile.drivetrain || summary.drivetrain || "",
        },
        resultGroups: decorateResultGroups(data.resultGroups),
        warnings: Array.isArray(data.warnings) ? data.warnings : [],
        nextTestPlan: Array.isArray(data.nextTestPlan) ? data.nextTestPlan : [],
      });
    } catch (err) {
      this.setData({
        loading: false,
        errorMessage: err && err.message ? err.message : "分享内容读取失败。",
      });
    }
  },

  openQuickTune() {
    wx.switchTab({
      url: "/pages/index/index",
    });
  },

  showTuneDescription(e) {
    const title = e.currentTarget.dataset.title || "调校说明";
    const description = e.currentTarget.dataset.description || "";
    if (!description) return;
    wx.showModal({
      title,
      content: description,
      showCancel: false,
      confirmText: "知道了",
    });
  },

  onShareAppMessage() {
    return {
      title: shareTitle(this.data.summary),
      path: this.data.shareId ? `/pages/share-detail/index?id=${this.data.shareId}` : "/pages/index/index",
    };
  },
});
