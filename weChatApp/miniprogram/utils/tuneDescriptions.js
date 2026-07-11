const GEAR_DESCRIPTION = "调整各单独传动比会影响加速和极速。较高的传动比可提升加速，而较低的传动比可获得更高的极速。请根据赛道需要的引擎马力与扭矩特性，设置正确的传动比。";

const descriptions = {
  frontTirePressure: "轮胎的理想胎压，取决于使用的踏面胶料。 胎压调整的目的，在于让轮胎获得最大的地面接触面积，同时避免轮胎过热或丧失响应性。较低的胎压可提升接触面积，但若太低则会降低相应速度，同时导致轮胎过热。较高的胎压可提升响应速度，但会导致轮胎温度降低，可能导致车辆瞬间失去抓地力。",
  rearTirePressure: "轮胎的理想胎压，取决于使用的踏面胶料。 胎压调整的目的，在于让轮胎获得最大的地面接触面积，同时避免轮胎过热或丧失响应性。较低的胎压可提升接触面积，但若太低则会降低相应速度，同时导致轮胎过热。较高的胎压可提升响应速度，但会导致轮胎温度降低，可能导致车辆瞬间失去抓地力。",
  finalDrive: "调整最终传动比值（差速器中环形齿轮与行星齿轮的比率）会影响加速和极速，您可藉由设置变速箱中所有齿轮的比率来调整最终传动比。最终传动比值越高，您的传动比越短；这样可提升加速，但会降低极速，或导致您的最终齿轮成为瓶颈。降低最终传动比值可产生较高的极速，但会降低加速性能。",
  gear1: GEAR_DESCRIPTION,
  gear2: GEAR_DESCRIPTION,
  gear3: GEAR_DESCRIPTION,
  gear4: GEAR_DESCRIPTION,
  gear5: GEAR_DESCRIPTION,
  gear6: GEAR_DESCRIPTION,
  gear7: GEAR_DESCRIPTION,
  gear8: GEAR_DESCRIPTION,
  gear9: GEAR_DESCRIPTION,
  gear10: GEAR_DESCRIPTION,
  frontCamber: "无论在直线还是弯道行驶，调整外倾角（轮胎由上到下的倾斜角度）都会影响车辆的抓地力。负值外倾角会令轮胎顶端向内倾斜，提升过弯时轮胎与地面的接触面积。然而，负值外倾角也会降低直线行驶时轮胎与地面的接触面积；负值过大时，会降低直线行驶时的加/减速性能。正值外倾角会令轮胎顶端向外倾斜，降低过弯时轮胎与地面的接触面积，进而降低车辆稳定性。",
  rearCamber: "无论在直线还是弯道行驶，调整外倾角（轮胎由上到下的倾斜角度）都会影响车辆的抓地力。负值外倾角会令轮胎顶端向内倾斜，提升过弯时轮胎与地面的接触面积。然而，负值外倾角也会降低直线行驶时轮胎与地面的接触面积；负值过大时，会降低直线行驶时的加/减速性能。正值外倾角会令轮胎顶端向外倾斜，降低过弯时轮胎与地面的接触面积，进而降低车辆稳定性。",
  frontToe: "调整束角（车轮向内或向外的角度）可让转向响应（直线驾驶与车辆转向之间的转换）变得更加灵敏。内束代表两侧轮胎前端比后端更靠近，这样可提高稳定性，但会降低转向反应速度。外束代表轮胎后端比前端更靠近，这样可提高转向反应速度，但也会降低稳定性。",
  rearToe: "调整束角（车轮向内或向外的角度）可让转向响应（直线驾驶与车辆转向之间的转换）变得更加灵敏。内束代表两侧轮胎前端比后端更靠近，这样可提高稳定性，但会降低转向反应速度。外束代表轮胎后端比前端更靠近，这样可提高转向反应速度，但也会降低稳定性。",
  caster: "调整后倾角（转向轴向前或向后的倾斜角度）可提升直线行驶的稳定性。使用正值后倾角会让转向轴向后倾。当悬挂系统压缩和/或轮胎转向锁定时，外倾角会向负值偏移；因此，增加正值后倾角可让您在负值外倾角较低的情况下驾驶，让轮胎在直线行驶时保持直线（有利于加速和剎车），但在转弯时提供适当的负值外倾角。",
  frontArb: "防倾杆可控制不必要的车体运动，并在稳定过弯时在转向不足与转向过度之间进行平衡。前后防倾杆的硬度会影响这一平衡：要缓解转向不足，可选择将前防倾杆调软、将后防倾杆调硬，或两者同时进行；要缓解转向过度，可选择将后防倾杆调软、将前防倾杆调硬，或两者同时进行。",
  rearArb: "防倾杆可控制不必要的车体运动，并在稳定过弯时在转向不足与转向过度之间进行平衡。前后防倾杆的硬度会影响这一平衡：要缓解转向不足，可选择将前防倾杆调软、将后防倾杆调硬，或两者同时进行；要缓解转向过度，可选择将后防倾杆调软、将前防倾杆调硬，或两者同时进行。",
  frontSpring: "弹簧软硬负责控制车辆在加速、刹车、过弯时的车重转移。较硬的前弹簧可转移较多车重，但太硬则会导致轮胎因负载过重而丧失抓地力。前弹簧比后弹簧软可增加前轮抓地力，并缓解转向不足，但太软则可能导致车底在重踩剎车时托底。前弹簧比后弹簧硬可缓解转向过度，但太硬则可能导致车辆在转弯时发生转向不足。",
  rearSpring: "弹簧软硬负责控制车辆在加速、刹车、过弯时的车重转移。较硬的前弹簧可转移较多车重，但太硬则会导致轮胎因负载过重而丧失抓地力。后弹簧比前弹簧软可增加后轮抓地力，并缓解转向过度。后弹簧比前弹簧硬会增加过度转向。",
  frontRideHeight: "车身高度会决定车辆的离地间隙，以及车体的重心位置。降低车身高度会降低重心，提升过弯性能；然而，高度太低可能导致车身托底而瞬间失控。一般来说，您需要在避免托底的情况下尽可能降低车身高度。",
  rearRideHeight: "车身高度会决定车辆的离地间隙，以及车体的重心位置。降低车身高度会降低重心，提升过弯性能；然而，高度太低可能导致车身托底而瞬间失控。一般来说，您需要在避免触地的情况下尽可能降低车身高度；不过，适当提高车身后部高度，也可以在急加速时帮助控制重量转移。",
  frontRebound: "调整车辆的阻尼可增加抓地力，改善车辆操控性。回弹阻尼可控制悬挂系统回弹远离轮弧内时的伸展速率。调整前轮的回弹阻尼，即可对车辆进出弯道时的平衡进行微调。增加前轮回弹阻尼的硬度会提高过渡性转向不足。降低前轮回弹阻尼的硬度则会提高过渡性转向过度。",
  rearRebound: "调整车辆的阻尼可增加抓地力，改善车辆操控性。回弹阻尼可控制悬挂系统回弹远离轮弧内时的伸展速率。调整后轮的回弹阻尼，即可对车辆进出弯道时的平衡进行微调。增加后轮回弹阻尼的硬度会提高过渡性转向不足。降低后轮回弹阻尼的硬度则会提高过渡性转向过度。",
  frontBump: "调整车辆的阻尼可增加抓地力，改善车辆操控性。压缩阻尼可控制悬挂系统向上进入轮弧内时的压缩速率。增加前轮压缩阻尼的硬度会提高过渡性转向不足，但过度的压缩阻尼则会导致车辆在非平坦路面行驶时出现不稳定情况。降低前轮压缩阻尼的硬度则会提高过渡性转向过度。",
  rearBump: "调整车辆的阻尼可增加抓地力，改善车辆操控性。压缩阻尼可控制悬挂系统向上进入轮弧内时的压缩速率。增加后轮压缩阻尼的硬度会提高过渡性转向不足，但过度的压缩阻尼则会导致车辆在非平坦路面行驶时出现不稳定情况。降低后轮压缩阻尼的硬度则会提高过渡性转向过度。",
  frontAero: "提升下压力，可让车辆与路面间保持更好的接触、更快热胎，同时改善高速行驶时的操控性。“效率”与车辆产生风阻的程度有关： 数值越高，车辆效率越高，产生的风阻越小。提升下压力也会导致车辆产生的风阻变大。“平衡“指车辆的空气动力学平衡：车身前部与后部之间的下压力差异越大，车辆空气动力学属性的变化也就越大。平衡值越低，车辆展现出的转向不足越严重；平衡值越高，车辆展现出的转向过度越严重。",
  rearAero: "提升下压力，可让车辆与路面间保持更好的接触、更快热胎，同时改善高速行驶时的操控性。“效率”与车辆产生风阻的程度有关： 数值越高，车辆效率越高，产生的风阻越小。提升下压力也会导致车辆产生的风阻变大。“平衡“指车辆的空气动力学平衡：车身前部与后部之间的下压力差异越大，车辆空气动力学属性的变化也就越大。平衡值越低，车辆展现出的转向不足越严重；平衡值越高，车辆展现出的转向过度越严重。",
  brakeBalance: "刹车平衡会影响制动力的分配比例，进而影响刹车距离，以及刹车时的转向不足/转向过度平衡。将刹车平衡向后调整会提升刹车时的转向过度，将刹车平衡向前调整则会提升转向不足、同时提升稳定性，但也可能导致刹车时出现严重的转向不足。",
  brakePressure: "刹车油压会根据踩下刹车的程度影响产生的制动力。降低整体刹车油压，可提升“产生显著制动力”所需的踏板行程；然而，若降低太多，车辆便无法有效减速。提升刹车油压，可降低“产生显著制动力”所需的踏板行程；然而，若提升太多，刹车便容易立即抱死。",
  frontDiffAccel: "加速差速器设置可调整车辆在加速时，触发速器锁定所需的轮胎转动差异值。提高加速设置会提升差速器在加速时的锁定速度。 条低加速设置会减慢差速器的锁定速度。降低前轮差速器的加速设置可缓解转向不足，但降低太多则会影响车辆的响应速度。",
  frontDiffDecel: "减速差速器设置可调整车辆在減速时，触发差速器锁定所需的轮胎转动差异值。提高减速设置会提升差速器在减速时的锁定速度，但差速器锁定过快会影响车辆的操控性。降低前轮的减速设置可缓解松开油门时的转向过度，但会提升前轮刹车抱死的频率（在未启动 ABS 系统的情况下）。",
  rearDiffAccel: "加速差速器设置可调整车辆在加速时，触发差速器锁定所需的轮胎转动差异值。提高加速设置会提升差速器在加速时的锁定速度。 提高后轮差速器的加速设置会加剧后轮驱动或四轮驱动车辆的转向过度。高马力车辆则必须提升加速设置，以维持足够的抓地力，但需要注意，锁定太快会导致操控性降低。 降低加速设置会减慢差速器的锁定速度。",
  rearDiffDecel: "减速差速器设置可调整车辆在减速时，触发差速器锁定所需的轮胎转动差异值。提高减速设置会提升差速器在减速时的锁定速度，但差速器锁定过快会影响车辆的操控性。提高后轮的减速设置可缓解松开油门时的转向过度。",
  centerDiffBalance: "中央差速器负责控制四轮驱动车辆前后轮轴之间的驱动扭矩分配。增加后轮扭矩可将更多车辆动力输送至后轮，提升响应速度与转向过度。后轮扭矩过高可能导致轮胎空转和大量转向过度。增加前轮扭矩可让车辆踩下油门时的转向过度表现更接近前轮驱动车辆，使车辆更加稳定。前轮扭矩过高可能导致车辆出现严重的转向不足。",
};

