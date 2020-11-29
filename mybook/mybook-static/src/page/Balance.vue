<template>
  <div>
    余额：{{ totalAmount }}
    <el-table :data="tableData" stripe style="width: 100%" height="800">
      <el-table-column sortable prop="ID" label="ID" width="60" align="right">
      </el-table-column>
      <el-table-column
        sortable
        prop="Name"
        label="账户"
        width="120"
        align="center"
      >
      </el-table-column>
      <el-table-column
        sortable
        prop="TypeId"
        label="账户"
        width="120"
        align="center"
      >
      </el-table-column>
      <el-table-column
        sortable
        prop="InitAmount"
        cell-style="text-align:right;"
        label="初始金额"
        width="120"
        align="right"
      >
      </el-table-column>
      <el-table-column
        sortable
        prop="CurrentAmount"
        cell-style="text-align:right;"
        label="余额"
        width="100"
        align="right"
      >
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
export default {
  data: function () {
    return {
      tableData: [],
      totalAmount: 0,
    };
  },
  methods: {
    async loadBalanceData() {
      const res = await this.$http.get("/api/account/balance");
      console.log("ren", res.data);
      this.tableData = [];
      this.totalAmount = 0;
      if (res.data && res.data.Data && res.data.Data.length > 0) {
        this.tableData = res.data.Data;
        for (let v of this.tableData) {
          this.totalAmount += v.CurrentAmount;
          v.CurrentAmount = (v.CurrentAmount / 100.0).toFixed(2);
        }
        this.totalAmount = (this.totalAmount / 100.0).toFixed(2);
      }
    },
  },
};
</script>