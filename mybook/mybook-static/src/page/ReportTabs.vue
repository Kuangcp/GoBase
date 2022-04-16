<template>
  <el-tabs v-model="activeName" @tab-click="handleClick" class="report-tab-box">
    <el-tab-pane label="时序统计" name="dashboard">
      <CategoryTimeReport ref="dashboardCom"/>
    </el-tab-pane>
  </el-tabs>
</template>

<script>
import CategoryTimeReport from "./CategoryTimeReport.vue";

export default {
  name: "ReportTabs",
  components: {
    CategoryTimeReport,
  },
  data() {
    return {
      activeName: "dashboard",
    };
  },
  methods: {
    handleClick() {

    },
  },
}
</script>

<style scoped>
.report-tab-box {
  padding-left: 6px;
}
</style>