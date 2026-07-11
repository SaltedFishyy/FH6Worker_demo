const runtimeConfig = require("../../config");
const { setSelectedTab, setTabBarHidden } = require("../../utils/tabBar");
const { addTuneRecord, updateTuneRecord } = require("../../utils/tuneHistory");
const {
  decorateResultGroups,
  findResultField,
  parseManualDisplayValue,
  updateGroupsFieldValue,
  updateResultFieldValue,
} = require("../../utils/tuneDescriptions");
const { buildSharePath, createTuneShare, shareTitle } = require("../../utils/tuneShare");

const storageKey = "quickTuneAdStateV1";
const resultEditTipStorageKey = "quickTuneResultEditTipSeenV1";

const useCaseOptions = [
  { label: "公路", value: "Road" },
  { label: "拉力", value: "Rally" },
  { label: "越野", value: "Offroad" },
  { label: "直线", value: "Drag" },
  { label: "漂移", value: "Drift" },
];

const drivetrainOptions = [
  { label: "四驱", value: "AWD" },
  { label: "后驱", value: "RWD" },
  { label: "前驱", value: "FWD" },
];

const tireCompoundOptions = [
  { label: "标准", value: "stock" },
  { label: "街车", value: "street" },
  { label: "跑车", value: "sport" },
  { label: "半热熔", value: "semi" },
  { label: "热熔胎", value: "slick" },
  { label: "拉力", value: "rally" },
  { label: "越野", value: "offroad" },
  { label: "漂移", value: "drift" },
  { label: "直线", value: "drag" },
  { label: "雪地", value: "snow" },
];

const defaultForm = {
  useCase: "Road",
  pi: "",
  drivetrain: "AWD",
  tireCompound: "sport",
  weightKG: "",
  frontWeightPct: "",
  powerKW: "",
  torqueNM: "",
  gearingEnabled: false,
  redlineRPM: "",
  gearCount: "",
  frontTireWidth: "",
  frontTireAspectRatio: "",
  frontTireRimInches: "",
  tireWidth: "",
  tireAspectRatio: "",
  tireRimInches: "",
  targetTopSpeedKmh: "",
  balanceBias: "100",
  stiffnessBias: "100",
  speedBias: "100",
  tractionBias: "100",
  shiftBias: "100",
  trackFeedbackScene: "",
  trackFeedbackSymptom: "",
  correctionStrengthBias: "100",
  driftFeedbackSymptom: "",
  driftCorrectionStrengthBias: "100",
  offroadFeedbackSymptom: "",
  offroadCorrectionStrengthBias: "100",
};

const preferenceSceneOptions = [
  { label: "低速弯", value: "lowSpeedCorner" },
  { label: "中高速弯", value: "highSpeedCorner" },
];

const preferenceSymptomOptions = {
  lowSpeedCorner: [
    { label: "入弯推头", value: "entryUndersteer" },
    { label: "弯中推头", value: "midUndersteer" },
    { label: "出弯打滑", value: "exitWheelspin" },
    { label: "减速不稳", value: "brakingInstability" },
  ],
  highSpeedCorner: [
    { label: "高速入弯不稳", value: "highSpeedEntryInstability" },
    { label: "弯中推头", value: "midUndersteer" },
    { label: "高速甩尾", value: "highSpeedOversteer" },
    { label: "车身跳动", value: "bodyBounce" },
  ],
};

const driftFeedbackOptions = [
  { label: "起漂困难", value: "initiationDifficult" },
  { label: "起漂太敏感", value: "initiationTooSensitive" },
  { label: "角度不够", value: "angleTooSmall" },
  { label: "角度过大难控制", value: "angleTooLarge" },
  { label: "后轮空转过多", value: "rearWheelspin" },
  { label: "后轮抓地过强", value: "rearGripTooStrong" },
];

const offroadFeedbackOptions = [
  { label: "过坑弹飞", value: "bumpKick" },
  { label: "悬挂触底", value: "bottomingOut" },
  { label: "落地不稳", value: "landingUnstable" },
  { label: "转向迟钝", value: "lazyTurnIn" },
  { label: "出弯空转", value: "exitWheelspin" },
  { label: "高速漂浮", value: "highSpeedFloat" },
];

const groupLabels = {
  tire: "轮胎",
  gearing: "齿轮",
  alignment: "轮胎定位",
  antiroll: "防倾杆",
  springs: "弹簧 / 车高",
  damping: "阻尼",
  aero: "空气动力学",
  brake: "刹车",
  differential: "差速器",
};

const groupOrder = ["tire", "gearing", "alignment", "antiroll", "springs", "damping", "aero", "brake", "differential"];

