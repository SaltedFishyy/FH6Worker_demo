const { setSelectedTab, setTabBarHidden } = require("../../utils/tabBar");
const config = require("../../config");
const fallbackRecommendations = require("../../data/recommendedCars");

const useCaseTabs = [
  { label: "公路", value: "Road", icon: "../../images/icons/use-road.svg", activeIcon: "../../images/icons/use-road-active.svg" },
  { label: "拉力", value: "Rally", icon: "../../images/icons/use-rally.svg", activeIcon: "../../images/icons/use-rally-active.svg" },
  { label: "越野", value: "Offroad", icon: "../../images/icons/use-offroad.svg", activeIcon: "../../images/icons/use-offroad-active.svg" },
  { label: "直线", value: "Drag", icon: "../../images/icons/use-drag.svg", activeIcon: "../../images/icons/use-drag-active.svg" },
  { label: "漂移", value: "Drift", icon: "../../images/icons/use-drift.svg", activeIcon: "../../images/icons/use-drift-active.svg" },
];

const classTabs = [
  { label: "X", value: "X" },
  { label: "R", value: "R" },
  { label: "S2", value: "S2" },
  { label: "S1", value: "S1" },
  { label: "A", value: "A" },
  { label: "B", value: "B" },
  { label: "C", value: "C" },
  { label: "D", value: "D" },
];

const defaultUseCase = "Road";
const fallbackCars = parseRecommendationsPayload(fallbackRecommendations);
const defaultClass = firstAvailableClass(defaultUseCase, fallbackCars);
const defaultClassIndex = classIndex(defaultClass);
const recommendationsCacheKey = "fh6_recommendations_cache_v1";
const recommendationsCacheTtlMs = 24 * 60 * 60 * 1000;