const editableFields = {
  frontTirePressure: true,
  rearTirePressure: true,
  finalDrive: true,
  gear1: true,
  gear2: true,
  gear3: true,
  gear4: true,
  gear5: true,
  gear6: true,
  gear7: true,
  gear8: true,
  gear9: true,
  gear10: true,
  frontCamber: true,
  rearCamber: true,
  frontToe: true,
  rearToe: true,
  caster: true,
  frontArb: true,
  rearArb: true,
  frontSpring: true,
  rearSpring: true,
  frontRebound: true,
  rearRebound: true,
  frontBump: true,
  rearBump: true,
  frontAero: true,
  rearAero: true,
  brakeBalance: true,
  brakePressure: true,
  frontDiffAccel: true,
  frontDiffDecel: true,
  rearDiffAccel: true,
  rearDiffDecel: true,
  centerDiffBalance: true,
};

function decorateResultGroups(groups) {
  if (!Array.isArray(groups)) return [];
  return groups.map((group, groupIndex) => {
    const safeGroup = group || {};
    const groupKey = safeGroup.key || safeGroup.group || `group${groupIndex}`;
    const items = Array.isArray(safeGroup.items) ? safeGroup.items : [];
    return {
      ...safeGroup,
      key: groupKey,
      items: items.map((field, fieldIndex) => decorateField(field, groupKey, fieldIndex)),
    };
  });
}

