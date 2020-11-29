<template>
  <el-tabs v-model="activeName" @tab-click="handleClick" class="tab-box">
    <el-tab-pane label="记录" name="record"> </el-tab-pane>
    <el-tab-pane label="明细" name="detail">
      <Detail />
    </el-tab-pane>
    <el-tab-pane label="余额" name="balance">
      <Balance />
    </el-tab-pane>
  </el-tabs>
</template>

<style>
.tab-box {
  padding-left: 36px;
}
</style>

<script>
import Detail from "./Detail.vue";
import Balance from "./Balance.vue";

export default {
  name: "Tabs",
  components: {
    Detail,
    Balance,
  },
  data() {
    return {
      activeName: "detail",
    };
  },
  methods: {
    handleClick() {
      console.log(this.activeName)
    },
  },
};
</script>