Page({
  data: {
    useCaseTabs,
    classTabs,
    recommendedCars: fallbackCars,
    activeUseCase: defaultUseCase,
    activeClass: defaultClass,
    activeTabIndex: 0,
    activeClassIndex: defaultClassIndex,
    expandedCarIds: {},
    cars: buildCars(defaultUseCase, defaultClass, {}, fallbackCars),
  },

  onLoad(options = {}) {
    this.enableShareMenu();
    this.applyRouteFilters(options);
    this.resolveVisibleCarImages();
    this.loadRecommendations();
  },

  onShow() {
    setSelectedTab(this, 2);
    setTabBarHidden(this, false);
  },

  onTabTap(e) {
    const index = Number(e.currentTarget.dataset.index);
    const tab = useCaseTabs[index];
    if (!tab) return;
    const nextClass = firstAvailableClass(tab.value, this.data.recommendedCars);
    this.setData({
      activeUseCase: tab.value,
      activeClass: nextClass,
      activeTabIndex: index,
      activeClassIndex: classIndex(nextClass),
      cars: buildCars(tab.value, nextClass, this.data.expandedCarIds, this.data.recommendedCars),
    }, () => this.resolveVisibleCarImages());
  },

  onClassTap(e) {
    const index = Number(e.currentTarget.dataset.index);
    const tab = classTabs[index];
    if (!tab) return;
    this.setData({
      activeClass: tab.value,
      activeClassIndex: index,
      cars: buildCars(this.data.activeUseCase, tab.value, this.data.expandedCarIds, this.data.recommendedCars),
    }, () => this.resolveVisibleCarImages());
  },

  toggleCar(e) {
    const id = e.currentTarget.dataset.id;
    if (!id) return;
    const expandedCarIds = { ...this.data.expandedCarIds };
    if (expandedCarIds[id]) {
      delete expandedCarIds[id];
    } else {
      expandedCarIds[id] = true;
    }
    this.setData({
      expandedCarIds,
      cars: applyExpandedState(this.data.cars, expandedCarIds),
    });
  },

  onCarImageError(e) {
    const id = e.currentTarget.dataset.id;
    const cars = this.data.cars.map((car) => (
      car.id === id ? { ...car, imageDisplaySrc: "" } : car
    ));
    this.setData({ cars });
  },

  enableShareMenu() {
    if (!wx.showShareMenu) return;
    wx.showShareMenu({
      withShareTicket: false,
      menus: ["shareAppMessage", "shareTimeline"],
    });
  },

  applyRouteFilters(options) {
    const activeUseCase = normalizeUseCase(options && options.useCase) || this.data.activeUseCase;
    const requestedClass = normalizeCarClass(options && options.carClass);
    const activeClass = requestedClass || firstAvailableClass(activeUseCase, this.data.recommendedCars);

    this.setData({
      activeUseCase,
      activeClass,
      activeTabIndex: useCaseIndex(activeUseCase),
      activeClassIndex: classIndex(activeClass),
      cars: buildCars(activeUseCase, activeClass, this.data.expandedCarIds, this.data.recommendedCars),
    });
  },

  onShareAppMessage() {
    return {
      title: buildRecommendShareTitle(this.data),
      path: buildRecommendSharePath(this.data),
    };
  },

  onShareTimeline() {
    return {
      title: buildRecommendShareTitle(this.data),
      query: buildRecommendShareQuery(this.data),
    };
  },

  async loadRecommendations() {
    const fileID = stringValue(config.recommendationsFileId);
    const cachedRecommendations = readCachedRecommendations(fileID);
    if (cachedRecommendations.cars.length) {
      this.applyRecommendedCars(cachedRecommendations.cars);
    }

    await this.loadRemoteRecommendations(fileID, cachedRecommendations);
  },

  async loadRemoteRecommendations(fileIDValue, cachedRecommendations = null) {
    const fileID = stringValue(fileIDValue || config.recommendationsFileId);
    if (!fileID || !wx.cloud || !wx.cloud.downloadFile) return;

    try {
      const downloadResult = await cloudDownloadFile(fileID);
      const text = await readFileText(downloadResult.tempFilePath);
      const remoteRecommendations = parseRecommendationsDocument(JSON.parse(text));
      if (!remoteRecommendations.cars.length) {
        console.warn("[recommend] 云端推荐为空，继续使用本地推荐。");
        return;
      }

      if (isSameRecommendationsVersion(cachedRecommendations, remoteRecommendations)) {
        touchCachedRecommendations(fileID, remoteRecommendations);
        return;
      }

      writeCachedRecommendations(fileID, remoteRecommendations);
      this.applyRecommendedCars(remoteRecommendations.cars);
    } catch (err) {
      console.warn("[recommend] 云端推荐下载失败，继续使用本地推荐。", err);
    }
  },

  applyRecommendedCars(sourceCars) {
    if (!Array.isArray(sourceCars) || !sourceCars.length) return;
    const activeUseCase = this.data.activeUseCase;
    const activeClass = hasCarsForClass(sourceCars, activeUseCase, this.data.activeClass)
      ? this.data.activeClass
      : firstAvailableClass(activeUseCase, sourceCars);

    this.setData({
      recommendedCars: sourceCars,
      activeClass,
      activeClassIndex: classIndex(activeClass),
      cars: buildCars(activeUseCase, activeClass, this.data.expandedCarIds, sourceCars),
    }, () => this.resolveVisibleCarImages());
  },

  resolveVisibleCarImages() {
    if (!wx.cloud || !wx.cloud.getTempFileURL) return;
    const cloudImageCars = this.data.cars.filter((car) => isCloudFileId(car.imageSrc));
    if (!cloudImageCars.length) return;

    wx.cloud.getTempFileURL({
      fileList: cloudImageCars.map((car) => car.imageSrc),
      success: (res) => {
        const urlByFileId = {};
        (res.fileList || []).forEach((item) => {
          if (item.status === 0 && item.tempFileURL) {
            urlByFileId[item.fileID] = item.tempFileURL;
          } else if (item.fileID) {
            console.warn("[recommend] 车辆图片临时链接获取失败。", {
              fileID: item.fileID,
              status: item.status,
              errMsg: item.errMsg,
            });
          }
        });

        const cars = this.data.cars.map((car) => {
          const tempUrl = urlByFileId[car.imageSrc];
          return tempUrl ? { ...car, imageDisplaySrc: tempUrl } : car;
        });
        this.setData({ cars });
      },
      fail: (err) => {
        console.warn("[recommend] 车辆图片临时链接批量获取失败。", err);
      },
    });
  },
});

