const { classTabs, guideData, useCaseTabs } = require("../../data/upgradeGuide");

const priorityMeta = [
  { key: "p0", label: "P0 必选", tone: "p0" },
  { key: "p1", label: "P1 高优先级", tone: "p1" },
  { key: "p2", label: "P2 视情况", tone: "p2" },
  { key: "p3", label: "P3 最后补", tone: "p3" },
  { key: "avoid", label: "不推荐", tone: "avoid" },
];

Page({
  data: {
    useCaseTabs,
    classTabs,
    priorityMeta,
    activeUseCase: "Road",
    activeClass: "A",
    useCaseIndex: 0,
    classIndex: 4,
    guide: buildGuide("Road", "A"),
  },

  onLoad(options) {
    const useCase = normalizeUseCase(options && options.useCase) || "Road";
    const carClass = normalizeClass(options && options.carClass) || "A";
    this.applyGuide(useCase, carClass);
  },

  onUseCaseTap(e) {
    const index = Number(e.currentTarget.dataset.index);
    const tab = useCaseTabs[index];
    if (!tab) return;
    this.applyGuide(tab.value, this.data.activeClass);
  },

  onClassTap(e) {
    const index = Number(e.currentTarget.dataset.index);
    const tab = classTabs[index];
    if (!tab) return;
    this.applyGuide(this.data.activeUseCase, tab.value);
  },

  openQuickTune() {
    wx.switchTab({
      url: "/pages/index/index",
    });
  },

  openMyTunes() {
    wx.switchTab({
      url: "/pages/my-tunes/index",
    });
  },

  openRecommend() {
    wx.switchTab({
      url: "/pages/recommend/index",
    });
  },

  applyGuide(useCase, carClass) {
    this.setData({
      activeUseCase: useCase,
      activeClass: carClass,
      useCaseIndex: indexOfValue(useCaseTabs, useCase),
      classIndex: indexOfValue(classTabs, carClass),
      guide: buildGuide(useCase, carClass),
    });
  },
});

function buildGuide(useCase, carClass) {
  const section = guideData[useCase] || guideData.Road;
  const selectedGroup = findClassGroup(section.groups, carClass) || section.groups[0];
  return {
    ...section,
    classLabel: carClass,
    groupSummary: selectedGroup.summary,
    classRangeLabel: selectedGroup.classes.join(" / "),
    priorities: priorityMeta
      .map((meta) => ({
        ...meta,
        items: selectedGroup.priorities[meta.key] || [],
      }))
      .filter((item) => item.items.length > 0),
  };
}

function findClassGroup(groups, carClass) {
  if (!Array.isArray(groups)) return null;
  return groups.find((group) => group.classes.includes(carClass));
}

function normalizeUseCase(value) {
  const raw = String(value || "").trim();
  return useCaseTabs.some((item) => item.value === raw) ? raw : "";
}

function normalizeClass(value) {
  const raw = String(value || "").trim().toUpperCase();
  return classTabs.some((item) => item.value === raw) ? raw : "";
}

function indexOfValue(list, value) {
  const index = list.findIndex((item) => item.value === value);
  return index === -1 ? 0 : index;
}
