<template>
  <el-tabs v-model="activeName" @tab-click="handleClick">
    <el-tab-pane label="记录" name="record"> </el-tab-pane>
    <el-tab-pane label="明细" name="detail">
      <Detail />
    </el-tab-pane>
  </el-tabs>
</template>

<script>
import Detail from "./Detail.vue";

export default {
  name: "Tabs",
  components: {
    Detail,
  },
  data() {
    return {
      activeName: "detail",
    };
  },
};
</script>