function decorateField(field, groupKey, fieldIndex) {
  const safeField = field || {};
  const fieldKey = safeField.fieldKey || "";
  return {
    ...safeField,
    rowId: safeField.rowId || `${groupKey}-${fieldKey || "field"}-${fieldIndex}-${safeField.tierItem ? "tier" : "value"}`,
    description: descriptions[fieldKey] || "",
    editable: isEditableTuneField(safeField),
    manualEdited: Boolean(safeField.manualEdited),
  };
}

function isEditableTuneField(field) {
  if (!field || field.tierItem) return false;
  if (!editableFields[field.fieldKey]) return false;
  if (field.displayValue === "--" || field.displayValue === "请开启齿轮调校") return false;
  return true;
}

function parseManualDisplayValue(fieldKey, rawValue) {
  const raw = String(rawValue === undefined || rawValue === null ? "" : rawValue).trim();
  if (!/^-?\d+(\.\d+)?$/.test(raw)) {
    throw new Error("请输入数字");
  }
  const numericValue = Number(raw);
  if (!Number.isFinite(numericValue)) {
    throw new Error("请输入数字");
  }
  const step = fieldStep(fieldKey);
  const decimals = fieldDecimals(fieldKey);
  const rounded = Math.round(numericValue / step) * step;
  const value = decimals === 0 ? Math.round(rounded) : Number(rounded.toFixed(decimals));
  return {
    value,
    displayValue: formatManualValue(fieldKey, value),
  };
}

