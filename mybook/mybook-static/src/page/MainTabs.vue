<template>
  <el-tabs v-model="activeName" @tab-click="handleClick" class="tab-box">
    <el-tab-pane label="面板" name="report">
      <ReportTabs/>
    </el-tab-pane>
    <el-tab-pane label="记录" name="record">
      <Record/>
    </el-tab-pane>
    <el-tab-pane label="明细" name="detail">
      <Detail/>
    </el-tab-pane>
    <el-tab-pane label="分类" name="category">
      <Category ref="categoryCom"/>
    </el-tab-pane>
    <el-tab-pane label="余额" name="balance">
      <Balance ref="balanceCom"/>
    </el-tab-pane>
  </el-tabs>
</template>

<style scoped>
.tab-box {
  padding-left: 22px;
}
</style>

<script>
import Record from "./Record.vue";
import Detail from "./Detail.vue";
import Balance from "./Balance.vue";
import Category from "./Category.vue";
import ReportTabs from "./ReportTabs.vue";

export default {
  name: "MainTabs",
  components: {
    Record,
    Detail,
    Balance,
    Category,
    ReportTabs,
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