function buildCars(useCase, carClass, expandedCarIds, sourceCars) {
  return filterCars(useCase, carClass, sourceCars).map((car) => ({
    ...car,
    expanded: Boolean(expandedCarIds && expandedCarIds[car.id]),
    imageDisplaySrc: initialImageDisplaySrc(car),
  }));
}

function applyExpandedState(cars, expandedCarIds) {
  return cars.map((car) => ({
    ...car,
    expanded: Boolean(expandedCarIds && expandedCarIds[car.id]),
  }));
}

function initialImageDisplaySrc(car) {
  if (car.imageSrc && car.imageSrc.indexOf("cloud://") !== 0) {
    return car.imageSrc;
  }
  return car.localImageSrc || "";
}

function isCloudFileId(value) {
  return stringValue(value).indexOf("cloud://") === 0;
}

function filterCars(useCase, carClass, sourceCars) {
  const cars = Array.isArray(sourceCars) ? sourceCars : [];
  return cars.filter((item) => {
    if (item.useCase !== useCase) return false;
    return item.carClass === carClass;
  });
}

function parseRecommendationsPayload(payload) {
  const cars = Array.isArray(payload) ? payload : payload && (payload.cars || payload.recommendedCars);
  return normalizeRecommendedCars(cars);
}

function parseRecommendationsDocument(payload) {
  return {
    version: Array.isArray(payload) ? "" : stringValue(payload && payload.version),
    updatedAt: Array.isArray(payload) ? "" : stringValue(payload && payload.updatedAt),
    cars: parseRecommendationsPayload(payload),
  };
}

function normalizeRecommendedCars(cars) {
  if (!Array.isArray(cars)) return [];
  return mergeRecommendedCars(cars.map((car, index) => normalizeRecommendedCar(car, index)).filter(Boolean));
}

function normalizeRecommendedCar(car, index) {
  if (!car || typeof car !== "object") return null;
  const useCase = normalizeUseCase(car.useCase);
  const carClass = normalizeCarClass(car.carClass || car.class);
  const name = stringValue(car.name);
  const drivetrain = normalizeDrivetrain(car.drivetrain);
  const tireCompound = stringValue(car.tireCompound || "sport");
  if (!useCase || !carClass || !name || !drivetrain) return null;

  const tuneCodes = normalizeTuneCodes(car.tuneCodes, car.tuneCode);
  return {
    id: stringValue(car.id) || `${useCase}-${carClass}-${index}`,
    name,
    useCase,
    useCaseLabel: stringValue(car.useCaseLabel) || useCaseLabel(useCase),
    pi: clampInt(car.pi, 100, 999, 700),
    carClass,
    drivetrain,
    tireCompound,
    tireCompoundLabel: stringValue(car.tireCompoundLabel) || tireCompoundLabel(tireCompound),
    weightKG: clampInt(car.weightKG, 300, 3000, 1400),
    frontWeightPct: clampInt(car.frontWeightPct, 1, 99, 54),
    tuneCode: tuneCodes[0] || "",
    tuneCodes,
    tuneCodeCount: tuneCodes.length,
    multipleTuneCodes: tuneCodes.length > 1,
    imageSrc: stringValue(car.imageSrc || car.imageFileId || car.imageUrl),
    localImageSrc: stringValue(car.localImageSrc),
    tags: Array.isArray(car.tags) ? car.tags.map(stringValue).filter(Boolean).slice(0, 6) : [],
    reason: stringValue(car.reason || car.recommendedFor),
  };
}