function updateGroupsFieldValue(groups, fieldKey, parsedValue) {
  const nextGroups = decorateResultGroups(groups).map((group) => ({
    ...group,
    items: group.items.map((field) => {
      if (field.fieldKey !== fieldKey || !field.editable) return field;
      return {
        ...field,
        value: parsedValue.value,
        displayValue: parsedValue.displayValue,
        manualEdited: true,
      };
    }),
  }));
  return decorateResultGroups(nextGroups);
}

function updateResultFieldValue(result, fieldKey, value) {
  if (!result || typeof result !== "object") return result;
  const nextResult = {
    ...result,
    profileDraft: {
      ...(result.profileDraft || {}),
      [fieldKey]: value,
    },
    generatedFields: Array.isArray(result.generatedFields)
      ? result.generatedFields.map((field) => (field.fieldKey === fieldKey ? { ...field, value, manualEdited: true } : field))
      : result.generatedFields,
  };
  return nextResult;
}

function findResultField(groups, fieldKey) {
  const decoratedGroups = decorateResultGroups(groups);
  for (let groupIndex = 0; groupIndex < decoratedGroups.length; groupIndex += 1) {
    const group = decoratedGroups[groupIndex];
    const field = group.items.find((item) => item.fieldKey === fieldKey && item.editable);
    if (field) return field;
  }
  return null;
}

function formatManualValue(fieldKey, value) {
  const decimals = fieldDecimals(fieldKey);
  return Number(value).toFixed(decimals);
}

function fieldStep(fieldKey) {
  if (fieldKey === "frontTirePressure" || fieldKey === "rearTirePressure") return 0.01;
  if (fieldKey === "finalDrive" || /^gear\d+$/.test(fieldKey)) return 0.01;
  if (fieldKey === "frontArb" || fieldKey === "rearArb") return 0.1;
  if (["frontAero", "rearAero", "brakeBalance", "brakePressure", "frontDiffAccel", "frontDiffDecel", "rearDiffAccel", "rearDiffDecel", "centerDiffBalance"].includes(fieldKey)) return 1;
  return 0.1;
}

function fieldDecimals(fieldKey) {
  const step = fieldStep(fieldKey);
  if (step >= 1) return 0;
  if (step >= 0.1) return 1;
  return 2;
}

module.exports = {
  decorateResultGroups,
  findResultField,
  parseManualDisplayValue,
  updateGroupsFieldValue,
  updateResultFieldValue,
};
