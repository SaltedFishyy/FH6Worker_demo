function setSelectedTab(page, selected) {
  if (!page || typeof page.getTabBar !== "function") return;
  const tabBar = page.getTabBar();
  if (!tabBar) return;
  tabBar.setData({ selected });
}

function setTabBarHidden(page, hidden) {
  if (!page || typeof page.getTabBar !== "function") return;
  const tabBar = page.getTabBar();
  if (!tabBar) return;
  tabBar.setData({ hidden: Boolean(hidden) });
}

module.exports = {
  setSelectedTab,
  setTabBarHidden,
};
