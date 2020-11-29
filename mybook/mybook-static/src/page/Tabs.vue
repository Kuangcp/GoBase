<template>
  <el-tabs v-model="activeName" @tab-click="handleClick" class="tab-box">
    <el-tab-pane label="记录" name="record">
      <Record />
    </el-tab-pane>
    <el-tab-pane label="明细" name="detail">
      <Detail />
    </el-tab-pane>
    <el-tab-pane label="余额" name="balance">
      <Balance ref="balanceCom" />
    </el-tab-pane>
  </el-tabs>
</template>

<style>
.tab-box {
  padding-left: 36px;
}
</style>

<script>
import Record from "./Record.vue";
import Detail from "./Detail.vue";
import Balance from "./Balance.vue";

export default {
  name: "Tabs",
  components: {
    Record,
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
      if(this.activeName === 'balance'){
        this.$refs.balanceCom.loadBalanceData()
      }
    },
  },
};
</script>
