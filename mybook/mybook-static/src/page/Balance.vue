<template>
  <div>
    <span style="color: green" title="净值">{{ (this.totalAmount / 100.0).toFixed(2) }}</span>
    = <span style="color: blue" title="资产">{{ (this.totalAssets / 100.0).toFixed(2) }}</span>
    - <span style="color: red" title="负债">{{ (this.totalLiabilities / 100.0).toFixed(2) }}</span>

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
        <template slot-scope="scope">
          <span>{{ (scope.row.InitAmount / 100.0).toFixed(2) }}</span>
        </template>
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
          <div v-if="scope.row.CurrentAmount<0">
            <span style="color: red">{{ (scope.row.CurrentAmount / 100.0).toFixed(2) }}</span>
          </div>
          <div v-else>
            <span style="color: blue">{{ (scope.row.CurrentAmount / 100.0).toFixed(2) }}</span>
          </div>
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
      totalLiabilities: 0,
      totalAssets: 0,
    };
  },
  methods: {
    async loadBalanceData() {
      const res = await this.$http.get(window.api.record.balance);
      this.tableData = [];
      this.totalAmount = 0;
      this.totalLiabilities = 0;
      this.totalAssets = 0;

      if (res.data && res.data.data && res.data.data.length > 0) {
        this.tableData = res.data.data;
        for (let v of this.tableData) {
          this.totalAmount += v.CurrentAmount;
          if (v.CurrentAmount < 0) {
            this.totalLiabilities += -1 * v.CurrentAmount
          }
        }
        this.totalAssets = this.totalAmount + this.totalLiabilities
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