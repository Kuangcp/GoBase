<template>
  <div>
    余额：{{ (this.totalAmount / 100.0).toFixed(2) }}
    <el-table :data="tableData" stripe class="main-box" height="880">
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
        <template slot-scope="scope">
          <span>{{ (scope.row.CurrentAmount / 100.0).toFixed(2) }}</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<style scoped>
.main-box {
  width: 550px;
  margin-left: 2vw;
  border-radius: 4px;
}
</style>

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
      const res = await this.$http.get(window.api.record.balance);
      this.tableData = [];
      this.totalAmount = 0;
      if (res.data && res.data.Data && res.data.Data.length > 0) {
        this.tableData = res.data.Data;
        for (let v of this.tableData) {
          this.totalAmount += v.CurrentAmount;
        }
      } else {
        this.$message({
          message: "数据为空",
          type: "warning",
        });
      }
    },
  },
};
</script>