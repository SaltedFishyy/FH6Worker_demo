const HISTORY_KEY = "fh6TuneHistoryV1";
const MAX_HISTORY_COUNT = 50;

function loadTuneHistory() {
  try {
    const raw = wx.getStorageSync(HISTORY_KEY);
    if (!Array.isArray(raw)) return [];
    return raw.map(normalizeRecord).filter(Boolean);
  } catch (err) {
    return [];
  }
}

function addTuneRecord(input) {
  const record = createTuneRecord(input);
  if (!record) return null;
  const nextRecords = [record, ...loadTuneHistory()].slice(0, MAX_HISTORY_COUNT);
  saveTuneHistory(nextRecords);
  return record;
}

function updateTuneRecord(id, patch) {
  if (!id || !patch || typeof patch !== "object") return null;
  let updatedRecord = null;
  const nextRecords = loadTuneHistory().map((record) => {
    if (record.id !== id) return record;
    updatedRecord = normalizeRecord({
      ...record,
      ...patch,
      summary: {
        ...(record.summary || {}),
        ...(patch.summary || {}),
      },
      updatedAt: Date.now(),
    });
    return updatedRecord;
  });
  if (updatedRecord) {
    saveTuneHistory(nextRecords);
  }
  return updatedRecord;
}

function removeTuneRecord(id) {
  const nextRecords = loadTuneHistory().filter((item) => item.id !== id);
  saveTuneHistory(nextRecords);
  return nextRecords;
}

function clearTuneHistory() {
  saveTuneHistory([]);
}

function saveTuneHistory(records) {
  try {
    wx.setStorageSync(HISTORY_KEY, records);
  } catch (err) {
    // History is local convenience data; storage failure should not block tuning.
  }
}

function createTuneRecord(input) {
  if (!input || !input.payload || !input.result) return null;
  const createdAt = Date.now();
  const profile = input.result.profileDraft || {};
  const payload = input.payload;
  const summary = {
    useCase: payload.useCase || profile.useCase || "",
    useCaseLabel: input.useCaseLabel || useCaseLabel(payload.useCase || profile.useCase),
    carClass: profile.carClass || classFromPi(payload.pi),
    pi: Number(profile.pi || payload.pi || 0),
    drivetrain: profile.drivetrain || payload.drivetrain || "",
    tireCompoundLabel: input.tireCompoundLabel || tireCompoundLabel(payload.tireCompound),
  };
  const displayName = normalizeDisplayName(input.displayName || summary.displayName || defaultDisplayName(summary));
  summary.displayName = displayName;
  return {
    id: `${createdAt}-${Math.random().toString(36).slice(2, 8)}`,
    createdAt,
    createdAtText: formatTime(createdAt),
    displayName,
    summary,
    payload,
    result: input.result,
    resultGroups: Array.isArray(input.resultGroups) ? input.resultGroups : [],
    warnings: Array.isArray(input.warnings) ? input.warnings : [],
    nextTestPlan: Array.isArray(input.nextTestPlan) ? input.nextTestPlan : [],
  };
}

function normalizeRecord(record) {
  if (!record || typeof record !== "object" || !record.id) return null;
  const summary = record.summary || {};
  const displayName = normalizeDisplayName(record.displayName || summary.displayName || defaultDisplayName(summary));
  return {
    ...record,
    displayName,
    createdAtText: record.createdAtText || formatTime(record.createdAt),
    summary: {
      ...summary,
      displayName,
    },
    resultGroups: Array.isArray(record.resultGroups) ? record.resultGroups : [],
    warnings: Array.isArray(record.warnings) ? record.warnings : [],
    nextTestPlan: Array.isArray(record.nextTestPlan) ? record.nextTestPlan : [],
  };
}

function normalizeDisplayName(value) {
  const name = String(value || "").trim();
  return name.slice(0, 24);
}

function defaultDisplayName(summary) {
  const useCase = summary && summary.useCaseLabel ? summary.useCaseLabel : "车辆";
  return `${useCase} 调校`;
}

function formatTime(value) {
  const date = new Date(Number(value) || Date.now());
  const pad = (num) => String(num).padStart(2, "0");
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`;
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

module.exports = {
  addTuneRecord,
  clearTuneHistory,
  loadTuneHistory,
  removeTuneRecord,
  updateTuneRecord,
};