function mergeRecommendedCars(cars) {
  const merged = [];
  const byKey = {};
  cars.forEach((car) => {
    const key = recommendedCarMergeKey(car);
    if (!byKey[key]) {
      byKey[key] = {
        ...car,
        tuneCode: "",
        tuneCodes: [],
        tuneCodeCount: 0,
        multipleTuneCodes: false,
      };
      merged.push(byKey[key]);
    }
    appendTuneCodes(byKey[key], car.tuneCodes);
  });

  return merged.map((car) => ({
    ...car,
    tuneCode: car.tuneCodes[0] || "",
    tuneCodeCount: car.tuneCodes.length,
    multipleTuneCodes: car.tuneCodes.length > 1,
  }));
}

function recommendedCarMergeKey(car) {
  return [
    normalizeMergeText(car.name),
    car.useCase,
    car.carClass,
    car.pi,
  ].join("|");
}

function normalizeMergeText(value) {
  return stringValue(value).replace(/\s+/g, " ").toLowerCase();
}

function appendTuneCodes(target, tuneCodes) {
  if (!Array.isArray(tuneCodes)) return;
  tuneCodes.forEach((code) => {
    if (code && target.tuneCodes.indexOf(code) === -1) {
      target.tuneCodes.push(code);
    }
  });
}

function normalizeUseCase(value) {
  const raw = stringValue(value);
  return useCaseTabs.some((item) => item.value === raw) ? raw : "";
}

function normalizeCarClass(value) {
  const raw = stringValue(value).toUpperCase();
  return classTabs.some((item) => item.value === raw) ? raw : "";
}

function normalizeDrivetrain(value) {
  const raw = stringValue(value).toUpperCase();
  return ["FWD", "RWD", "AWD"].includes(raw) ? raw : "";
}

function normalizeTuneCode(value) {
  const raw = stringValue(value).replace(/\D/g, "");
  return raw.length === 9 ? raw : "";
}

function normalizeTuneCodes(tuneCodes, tuneCode) {
  const rawValues = [];
  if (tuneCode) rawValues.push(tuneCode);
  if (Array.isArray(tuneCodes)) {
    rawValues.push(...tuneCodes);
  } else if (tuneCodes) {
    rawValues.push(tuneCodes);
  }

  const normalized = [];
  rawValues.forEach((value) => {
    const code = normalizeTuneCode(value);
    if (code && normalized.indexOf(code) === -1) {
      normalized.push(code);
    }
  });
  return normalized;
}

function useCaseLabel(value) {
  const found = useCaseTabs.find((item) => item.value === value);
  return found ? found.label : value;
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
  return labels[value] || value;
}

function clampInt(value, min, max, fallback) {
  const parsed = Number(value);
  if (!Number.isFinite(parsed)) return fallback;
  return Math.min(max, Math.max(min, Math.round(parsed)));
}

function stringValue(value) {
  if (value === null || value === undefined) return "";
  return String(value).trim();
}

function firstAvailableClass(useCase, sourceCars) {
  const cars = Array.isArray(sourceCars) ? sourceCars : [];
  const found = classTabs.find((tab) => cars.some((item) => item.useCase === useCase && item.carClass === tab.value));
  return found ? found.value : classTabs[0].value;
}

function classIndex(value) {
  const index = classTabs.findIndex((item) => item.value === value);
  return index === -1 ? 0 : index;
}

function useCaseIndex(value) {
  const index = useCaseTabs.findIndex((item) => item.value === value);
  return index === -1 ? 0 : index;
}

function hasCarsForClass(sourceCars, useCase, carClass) {
  return Array.isArray(sourceCars) && sourceCars.some((item) => item.useCase === useCase && item.carClass === carClass);
}

function buildRecommendShareTitle(data) {
  const prefix = "FH6\u8f66\u8f86\u63a8\u8350";
  const useCase = useCaseLabel(data.activeUseCase || defaultUseCase);
  const carClass = data.activeClass ? `${data.activeClass}\u7ea7` : "";
  return carClass ? `${prefix}\uff1a${useCase} ${carClass}` : prefix;
}

function buildRecommendSharePath(data) {
  return `/pages/recommend/index?${buildRecommendShareQuery(data)}`;
}

