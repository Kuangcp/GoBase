<template>
  <el-tabs v-model="activeName" @tab-click="handleClick" class="report-tab-box">
    <el-tab-pane label="时序统计" name="dashboard">
      <ReportCategoryTime ref="dashboardCom"/>
    </el-tab-pane>
    <el-tab-pane label="预算" name="budget">
    </el-tab-pane>
  </el-tabs>
</template>

<script>
import ReportCategoryTime from "./ReportCategoryTime.vue";

export default {
  name: "ReportTabs",
  components: {
    ReportCategoryTime,
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