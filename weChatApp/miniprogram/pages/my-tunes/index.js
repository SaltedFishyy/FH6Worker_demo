const { setSelectedTab, setTabBarHidden } = require("../../utils/tabBar");
const { clearTuneHistory, loadTuneHistory, removeTuneRecord, updateTuneRecord } = require("../../utils/tuneHistory");
const {
  decorateResultGroups,
  findResultField,
  parseManualDisplayValue,
  updateGroupsFieldValue,
  updateResultFieldValue,
} = require("../../utils/tuneDescriptions");
const { buildSharePath, createTuneShare, shareTitle } = require("../../utils/tuneShare");

Page({
  data: {
    records: [],
    expandedRecordIds: {},
    editingRecordId: "",
    editingFieldKey: "",
    editingValue: "",
    shareLinks: {},
    shareLoadingRecordId: "",
  },

  onShow() {
    setSelectedTab(this, 1);
    setTabBarHidden(this, false);
    this.refreshRecords();
  },

  refreshRecords() {
    this.setData({
      records: applyRecordState(loadTuneHistory(), this.data.expandedRecordIds, this.data.shareLinks, this.data.shareLoadingRecordId),
    });
  },

  toggleRecord(e) {
    const id = e.currentTarget.dataset.id;
    if (!id) return;
    const expandedRecordIds = { ...this.data.expandedRecordIds };
    if (expandedRecordIds[id]) {
      delete expandedRecordIds[id];
    } else {
      expandedRecordIds[id] = true;
    }
    this.setData({
      expandedRecordIds,
      records: applyRecordState(this.data.records, expandedRecordIds, this.data.shareLinks, this.data.shareLoadingRecordId),
    });
  },

  deleteRecord(e) {
    const id = e.currentTarget.dataset.id;
    if (!id) return;
    wx.showModal({
      title: "删除调校",
      content: "删除后仅从本机移除，不影响已生成的当前结果。",
      confirmText: "删除",
      confirmColor: "#ef233c",
      success: (res) => {
        if (!res.confirm) return;
        const expandedRecordIds = { ...this.data.expandedRecordIds };
        const shareLinks = { ...this.data.shareLinks };
        delete expandedRecordIds[id];
        delete shareLinks[id];
        const records = applyRecordState(removeTuneRecord(id), expandedRecordIds, shareLinks, this.data.shareLoadingRecordId);
        this.setData({ expandedRecordIds, shareLinks, records });
      },
    });
  },

  renameRecord(e) {
    const id = e.currentTarget.dataset.id;
    if (!id) return;
    const record = this.data.records.find((item) => item.id === id);
    if (!record) return;
    const currentName = record.displayName || (record.summary && record.summary.displayName) || `${record.summary.useCaseLabel || "车辆"} 调校`;
    wx.showModal({
      title: "修改名称",
      editable: true,
      placeholderText: "最多 24 个字",
      content: currentName,
      confirmText: "保存",
      success: (res) => {
        if (!res.confirm) return;
        const displayName = normalizeDisplayName(res.content);
        if (!displayName) {
          wx.showToast({
            title: "名称不能为空",
            icon: "none",
          });
          return;
        }
        const updated = updateTuneRecord(id, {
          displayName,
          summary: {
            ...record.summary,
            displayName,
          },
        });
        if (!updated) return;
        const shareLinks = { ...this.data.shareLinks };
        delete shareLinks[id];
        const records = this.data.records.map((item) => (
          item.id === id
            ? {
                ...item,
                displayName,
                summary: {
                  ...(item.summary || {}),
                  displayName,
                },
              }
            : item
        ));
        this.setData({
          records: applyRecordState(records, this.data.expandedRecordIds, shareLinks, this.data.shareLoadingRecordId),
          shareLinks,
        });
      },
    });
  },

  onRecordValueTap(e) {
    const recordId = e.currentTarget.dataset.id;
    const fieldKey = e.currentTarget.dataset.fieldKey;
    const editable = e.currentTarget.dataset.editable;
    if (!recordId || !fieldKey || !(editable === true || editable === "true")) return;

    const now = Date.now();
    const lastTap = this.lastValueTap || {};
    this.lastValueTap = { recordId, fieldKey, time: now };
    if (lastTap.recordId !== recordId || lastTap.fieldKey !== fieldKey || now - lastTap.time > 360) return;

    const record = this.data.records.find((item) => item.id === recordId);
    const field = record ? findResultField(record.resultGroups, fieldKey) : null;
    if (!field) return;
    this.setData({
      editingRecordId: recordId,
      editingFieldKey: fieldKey,
      editingValue: field.displayValue,
    });
  },

  onEditInput(e) {
    this.setData({
      editingValue: e.detail.value,
    });
  },

  finishRecordFieldEdit(e) {
    const recordId = e.currentTarget.dataset.id || this.data.editingRecordId;
    const fieldKey = e.currentTarget.dataset.fieldKey || this.data.editingFieldKey;
    if (!recordId || !fieldKey || this.data.editingRecordId !== recordId || this.data.editingFieldKey !== fieldKey) return;

    let parsed;
    try {
      parsed = parseManualDisplayValue(fieldKey, this.data.editingValue);
    } catch (err) {
      wx.showToast({
        title: err.message || "请输入数字",
        icon: "none",
      });
      this.setData({
        editingRecordId: "",
        editingFieldKey: "",
        editingValue: "",
      });
      return;
    }

    const record = this.data.records.find((item) => item.id === recordId);
    if (!record) return;
    const resultGroups = updateGroupsFieldValue(record.resultGroups, fieldKey, parsed);
    const result = updateResultFieldValue(record.result, fieldKey, parsed.value);
    updateTuneRecord(recordId, {
      result,
      resultGroups,
    });

    const shareLinks = { ...this.data.shareLinks };
    delete shareLinks[recordId];
    const records = this.data.records.map((item) => (
      item.id === recordId
        ? { ...item, result, resultGroups }
        : item
    ));
    this.setData({
      records: applyRecordState(records, this.data.expandedRecordIds, shareLinks, this.data.shareLoadingRecordId),
      shareLinks,
      editingRecordId: "",
      editingFieldKey: "",
      editingValue: "",
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

  async prepareRecordShare(e) {
    if (this.data.shareLoadingRecordId) return;
    const id = e.currentTarget.dataset.id;
    const record = this.data.records.find((item) => item.id === id);
    if (!record) return;

    this.setData({
      shareLoadingRecordId: id,
      records: applyRecordState(this.data.records, this.data.expandedRecordIds, this.data.shareLinks, id),
    });
    try {
      const data = await createTuneShare({
        summary: {
          ...(record.summary || {}),
          displayName: record.displayName || (record.summary && record.summary.displayName) || "",
        },
        result: record.result || {},
        resultGroups: record.resultGroups || [],
        warnings: record.warnings || [],
        nextTestPlan: record.nextTestPlan || [],
      });
      const shareLinks = {
        ...this.data.shareLinks,
        [id]: {
          shareId: data.shareId,
          sharePath: buildSharePath(data.shareId),
        },
      };
      this.setData({
        shareLinks,
        shareLoadingRecordId: "",
        records: applyRecordState(this.data.records, this.data.expandedRecordIds, shareLinks, ""),
      });
      wx.showToast({
        title: "链接已生成",
        icon: "success",
      });
    } catch (err) {
      this.setData({
        shareLoadingRecordId: "",
        records: applyRecordState(this.data.records, this.data.expandedRecordIds, this.data.shareLinks, ""),
      });
      wx.showModal({
        title: "分享失败",
        content: err && err.message ? err.message : "分享链接生成失败，请稍后再试。",
        showCancel: false,
      });
    }
  },

  onShareAppMessage(e) {
    const id = e && e.target && e.target.dataset ? e.target.dataset.id : "";
    const record = this.data.records.find((item) => item.id === id);
    const link = id ? this.data.shareLinks[id] : null;
    return {
      title: shareTitle(record ? record.summary : {}),
      path: link && link.sharePath ? link.sharePath : "/pages/my-tunes/index",
    };
  },

  clearAll() {
    if (!this.data.records.length) return;
    wx.showModal({
      title: "清空我的调校",
      content: "将删除本机保存的全部调校记录。",
      confirmText: "清空",
      confirmColor: "#ef233c",
      success: (res) => {
        if (!res.confirm) return;
        clearTuneHistory();
        this.setData({
          records: [],
          expandedRecordIds: {},
          shareLinks: {},
          shareLoadingRecordId: "",
        });
      },
    });
  },
});

function applyRecordState(records, expandedRecordIds, shareLinks, shareLoadingRecordId) {
  return records.map((record) => {
    const shareLink = shareLinks && shareLinks[record.id] ? shareLinks[record.id] : null;
    return {
      ...record,
      resultGroups: decorateResultGroups(record.resultGroups),
      expanded: Boolean(expandedRecordIds && expandedRecordIds[record.id]),
      shareReady: Boolean(shareLink),
      shareLoading: shareLoadingRecordId === record.id,
    };
  });
}

function normalizeDisplayName(value) {
  return String(value || "").trim().slice(0, 24);
}
