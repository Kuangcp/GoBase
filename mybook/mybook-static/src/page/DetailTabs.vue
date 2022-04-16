<template>
  <el-tabs v-model="activeName" class="report-tab-box">
    <el-tab-pane label="明细" name="details">
      <Detail/>
    </el-tab-pane>
    <el-tab-pane label="分类" name="category">
      <Category/>
    </el-tab-pane>
  </el-tabs>
</template>

<script>
import Detail from "./Detail.vue";
import Category from "./Category";

export default {
  name: "DetailTabs",
  components: {
    Category,
    Detail,
  },
  data() {
    return {
      activeName: "details",
    }
  }
}
</script>

<style scoped>

</style>