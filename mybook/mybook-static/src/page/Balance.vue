<template>
  <div>
    <span style="color: green" title="净值">{{ (this.totalAmount / 100.0).toFixed(2) }}</span>
    &nbsp;<span style="color: blue" title="资产">{{ (this.totalAssets / 100.0).toFixed(2) }}</span>
    &nbsp;️<span style="color: red" title="负债">{{ (this.totalLiabilities / 100.0).toFixed(2) }}</span>

    <br>
    <el-table :data="tableData" stripe class="balance-box">
      <el-table-column sortable prop="ID" label="ID" width="70" align="right">
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
          prop="TypeName"
          label="类型"
          width="80"
          align="center"
      >
      </el-table-column>
      <el-table-column
          sortable
          prop="InitAmount"
          cell-style="text-align:right;"
          label="初始"
          width="100"
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
          width="140"
          align="right"
      >
        <template slot-scope="scope">
          <div v-if="scope.row.CurrentAmount<0">
            <span style="color: red">{{ (scope.row.CurrentAmount / 100.0).toFixed(2) }}</span>
          </div>
          <div v-else-if="scope.row.CurrentAmount===0">
            <span style="color: grey">0</span>
          </div>
          <div v-else>
            <span style="color: green">{{ (scope.row.CurrentAmount / 100.0).toFixed(2) }}</span>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-table :data="loanUser" stripe class="loan-box">
      <el-table-column sortable prop="UserId" label="ID" width="70" align="right"/>
      <el-table-column sortable prop="Name" label="用户" width="120" align="right"/>
      <el-table-column sortable prop="Amount" label="金额" width="140" align="right">
        <template slot-scope="scope">
          <div v-if="scope.row.Amount<0">
            <span style="color: red">借入 {{ (scope.row.Amount / 100.0 * -1).toFixed(2) }}</span>
          </div>
          <div v-else-if="scope.row.Amount===0">
            <span style="color: grey">{{ (scope.row.Amount / 100.0 * -1).toFixed(2) }}</span>
          </div>
          <div v-else>
            <span style="color: green">贷出 {{ (scope.row.Amount / 100.0).toFixed(2) }}</span>
          </div>
        </template>
      </el-table-column>
    </el-table>

  </div>
</template>

<style scoped>
.balance-box {
  width: 520px;
  height: 800px;
  border-radius: 4px;
  float: left;
}

.loan-box {
  margin-left: 1vw;
  width: 340px;
  height: 800px;
  border-radius: 4px;
  float: left;
}
</style>

<script>
export default {
  data: function () {
    return {
      tableData: [],
      loanUser: [],
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

      const loans = await this.$http.get(window.api.loan.query);
      console.log(loans)
      if (loans.data && loans.data.data && loans.data.data.length > 0) {
        this.loanUser = loans.data.data;
      }
    },
  },
};
</script>