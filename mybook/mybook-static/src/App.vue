<template>
  <el-tabs v-model="activeName" @tab-click="handleClick" class="tab-box">
    <el-tab-pane label="报表" name="report">
      <ReportCategoryTime/>
    </el-tab-pane>
    <el-tab-pane label="记账" name="record">
      <Record/>
    </el-tab-pane>
    <el-tab-pane label="明细" name="detail">
      <Detail/>
      <!--      <DetailTabs/>-->
    </el-tab-pane>
    <el-tab-pane label="分类" name="category">
      <Category ref="categoryCom"/>
    </el-tab-pane>
    <el-tab-pane label="财务状况" name="balance">
      <Balance ref="balanceCom"/>
    </el-tab-pane>
    <el-tab-pane label="设置" name="setting">
      <Setting/>
    </el-tab-pane>
  </el-tabs>
</template>

<script>

import Record from "./page/Record.vue";
import Detail from "./page/Detail.vue";
import Balance from "./page/Balance.vue";
import Category from "./page/Category.vue";
import Setting from "./page/Setting";
// import ReportTabs from "./page/ReportTabs.vue";
// import DetailTabs from "./page/DetailTabs";
import ReportCategoryTime from "./page/ReportCategoryTime.vue";

export default {
  name: "App",
  components: {
    // DetailTabs,
    Record,
    Detail,
    Balance,
    Category,
    Setting,
    // ReportTabs,
    ReportCategoryTime,
  },
  data() {
    return {
      activeName: "report",
    };
  },
  methods: {
    handleClick() {
      if (this.activeName === "balance") {
        this.$refs.balanceCom.loadBalanceData();
      }
    },
  },
};

</script>

<style>
* {
  margin: 0;
  padding: 0;
}

.tab-box {
  padding-left: 22px;
}

body {
  height: 100vh;
  /*https://pxlab.cn/color/index.html*/
  /*background-image: linear-gradient(25deg, #65b8b5, #96b4ae, #bbb0a7, #dcaaa0);*/
  /*background-image: linear-gradient(25deg, #227495, #76859d, #b095a5, #e7a5ad)*/

  /*background-image: linear-gradient(60deg, #29323c 0%, #485563 100%);*/

  background: linear-gradient(-180deg, #BCC5CE 0%, #929EAD 98%), radial-gradient(at top left, rgba(255, 255, 255, 0.30) 0%, rgba(0, 0, 0, 0.30) 100%);
  background-blend-mode: screen;

  /*background: linear-gradient(to bottom, #323232 0%, #3F3F3F 40%, #1C1C1C 150%), linear-gradient(to top, rgba(255, 255, 255, 0.40) 0%, rgba(0, 0, 0, 0.25) 200%);*/
  /*background-blend-mode: multiply;*/

}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  /* text-align: center; */
  color: #2c3e50;
  padding-top: 8px;
}
</style>