const fieldLabels = {
  weightKG: "车重",
  frontWeightPct: "前配重",
  powerKW: "马力",
  torqueNM: "扭矩",
  redlineRPM: "红线转速",
  frontTirePressure: "前胎压",
  rearTirePressure: "后胎压",
  frontCamber: "前外倾角",
  rearCamber: "后外倾角",
  frontToe: "前束",
  rearToe: "后束",
  caster: "主销后倾角",
  frontArb: "前防倾杆",
  rearArb: "后防倾杆",
  frontSpring: "前弹簧",
  rearSpring: "后弹簧",
  frontRideHeight: "前车高",
  rearRideHeight: "后车高",
  frontRebound: "前回弹阻尼",
  rearRebound: "后回弹阻尼",
  frontBump: "前压缩阻尼",
  rearBump: "后压缩阻尼",
  frontAero: "前下压力",
  rearAero: "后下压力",
  brakeBalance: "刹车平衡",
  brakePressure: "刹车压力",
  frontDiffAccel: "前差速加速",
  frontDiffDecel: "前差速减速",
  rearDiffAccel: "后差速加速",
  rearDiffDecel: "后差速减速",
  centerDiffBalance: "中央差速分配",
  finalDrive: "最终传动齿轮",
  gear1: "1 挡齿轮",
  gear2: "2 挡齿轮",
  gear3: "3 挡齿轮",
  gear4: "4 挡齿轮",
  gear5: "5 挡齿轮",
  gear6: "6 挡齿轮",
  gear7: "7 挡齿轮",
  gear8: "8 挡齿轮",
  gear9: "9 挡齿轮",
  gear10: "10 挡齿轮",
};

const fieldOrder = Object.keys(fieldLabels).reduce((acc, key, index) => {
  acc[key] = index;
  return acc;
}, {});