function buildRecommendShareQuery(data) {
  const useCase = encodeURIComponent(data.activeUseCase || defaultUseCase);
  const carClass = encodeURIComponent(data.activeClass || defaultClass);
  return `useCase=${useCase}&carClass=${carClass}`;
}

function readCachedRecommendations(fileID) {
  const empty = { version: "", updatedAt: "", cachedAt: 0, stale: false, cars: [] };
  if (!fileID) return empty;
  try {
    const cache = wx.getStorageSync(recommendationsCacheKey);
    if (!cache || cache.fileID !== fileID) return empty;
    if (!Number.isFinite(cache.cachedAt)) return empty;
    return {
      version: stringValue(cache.version),
      updatedAt: stringValue(cache.updatedAt),
      cachedAt: cache.cachedAt,
      stale: Date.now() - cache.cachedAt > recommendationsCacheTtlMs,
      cars: normalizeRecommendedCars(cache.cars),
    };
  } catch (err) {
    console.warn("[recommend] 推荐缓存读取失败。", err);
    return empty;
  }
}

function writeCachedRecommendations(fileID, recommendations) {
  const safeRecommendations = normalizeRecommendationsForCache(recommendations);
  if (!fileID || !safeRecommendations.cars.length) return;
  try {
    wx.setStorageSync(recommendationsCacheKey, {
      fileID,
      version: safeRecommendations.version,
      updatedAt: safeRecommendations.updatedAt,
      cachedAt: Date.now(),
      cars: safeRecommendations.cars,
    });
  } catch (err) {
    console.warn("[recommend] 推荐缓存写入失败。", err);
  }
}

function touchCachedRecommendations(fileID, recommendations) {
  if (!fileID) return;
  try {
    const cache = wx.getStorageSync(recommendationsCacheKey);
    if (!cache || cache.fileID !== fileID || !Array.isArray(cache.cars) || !cache.cars.length) return;
    wx.setStorageSync(recommendationsCacheKey, {
      ...cache,
      version: stringValue(recommendations && recommendations.version) || stringValue(cache.version),
      updatedAt: stringValue(recommendations && recommendations.updatedAt) || stringValue(cache.updatedAt),
      cachedAt: Date.now(),
    });
  } catch (err) {
    console.warn("[recommend] 推荐缓存更新时间失败。", err);
  }
}

function normalizeRecommendationsForCache(recommendations) {
  if (Array.isArray(recommendations)) {
    return {
      version: "",
      updatedAt: "",
      cars: normalizeRecommendedCars(recommendations),
    };
  }
  return {
    version: stringValue(recommendations && recommendations.version),
    updatedAt: stringValue(recommendations && recommendations.updatedAt),
    cars: normalizeRecommendedCars(recommendations && recommendations.cars),
  };
}

function isSameRecommendationsVersion(cachedRecommendations, remoteRecommendations) {
  if (!cachedRecommendations || !cachedRecommendations.cars || !cachedRecommendations.cars.length) return false;
  const remoteVersion = stringValue(remoteRecommendations && remoteRecommendations.version);
  const cachedVersion = stringValue(cachedRecommendations && cachedRecommendations.version);
  if (remoteVersion && cachedVersion) return remoteVersion === cachedVersion;

  const remoteUpdatedAt = stringValue(remoteRecommendations && remoteRecommendations.updatedAt);
  const cachedUpdatedAt = stringValue(cachedRecommendations && cachedRecommendations.updatedAt);
  if (remoteUpdatedAt && cachedUpdatedAt) return remoteUpdatedAt === cachedUpdatedAt;

  return false;
}

function cloudDownloadFile(fileID) {
  return new Promise((resolve, reject) => {
    wx.cloud.downloadFile({
      fileID,
      success: resolve,
      fail: reject,
    });
  });
}

function readFileText(filePath) {
  return new Promise((resolve, reject) => {
    wx.getFileSystemManager().readFile({
      filePath,
      encoding: "utf8",
      success: (res) => resolve(res.data || ""),
      fail: reject,
    });
  });
}
