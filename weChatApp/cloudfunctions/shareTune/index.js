const crypto = require("crypto");

let cloud = null;

try {
  cloud = require("wx-server-sdk");
  cloud.init({
    env: cloud.DYNAMIC_CURRENT_ENV,
  });
} catch (err) {
  cloud = null;
}

const COLLECTION = "tune_shares";
const EXPIRE_MS = 30 * 24 * 60 * 60 * 1000;

exports.main = async (event = {}) => {
  try {
    if (!cloud) {
      throw new AppError("CLOUD_UNAVAILABLE", "云函数运行环境不可用。");
    }

    const type = event.type || "";
    if (type === "create") {
      return success(await createShare(event.payload));
    }
    if (type === "get") {
      return success(await getShare(event.shareId));
    }
    throw new AppError("VALIDATION_ERROR", "未知分享操作。");
  } catch (err) {
    const normalizedError = normalizeError(err);
    return {
      success: false,
      errorCode: normalizedError.errorCode,
      message: normalizedError.message,
    };
  }
};

async function createShare(payload) {
  validatePayload(payload);
  const db = cloud.database();
  const shareId = createShareId();
  const now = Date.now();
  const expiresAt = now + EXPIRE_MS;

  await ensureCollection(db);
  await cleanupExpiredShares(db, now);
  await setShareDocument(db, shareId, {
    data: {
      createdAt: now,
      expiresAt,
      snapshot: sanitizePayload(payload),
    },
  });

  return {
    shareId,
    expiresAt,
  };
}

async function getShare(shareId) {
  const id = String(shareId || "").trim();
  if (!/^[a-z0-9]{12,40}$/.test(id)) {
    throw new AppError("VALIDATION_ERROR", "分享链接参数无效。");
  }

  const db = cloud.database();
  await ensureCollection(db);
  let doc;
  try {
    doc = await db.collection(COLLECTION).doc(id).get();
  } catch (err) {
    throw new AppError("NOT_FOUND", "分享内容不存在或已被删除。");
  }

  const data = doc && doc.data ? doc.data : null;
  if (!data || !data.snapshot) {
    throw new AppError("NOT_FOUND", "分享内容不存在或已被删除。");
  }
  if (Number(data.expiresAt) <= Date.now()) {
    await db.collection(COLLECTION).doc(id).remove().catch(() => {});
    throw new AppError("EXPIRED", "分享链接已过期，请让好友重新分享。");
  }

  return {
    ...data.snapshot,
    createdAt: data.createdAt,
    expiresAt: data.expiresAt,
  };
}

async function setShareDocument(db, shareId, data) {
  try {
    await db.collection(COLLECTION).doc(shareId).set(data);
  } catch (err) {
    if (!isCollectionNotExist(err)) throw err;
    await ensureCollection(db);
    await db.collection(COLLECTION).doc(shareId).set(data);
  }
}

async function ensureCollection(db) {
  try {
    await db.createCollection(COLLECTION);
  } catch (err) {
    if (isCollectionAlreadyExists(err)) return;
    if (isCollectionNotExist(err)) return;
    const message = errorMessage(err);
    if (message.indexOf("already exists") !== -1 || message.indexOf("collection already exists") !== -1) return;
    throw err;
  }
}

function validatePayload(payload) {
  if (!payload || typeof payload !== "object") {
    throw new AppError("VALIDATION_ERROR", "分享内容为空。");
  }
  if (!payload.summary || typeof payload.summary !== "object") {
    throw new AppError("VALIDATION_ERROR", "分享摘要缺失。");
  }
  if (!Array.isArray(payload.resultGroups) || payload.resultGroups.length === 0) {
    throw new AppError("VALIDATION_ERROR", "调校结果为空。");
  }
  const size = Buffer.byteLength(JSON.stringify(payload), "utf8");
  if (size > 800 * 1024) {
    throw new AppError("VALIDATION_ERROR", "分享内容过大。");
  }
}

function sanitizePayload(payload) {
  return JSON.parse(JSON.stringify({
    summary: payload.summary || {},
    result: payload.result || {},
    resultGroups: payload.resultGroups || [],
    warnings: Array.isArray(payload.warnings) ? payload.warnings : [],
    nextTestPlan: Array.isArray(payload.nextTestPlan) ? payload.nextTestPlan : [],
  }));
}

function createShareId() {
  return `${Date.now().toString(36)}${crypto.randomBytes(8).toString("hex")}`;
}

async function cleanupExpiredShares(db, now) {
  try {
    const _ = db.command;
    await db.collection(COLLECTION).where({
      expiresAt: _.lte(now),
    }).remove();
  } catch (err) {
    // Expired snapshot cleanup must not block creating a new share.
  }
}

function normalizeError(err) {
  if (err && err.errorCode) {
    return err;
  }
  if (isCollectionNotExist(err)) {
    return new AppError("DATABASE_COLLECTION_NOT_EXIST", "分享数据集合 tune_shares 不存在，请在云开发数据库中创建，或重新部署 shareTune 后重试。");
  }
  return new AppError("VALIDATION_ERROR", err && err.message ? err.message : "分享服务暂不可用。");
}

function isCollectionNotExist(err) {
  const message = errorMessage(err);
  return err && (
    err.errCode === -502005
    || err.code === "DATABASE_COLLECTION_NOT_EXIST"
    || message.indexOf("DATABASE_COLLECTION_NOT_EXIST") !== -1
    || message.indexOf("database collection not exists") !== -1
    || message.indexOf("Db or Table not exist") !== -1
    || message.indexOf("ResourceNotFound") !== -1
  );
}

function isCollectionAlreadyExists(err) {
  const message = errorMessage(err);
  return err && (
    err.errCode === -501001
    || err.code === "DATABASE_COLLECTION_ALREADY_EXIST"
    || message.indexOf("DATABASE_COLLECTION_ALREADY_EXIST") !== -1
    || message.indexOf("already exist") !== -1
  );
}

function errorMessage(err) {
  if (!err) return "";
  return String(err.errMsg || err.message || err.toString() || "");
}

function success(data) {
  return {
    success: true,
    data,
  };
}

function AppError(errorCode, message) {
  this.name = "AppError";
  this.errorCode = errorCode;
  this.message = message;
}

AppError.prototype = Object.create(Error.prototype);
AppError.prototype.constructor = AppError;