Page({
  data: {
    config: runtimeConfig,
    form: { ...defaultForm },
    useCaseOptions,
    drivetrainOptions,
    tireCompoundOptions,
    useCaseLabel: optionLabel(useCaseOptions, "Road"),
    drivetrainLabel: optionLabel(drivetrainOptions, "AWD"),
    tireCompoundLabel: optionLabel(tireCompoundOptions, "sport"),
    useCaseIndex: 0,
    drivetrainIndex: 0,
    tireCompoundIndex: 2,
    loading: false,
    preferenceOpen: false,
    preferenceSceneOptions,
    driftFeedbackOptions,
    offroadFeedbackOptions,
    tempPreferenceSymptoms: [],
    tempPreference: {
      balanceBias: 100,
      stiffnessBias: 100,
      speedBias: 100,
      tractionBias: 100,
      shiftBias: 100,
      trackFeedbackScene: "",
      trackFeedbackSymptom: "",
      correctionStrengthBias: 100,
      driftFeedbackSymptom: "",
      driftCorrectionStrengthBias: 100,
      offroadFeedbackSymptom: "",
      offroadCorrectionStrengthBias: 100,
    },
    result: null,
    resultPayload: null,
    resultGroups: [],
    warnings: [],
    nextTestPlan: [],
    savedRecordId: null,
    editingFieldKey: "",
    editingValue: "",
    shareLoading: false,
    shareLinkReady: false,
    shareId: "",
    sharePath: "",
    shareSummary: null,
    resultEditTipVisible: false,
    quotaText: "",
    bannerAdUnitId: runtimeConfig.bannerAdUnitId,
    resultBannerAdUnitId: runtimeConfig.resultBannerAdUnitId,
    hasRewardedAd: Boolean(runtimeConfig.rewardedAdUnitId),
  },

  onLoad() {
    this.rewardedVideoAd = null;
    this.pendingAdResolve = null;
    this.initRewardedAd();
    this.refreshQuotaText();
  },

  onShow() {
    setSelectedTab(this, 0);
    setTabBarHidden(this, this.data.preferenceOpen);
    this.refreshQuotaText();
  },

  onUnload() {
    this.pendingAdResolve = null;
  },

  openRecommendPage() {
    wx.switchTab({
      url: "/pages/recommend/index",
    });
  },

  openUpgradeGuidePage() {
    const useCase = this.data.form.useCase || "Road";
    const carClass = classFromPi(this.data.form.pi || 700);
    wx.navigateTo({
      url: `/pages/upgrade-guide/index?useCase=${useCase}&carClass=${carClass}`,
    });
  },

  onInput(e) {
    const field = e.currentTarget.dataset.field;
    this.setData({
      [`form.${field}`]: e.detail.value,
    });
  },

  onSwitchChange(e) {
    const field = e.currentTarget.dataset.field;
    this.setData({
      [`form.${field}`]: e.detail.value,
    });
  },

  onPickerChange(e) {
    const field = e.currentTarget.dataset.field;
    const index = Number(e.detail.value);
    let options = [];
    if (field === "useCase") options = useCaseOptions;
    if (field === "drivetrain") options = drivetrainOptions;
    if (field === "tireCompound") options = tireCompoundOptions;
    const selected = options[index];
    if (!selected) return;

    const updates = {
      [`form.${field}`]: selected.value,
    };
    if (field === "useCase") {
      const defaultTire = defaultTireCompoundForUseCase(selected.value);
      const tireIndex = tireCompoundOptions.findIndex((item) => item.value === defaultTire);
      updates.useCaseIndex = index;
      updates.useCaseLabel = selected.label;
      updates["form.tireCompound"] = defaultTire;
      updates.tireCompoundIndex = tireIndex === -1 ? this.data.tireCompoundIndex : tireIndex;
      updates.tireCompoundLabel = optionLabel(tireCompoundOptions, defaultTire);
      if (selected.value === "Drift") {
        const rwdIndex = drivetrainOptions.findIndex((item) => item.value === "RWD");
        updates["form.drivetrain"] = "RWD";
        updates.drivetrainIndex = rwdIndex === -1 ? this.data.drivetrainIndex : rwdIndex;
        updates.drivetrainLabel = optionLabel(drivetrainOptions, "RWD");
      }
    }
    if (field === "drivetrain") {
      updates.drivetrainIndex = index;
      updates.drivetrainLabel = selected.label;
    }
    if (field === "tireCompound") {
      updates.tireCompoundIndex = index;
      updates.tireCompoundLabel = selected.label;
    }
    this.setData(updates);
  },

  resetForm() {
    this.setData({
      form: { ...defaultForm },
      useCaseLabel: optionLabel(useCaseOptions, "Road"),
      drivetrainLabel: optionLabel(drivetrainOptions, "AWD"),
      tireCompoundLabel: optionLabel(tireCompoundOptions, "sport"),
      useCaseIndex: 0,
      drivetrainIndex: 0,
      tireCompoundIndex: 2,
      result: null,
      resultPayload: null,
      resultGroups: [],
      warnings: [],
      nextTestPlan: [],
      savedRecordId: null,
      editingFieldKey: "",
      editingValue: "",
      shareLoading: false,
      shareLinkReady: false,
      shareId: "",
      sharePath: "",
      shareSummary: null,
      resultEditTipVisible: false,
      preferenceOpen: false,
      tempPreferenceSymptoms: [],
      tempPreference: {
        balanceBias: 100,
        stiffnessBias: 100,
        speedBias: 100,
        tractionBias: 100,
        shiftBias: 100,
        trackFeedbackScene: "",
        trackFeedbackSymptom: "",
        correctionStrengthBias: 100,
        driftFeedbackSymptom: "",
        driftCorrectionStrengthBias: 100,
        offroadFeedbackSymptom: "",
        offroadCorrectionStrengthBias: 100,
      },
    });
  },

  async generateTune() {
    if (this.data.loading) return;

    let payload;
    try {
      payload = buildPayload(this.data.form);
    } catch (err) {
      wx.showToast({
        title: err.message,
        icon: "none",
        duration: 2600,
      });
      return;
    }

    const app = getApp();
    if (!app.globalData.env) {
      wx.showModal({
        title: "云环境未配置",
        content: "请在 miniprogram/config.js 填写 cloudEnvId，并上传 calculateTune 云函数后再使用。",
        showCancel: false,
      });
      return;
    }

    const allowed = await this.ensureTuneQuota();
    if (!allowed) return;

    this.setData({ loading: true });
    try {
      const response = await wx.cloud.callFunction({
        name: "calculateTune",
        data: {
          type: "calculateTune",
          payload,
        },
      });
      const body = response.result || {};
      if (!body.success) {
        throw new Error(body.message || "云函数计算失败");
      }
      const result = body.data;
      const resultGroups = buildResultGroups(result.generatedFields || [], result.tierRecommendations || [], payload);
      const warnings = result.warnings || [];
      const nextTestPlan = result.nextTestPlan || [];
      const resultEditTipVisible = shouldShowResultEditTip();
      const savedRecord = addTuneRecord({
        payload,
        result,
        resultGroups,
        warnings,
        nextTestPlan,
        useCaseLabel: this.data.useCaseLabel,
        tireCompoundLabel: this.data.tireCompoundLabel,
      });
      this.consumeTuneQuota();
      this.setData({
        result,
        resultPayload: payload,
        resultGroups,
        warnings,
        nextTestPlan,
        savedRecordId: savedRecord ? savedRecord.id : null,
        editingFieldKey: "",
        editingValue: "",
        shareLoading: false,
        shareLinkReady: false,
        shareId: "",
        sharePath: "",
        shareSummary: null,
        resultEditTipVisible,
      });
      if (resultEditTipVisible) {
        markResultEditTipSeen();
      }
      wx.showToast({
        title: "已保存到我的调校",
        icon: "success",
      });
    } catch (err) {
      wx.showModal({
        title: "生成失败",
        content: err && err.message ? err.message : "请检查云函数是否已上传部署。",
        showCancel: false,
      });
    } finally {
      this.setData({ loading: false });
      this.refreshQuotaText();
    }
  },

  openPreferencePanel() {
    const trackFeedbackEnabled = supportsTrackFeedbackUseCase(this.data.form.useCase);
    const scene = trackFeedbackEnabled ? normalizeTrackFeedbackScene(this.data.form.trackFeedbackScene) : "";
    const symptom = normalizeTrackFeedbackSymptom(scene, this.data.form.trackFeedbackSymptom);
    this.setData({
      preferenceOpen: true,
      tempPreferenceSymptoms: preferenceSymptomsForScene(scene),
      tempPreference: {
        balanceBias: preferenceValue(this.data.form.balanceBias),
        stiffnessBias: preferenceValue(this.data.form.stiffnessBias),
        speedBias: preferenceValue(this.data.form.speedBias),
        tractionBias: preferenceValue(this.data.form.tractionBias),
        shiftBias: preferenceValue(this.data.form.shiftBias),
        trackFeedbackScene: scene,
        trackFeedbackSymptom: symptom,
        correctionStrengthBias: preferenceValue(this.data.form.correctionStrengthBias),
        driftFeedbackSymptom: normalizeDriftFeedbackSymptom(this.data.form.driftFeedbackSymptom),
        driftCorrectionStrengthBias: preferenceValue(this.data.form.driftCorrectionStrengthBias),
        offroadFeedbackSymptom: normalizeOffroadFeedbackSymptom(this.data.form.offroadFeedbackSymptom),
        offroadCorrectionStrengthBias: preferenceValue(this.data.form.offroadCorrectionStrengthBias),
      },
    });
    setTabBarHidden(this, true);
  },

  closePreferencePanel() {
    this.setData({ preferenceOpen: false });
    setTabBarHidden(this, false);
  },

  noop() {},

  onPreferenceSliderChange(e) {
    const field = e.currentTarget.dataset.field;
    this.setData({
      [`tempPreference.${field}`]: Number(e.detail.value),
    });
  },

  onPreferenceSceneTap(e) {
    const scene = e.currentTarget.dataset.value || "";
    const current = this.data.tempPreference.trackFeedbackScene;
    const nextScene = current === scene ? "" : scene;
    this.setData({
      tempPreferenceSymptoms: preferenceSymptomsForScene(nextScene),
      "tempPreference.trackFeedbackScene": nextScene,
      "tempPreference.trackFeedbackSymptom": "",
    });
  },

  onPreferenceSymptomTap(e) {
    const symptom = e.currentTarget.dataset.value || "";
    const current = this.data.tempPreference.trackFeedbackSymptom;
    this.setData({
      "tempPreference.trackFeedbackSymptom": current === symptom ? "" : symptom,
    });
  },

  onDriftFeedbackTap(e) {
    const symptom = e.currentTarget.dataset.value || "";
    const current = this.data.tempPreference.driftFeedbackSymptom;
    this.setData({
      "tempPreference.driftFeedbackSymptom": current === symptom ? "" : symptom,
    });
  },

  onOffroadFeedbackTap(e) {
    const symptom = e.currentTarget.dataset.value || "";
    const current = this.data.tempPreference.offroadFeedbackSymptom;
    this.setData({
      "tempPreference.offroadFeedbackSymptom": current === symptom ? "" : symptom,
    });
  },

  applyPreferencePanel() {
    const prefs = this.data.tempPreference;
    const trackFeedbackEnabled = supportsTrackFeedbackUseCase(this.data.form.useCase);
    const scene = trackFeedbackEnabled ? normalizeTrackFeedbackScene(prefs.trackFeedbackScene) : "";
    const symptom = normalizeTrackFeedbackSymptom(scene, prefs.trackFeedbackSymptom);
    const driftSymptom = this.data.form.useCase === "Drift" ? normalizeDriftFeedbackSymptom(prefs.driftFeedbackSymptom) : "";
    const offroadSymptom = this.data.form.useCase === "Offroad" ? normalizeOffroadFeedbackSymptom(prefs.offroadFeedbackSymptom) : "";
    this.setData({
      preferenceOpen: false,
      "form.balanceBias": String(prefs.balanceBias),
      "form.stiffnessBias": String(prefs.stiffnessBias),
      "form.speedBias": String(prefs.speedBias),
      "form.tractionBias": String(prefs.tractionBias),
      "form.shiftBias": String(prefs.shiftBias),
      "form.trackFeedbackScene": scene,
      "form.trackFeedbackSymptom": symptom,
      "form.correctionStrengthBias": String(prefs.correctionStrengthBias),
      "form.driftFeedbackSymptom": driftSymptom,
      "form.driftCorrectionStrengthBias": String(prefs.driftCorrectionStrengthBias),
      "form.offroadFeedbackSymptom": offroadSymptom,
      "form.offroadCorrectionStrengthBias": String(prefs.offroadCorrectionStrengthBias),
    });
    setTabBarHidden(this, false);
    setTimeout(() => this.generateTune(), 0);
  },

  onResultValueTap(e) {
    const fieldKey = e.currentTarget.dataset.fieldKey;
    const editable = e.currentTarget.dataset.editable;
    if (!fieldKey || !(editable === true || editable === "true")) return;

    const now = Date.now();
    const lastTap = this.lastValueTap || {};
    this.lastValueTap = { fieldKey, time: now };
    if (lastTap.fieldKey !== fieldKey || now - lastTap.time > 360) return;

    const field = findResultField(this.data.resultGroups, fieldKey);
    if (!field) return;
    this.setData({
      editingFieldKey: fieldKey,
      editingValue: field.displayValue,
    });
  },

  onEditInput(e) {
    this.setData({
      editingValue: e.detail.value,
    });
  },

  finishFieldEdit(e) {
    const fieldKey = e.currentTarget.dataset.fieldKey || this.data.editingFieldKey;
    if (!fieldKey || this.data.editingFieldKey !== fieldKey) return;

    let parsed;
    try {
      parsed = parseManualDisplayValue(fieldKey, this.data.editingValue);
    } catch (err) {
      wx.showToast({
        title: err.message || "请输入数字",
        icon: "none",
      });
      this.setData({
        editingFieldKey: "",
        editingValue: "",
      });
      return;
    }

    const resultGroups = updateGroupsFieldValue(this.data.resultGroups, fieldKey, parsed);
    const result = updateResultFieldValue(this.data.result, fieldKey, parsed.value);
    if (this.data.savedRecordId) {
      updateTuneRecord(this.data.savedRecordId, {
        result,
        resultGroups,
      });
    }
    this.setData({
      result,
      resultGroups,
      editingFieldKey: "",
      editingValue: "",
      shareLinkReady: false,
      shareId: "",
      sharePath: "",
      shareSummary: null,
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

  async prepareShareTune() {
    if (this.data.shareLoading) return;
    if (!this.data.result) {
      wx.showToast({
        title: "请先生成调校",
        icon: "none",
      });
      return;
    }

    this.setData({ shareLoading: true });
    try {
      const summary = buildShareSummary(this.data);
      const data = await createTuneShare({
        summary,
        result: this.data.result,
        resultGroups: this.data.resultGroups,
        warnings: this.data.warnings,
        nextTestPlan: this.data.nextTestPlan,
      });
      const sharePath = buildSharePath(data.shareId);
      this.setData({
        shareLinkReady: true,
        shareId: data.shareId,
        sharePath,
        shareSummary: summary,
      });
      wx.showToast({
        title: "链接已生成",
        icon: "success",
      });
    } catch (err) {
      wx.showModal({
        title: "分享失败",
        content: err && err.message ? err.message : "分享链接生成失败，请稍后再试。",
        showCancel: false,
      });
    } finally {
      this.setData({ shareLoading: false });
    }
  },

  onShareAppMessage() {
    const summary = this.data.shareSummary || buildShareSummary(this.data);
    return {
      title: shareTitle(summary),
      path: this.data.sharePath || "/pages/index/index",
    };
  },

  closeResultEditTip() {
    markResultEditTipSeen();
    this.setData({
      resultEditTipVisible: false,
    });
  },

  initRewardedAd() {
    if (!runtimeConfig.rewardedAdUnitId || !wx.createRewardedVideoAd) {
      return;
    }
    this.rewardedVideoAd = wx.createRewardedVideoAd({
      adUnitId: runtimeConfig.rewardedAdUnitId,
    });
    this.rewardedVideoAd.onClose((res) => {
      if (!this.pendingAdResolve) return;
      const completed = !res || res.isEnded;
      const resolve = this.pendingAdResolve;
      this.pendingAdResolve = null;
      resolve(completed ? "completed" : "closed");
    });
    this.rewardedVideoAd.onError(() => {
      if (!this.pendingAdResolve) return;
      const resolve = this.pendingAdResolve;
      this.pendingAdResolve = null;
      this.showRewardedAdUnavailable();
      resolve("unavailable");
    });
  },

  async ensureTuneQuota() {
    const state = loadAdState();
    if (state.rewardRemaining > 0 || state.freeUsed < runtimeConfig.freeTuneCount) {
      return true;
    }

    if (!runtimeConfig.rewardedAdUnitId) {
      return true;
    }

    const adResult = await this.showRewardedAd();
    if (adResult !== "completed") {
      if (adResult === "closed") {
        wx.showToast({
          title: "看完广告后才能继续",
          icon: "none",
        });
      }
      this.refreshQuotaText();
      return false;
    }

    const refreshedState = loadAdState();
    refreshedState.rewardRemaining = runtimeConfig.rewardTuneCount;
    saveAdState(refreshedState);
    this.refreshQuotaText();
    return true;
  },

  showRewardedAd() {
    if (!this.rewardedVideoAd) {
      this.showRewardedAdUnavailable();
      return Promise.resolve("unavailable");
    }
    return new Promise((resolve) => {
      this.pendingAdResolve = resolve;
      this.rewardedVideoAd.show().catch(() => {
        this.rewardedVideoAd.load()
          .then(() => this.rewardedVideoAd.show())
          .catch(() => {
            if (!this.pendingAdResolve) return;
            this.pendingAdResolve = null;
            this.showRewardedAdUnavailable();
            resolve("unavailable");
          });
      });
    });
  },

  showRewardedAdUnavailable() {
    wx.showModal({
      title: "广告暂不可用",
      content: "激励视频加载失败，请稍后再试。",
      showCancel: false,
    });
  },

  consumeTuneQuota() {
    const state = loadAdState();
    if (state.rewardRemaining > 0) {
      state.rewardRemaining -= 1;
    } else if (state.freeUsed < runtimeConfig.freeTuneCount) {
      state.freeUsed += 1;
    }
    saveAdState(state);
  },

  refreshQuotaText() {
    const state = loadAdState();
    let quotaText = "";
    if (!runtimeConfig.rewardedAdUnitId) {
      quotaText = "";
    } else if (state.rewardRemaining > 0) {
      quotaText = `广告解锁剩余 ${state.rewardRemaining} 次`;
    } else {
      const freeLeft = Math.max(0, runtimeConfig.freeTuneCount - state.freeUsed);
      quotaText = freeLeft > 0 ? `免费调校剩余 ${freeLeft} 次` : "本次需观看激励视频";
    }
    this.setData({ quotaText });
  },

});

function buildPayload(form) {
  const gearingEnabled = Boolean(form.gearingEnabled);
  const rearTireDiameterCm = gearingEnabled ? tireDiameterCmFromParts(form.tireWidth, form.tireAspectRatio, form.tireRimInches) : null;
  const frontTireDiameterCm = gearingEnabled ? optionalTireDiameterCmFromParts(form.frontTireWidth, form.frontTireAspectRatio, form.frontTireRimInches, "前轮尺寸") : null;
  return {
    carName: "",
    versionName: "",
    useCase: form.useCase,
    pi: parseIntegerField(form.pi, "PI", 100, 999),
    drivetrain: form.drivetrain,
    tireCompound: form.tireCompound,
    weightKG: parseIntegerField(form.weightKG, "车重", 300, 3000),
    frontWeightPct: parseIntegerField(form.frontWeightPct, "前配重", 1, 99),
    powerKW: gearingEnabled ? parseOptionalNumberField(form.powerKW, "马力") : null,
    torqueNM: gearingEnabled ? parseOptionalNumberField(form.torqueNM, "扭矩") : null,
    redlineRPM: gearingEnabled ? parseIntegerField(form.redlineRPM, "红线转速", 1000, 20000) : null,
    gearCount: gearingEnabled ? parseIntegerField(form.gearCount, "挡位数", 2, 10) : null,
    tireDiameterCm: rearTireDiameterCm,
    frontTireDiameterCm,
    rearTireDiameterCm,
    frontTireWidthMm: gearingEnabled ? parseOptionalTireWidthField(form.frontTireWidth, "前轮宽度") : null,
    rearTireWidthMm: gearingEnabled ? parseIntegerField(form.tireWidth, "后轮宽度", 100, 455) : null,
    targetTopSpeedKmh: gearingEnabled ? parseIntegerField(
      form.targetTopSpeedKmh,
      form.useCase === "Drift" ? "目标漂移速度" : "目标速度",
      form.useCase === "Drift" ? 20 : 1,
      form.useCase === "Drift" ? 260 : 600
    ) : null,
    frontRideHeightMinCm: null,
    frontRideHeightMaxCm: null,
    rearRideHeightMinCm: null,
    rearRideHeightMaxCm: null,
    frontAeroMinKgf: null,
    frontAeroMaxKgf: null,
    rearAeroMinKgf: null,
    rearAeroMaxKgf: null,
    frontRideHeightAdjustable: false,
    rearRideHeightAdjustable: false,
    frontAeroAdjustable: false,
    rearAeroAdjustable: false,
    balanceBias: parseIntegerField(form.balanceBias, "驾驶风格", 50, 150),
    stiffnessBias: parseIntegerField(form.stiffnessBias, "悬挂支撑", 50, 150),
    speedBias: parseIntegerField(form.speedBias, "齿轮取向", 50, 150),
    tractionBias: parseIntegerField(form.tractionBias, "牵引偏好", 50, 150),
    shiftBias: parseIntegerField(form.shiftBias, "换挡偏好", 50, 150),
    trackFeedbackScene: supportsTrackFeedbackUseCase(form.useCase) ? normalizeTrackFeedbackScene(form.trackFeedbackScene) : "",
    trackFeedbackSymptom: supportsTrackFeedbackUseCase(form.useCase) ? normalizeTrackFeedbackSymptom(form.trackFeedbackScene, form.trackFeedbackSymptom) : "",
    correctionStrengthBias: parseIntegerField(form.correctionStrengthBias, "赛道修正强度", 50, 150),
    driftFeedbackSymptom: form.useCase === "Drift" ? normalizeDriftFeedbackSymptom(form.driftFeedbackSymptom) : "",
    driftCorrectionStrengthBias: parseIntegerField(form.driftCorrectionStrengthBias, "漂移修正强度", 50, 150),
    offroadFeedbackSymptom: form.useCase === "Offroad" ? normalizeOffroadFeedbackSymptom(form.offroadFeedbackSymptom) : "",
    offroadCorrectionStrengthBias: parseIntegerField(form.offroadCorrectionStrengthBias, "越野修正强度", 50, 150),
  };
}

function parseIntegerField(value, label, min, max) {
  const raw = String(value || "").trim();
  if (!/^-?\d+$/.test(raw)) {
    throw new Error(`${label}必须是整数`);
  }
  const parsed = Number(raw);
  if (parsed < min || parsed > max) {
    throw new Error(`${label}范围 ${min}-${max}`);
  }
  return parsed;
}

function tireDiameterCmFromParts(widthValue, aspectValue, rimValue) {
  const widthMm = parseIntegerField(widthValue, "轮胎宽度", 100, 455);
  const aspectRatio = parseIntegerField(aspectValue, "轮胎扁平比", 20, 100);
  const rimInches = parseIntegerField(rimValue, "轮毂尺寸", 10, 30);
  const diameterCm = (rimInches * 25.4 + 2 * widthMm * (aspectRatio / 100)) / 10;
  if (diameterCm < 40 || diameterCm > 120) {
    throw new Error("轮胎尺寸换算直径需在 40-120 cm");
  }
  return Math.round(diameterCm * 100) / 100;
}

function optionalTireDiameterCmFromParts(widthValue, aspectValue, rimValue, label) {
  const parts = [widthValue, aspectValue, rimValue].map((value) => String(value || "").trim());
  if (parts.every((value) => !value)) return null;
  if (parts.some((value) => !value)) {
    throw new Error(`${label}需完整填写宽度、扁平比和轮毂尺寸`);
  }
  return tireDiameterCmFromParts(widthValue, aspectValue, rimValue);
}

function parseOptionalTireWidthField(value, label) {
  const raw = String(value || "").trim();
  if (!raw) return null;
  return parseIntegerField(raw, label, 100, 455);
}

function parseOptionalIntegerField(value, label, min, max) {
  const raw = String(value || "").trim();
  if (!raw) return null;
  return parseIntegerField(raw, label, min, max);
}

function parseOptionalNumberField(value, label, min, max) {
  const raw = String(value || "").trim();
  if (!raw) return null;
  const parsed = Number(raw);
  if (!Number.isFinite(parsed)) {
    throw new Error(`${label}必须是数字`);
  }
  if (typeof min === "number" && parsed < min) {
    throw new Error(`${label}不能小于 ${min}`);
  }
  if (typeof max === "number" && parsed > max) {
    throw new Error(`${label}不能大于 ${max}`);
  }
  return parsed;
}

function loadAdState() {
  try {
    const raw = wx.getStorageSync(storageKey);
    if (!raw) return { freeUsed: 0, rewardRemaining: 0 };
    return {
      freeUsed: Number(raw.freeUsed) || 0,
      rewardRemaining: Number(raw.rewardRemaining) || 0,
    };
  } catch (err) {
    return { freeUsed: 0, rewardRemaining: 0 };
  }
}

function saveAdState(state) {
  try {
    wx.setStorageSync(storageKey, {
      freeUsed: Math.max(0, Number(state.freeUsed) || 0),
      rewardRemaining: Math.max(0, Number(state.rewardRemaining) || 0),
    });
  } catch (err) {
    // Local quota is a soft limit; storage failure should not break tuning.
  }
}

function buildResultGroups(fields, tierRecommendations, payload) {
  const groups = {};
  const visibleFields = fields.filter((field) => field.group !== "power");
  visibleFields.forEach((field) => {
    const group = field.group || "other";
    if (!groups[group]) {
      groups[group] = {
        key: group,
        label: groupLabels[group] || group,
        items: [],
      };
    }
    groups[group].items.push({
      ...field,
      label: fieldLabels[field.fieldKey] || field.fieldKey,
      displayValue: formatFieldValue(field),
    });
  });

  if (!payload.gearCount) {
    groups.gearing = {
      key: "gearing",
      label: groupLabels.gearing,
      items: [{
        fieldKey: "gearingDisabled",
        label: "齿轮",
        displayValue: "请开启齿轮调校",
        unit: "",
      }],
    };
  }

  const rideItems = tierRecommendations
    .filter((item) => item.group === "springs" && (item.fieldKey === "frontRideHeight" || item.fieldKey === "rearRideHeight"))
    .map((item) => ({
      fieldKey: item.fieldKey,
      label: fieldLabels[item.fieldKey] || item.fieldKey,
      displayValue: fiveLevelLabel(item.tier, ["最低", "低", "中", "高", "最高"]),
      unit: "",
      tierItem: true,
    }));
  if (rideItems.length > 0) {
    if (!groups.springs) {
      groups.springs = { key: "springs", label: groupLabels.springs, items: [] };
    }
    groups.springs.items.push(...rideItems);
  }

  const aeroItems = tierRecommendations
    .filter((item) => item.group === "aero")
    .map((item) => ({
      fieldKey: item.fieldKey,
      label: fieldLabels[item.fieldKey] || item.fieldKey,
      displayValue: fiveLevelLabel(item.tier, ["速度", "偏速度", "均衡", "偏过弯", "过弯"]),
      unit: "",
      tierItem: true,
    }));
  if (aeroItems.length > 0) {
    groups.aero = {
      key: "aero",
      label: groupLabels.aero,
      items: aeroItems,
    };
  }

  const resultGroups = Object.keys(groups)
    .sort((a, b) => groupIndex(a) - groupIndex(b))
    .map((key) => {
      groups[key].items.sort((a, b) => fieldIndex(a.fieldKey) - fieldIndex(b.fieldKey));
      return groups[key];
    });
  return decorateResultGroups(resultGroups);
}

function formatFieldValue(field) {
  const value = field.value;
  if (value === null || value === undefined) return "--";
  if (typeof value !== "number") return String(value);
  const step = fieldStep(field.fieldKey);
  const decimals = step >= 1 ? 0 : step >= 0.1 ? 1 : 2;
  return value.toFixed(decimals);
}

function fieldStep(key) {
  if (key === "frontTirePressure" || key === "rearTirePressure") return 0.01;
  if (key === "finalDrive" || /^gear\d+$/.test(key)) return 0.01;
  if (key === "frontArb" || key === "rearArb") return 0.1;
  if (["frontAero", "rearAero", "brakeBalance", "brakePressure", "frontDiffAccel", "frontDiffDecel", "rearDiffAccel", "rearDiffDecel", "centerDiffBalance"].includes(key)) return 1;
  return 0.1;
}

function fiveLevelLabel(tier, labels) {
  if (tier === "lowest") return labels[0];
  if (tier === "low") return labels[1];
  if (tier === "medium") return labels[2];
  if (tier === "high") return labels[3];
  if (tier === "highest") return labels[4];
  return labels[2];
}

function groupIndex(group) {
  const index = groupOrder.indexOf(group);
  return index === -1 ? 99 : index;
}

function fieldIndex(key) {
  return fieldOrder[key] === undefined ? 999 : fieldOrder[key];
}

function optionLabel(options, value) {
  const found = options.find((item) => item.value === value);
  return found ? found.label : "";
}

function defaultTireCompoundForUseCase(useCase) {
  if (useCase === "Rally") return "rally";
  if (useCase === "Offroad") return "offroad";
  if (useCase === "Drag") return "drag";
  if (useCase === "Drift") return "drift";
  return "sport";
}

function preferenceValue(value) {
  const parsed = Number(value);
  if (!Number.isFinite(parsed)) return 100;
  return Math.min(150, Math.max(50, Math.round(parsed)));
}

function preferenceSymptomsForScene(scene) {
  return preferenceSymptomOptions[scene] || [];
}

function normalizeTrackFeedbackScene(scene) {
  return preferenceSceneOptions.some((item) => item.value === scene) ? scene : "";
}

function normalizeTrackFeedbackSymptom(scene, symptom) {
  const normalizedScene = normalizeTrackFeedbackScene(scene);
  if (!normalizedScene) return "";
  return preferenceSymptomsForScene(normalizedScene).some((item) => item.value === symptom) ? symptom : "";
}

function normalizeDriftFeedbackSymptom(symptom) {
  return driftFeedbackOptions.some((item) => item.value === symptom) ? symptom : "";
}

function normalizeOffroadFeedbackSymptom(symptom) {
  return offroadFeedbackOptions.some((item) => item.value === symptom) ? symptom : "";
}

function supportsTrackFeedbackUseCase(useCase) {
  return useCase === "Road" || useCase === "Rally";
}

function shouldShowResultEditTip() {
  try {
    return !wx.getStorageSync(resultEditTipStorageKey);
  } catch (err) {
    return true;
  }
}

function markResultEditTipSeen() {
  try {
    wx.setStorageSync(resultEditTipStorageKey, true);
  } catch (err) {
    // This hint is optional; storage failure should not affect tuning.
  }
}

function buildShareSummary(data) {
  const result = data.result || {};
  const profile = result.profileDraft || {};
  const payload = data.resultPayload || {};
  const useCase = payload.useCase || profile.useCase || "";
  const pi = Number(profile.pi || payload.pi || 0);
  return {
    useCase,
    useCaseLabel: data.useCaseLabel || useCaseLabel(useCase),
    carClass: profile.carClass || (pi ? classFromPi(pi) : ""),
    pi,
    drivetrain: profile.drivetrain || payload.drivetrain || "",
    tireCompoundLabel: data.tireCompoundLabel || tireCompoundLabel(payload.tireCompound),
  };
}

function classFromPi(pi) {
  const value = Number(pi) || 0;
  if (value >= 999) return "X";
  if (value >= 901) return "R";
  if (value >= 801) return "S2";
  if (value >= 701) return "S1";
  if (value >= 601) return "A";
  if (value >= 501) return "B";
  if (value >= 401) return "C";
  return "D";
}

function useCaseLabel(value) {
  const labels = {
    Road: "公路",
    Rally: "拉力",
    Offroad: "越野",
    Drag: "直线",
    Drift: "漂移",
  };
  return labels[value] || value || "";
}

function tireCompoundLabel(value) {
  const labels = {
    stock: "标准",
    street: "街车",
    sport: "跑车",
    semi: "半热熔",
    slick: "热熔胎",
    rally: "拉力",
    offroad: "越野",
    drift: "漂移",
    drag: "直线",
    snow: "雪地",
  };
  return labels[value] || value || "";
}
