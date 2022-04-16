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
// import ReportTabs from "./ReportTabs.vue";
// import DetailTabs from "./DetailTabs";
import ReportCategoryTime from "./ReportCategoryTime.vue";

export default {
  name: "MainTabs",
  components: {
    // DetailTabs,
    Record,
    Detail,
    Balance,
    Category,
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
