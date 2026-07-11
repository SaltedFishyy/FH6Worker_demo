Component({
  data: {
    hidden: false,
    selected: 0,
    list: [
      {
        pagePath: "/pages/index/index",
        text: "快速调校",
        iconPath: "/images/icons/tune.png",
        selectedIconPath: "/images/icons/tune-active.png",
      },
      {
        pagePath: "/pages/my-tunes/index",
        text: "我的调校",
        iconPath: "/images/icons/garage.png",
        selectedIconPath: "/images/icons/garage-active.png",
      },
      {
        pagePath: "/pages/recommend/index",
        text: "车辆推荐",
        iconPath: "/images/icons/car.png",
        selectedIconPath: "/images/icons/car-active.png",
      },
    ],
  },

  methods: {
    switchTab(e) {
      const index = Number(e.currentTarget.dataset.index);
      const item = this.data.list[index];
      if (!item) return;
      wx.switchTab({
        url: item.pagePath,
      });
    },
  },
});
