<template>
  <div>
    <el-form :inline="true" class="demo-form-inline">
      <el-form-item label="时间">
        <el-date-picker
          v-model="dateArray"
          type="daterange"
          align="right"
          unlink-panels
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          size="mini"
          :picker-options="pickerOptions"
        >
        </el-date-picker>
      </el-form-item>
      <el-form-item label="类型">
        <el-select
          v-model="accountType"
          size="mini"
          clearable
          placeholder="请选择"
        >
          <el-option
            v-for="item in accountTypes"
            :key="item.ID"
            :label="item.Name"
            :value="item.ID"
          >
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="账户">
        <AccountSelect ref="accountCom"/>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="onSubmit" size="mini">查询</el-button>
      </el-form-item>
      <el-form-item>
        <div>总金额:{{ totalAmount }}</div>
      </el-form-item>
    </el-form>

    <el-table :data="tableData" stripe style="width: 100%" height="800">
      <el-table-column sortable prop="ID" label="ID" width="60" align="right">
      </el-table-column>
      <el-table-column
        sortable
        prop="AccountName"
        label="账户"
        width="120"
        align="center"
      >
      </el-table-column>
      <el-table-column
        prop="RecordTypeName"
        label="类型"
        width="60"
      ></el-table-column>
      <el-table-column
        prop="CategoryName"
        label="明细类型"
        width="100"
        align="right"
      ></el-table-column>

      <!-- <el-table-column prop="RecordType" label="类型" width="180"> -->

      <el-table-column
        sortable
        prop="Amount"
        cell-style="text-align:right;"
        label="金额"
        width="100"
        align="right"
      >
      </el-table-column>
      <el-table-column sortable prop="RecordTime" label="时间" width="190">
      </el-table-column>
      <el-table-column prop="Comment" label="备注" width="200">
      </el-table-column>
    </el-table>
  </div>
</template>

<style>
</style>
<script>
import DateUtil from "../util/DateUtil.js";
import AccountSelect from "../components/AccountSelect";

function fillDate(picker, offset) {
  const end = new Date();
  const start = new Date();
  start.setTime(start.getTime() - offset);
  picker.$emit("pick", [start, end]);
}

export default {
  components: {
    AccountSelect,
  },
  data: function () {
    return {
      pickerOptions: {
        shortcuts: [
          {
            text: "最近一周",
            onClick(picker) {
              fillDate(picker, 3600 * 1000 * 24 * 7);
            },
          },
          {
            text: "最近一个月",
            onClick(picker) {
              fillDate(picker, 3600 * 1000 * 24 * 30);
            },
          },
          {
            text: "最近三个月",
            onClick(picker) {
              fillDate(picker, 3600 * 1000 * 24 * 90);
            },
          },
          {
            text: "最近半年",
            onClick(picker) {
              fillDate(picker, 3600 * 1000 * 24 * 180);
            },
          },
        ],
      },
      dateArray: [],
      visible: false,
      tableData: [],
      accountType: "",
      accountTypes: [
        { ID: 1, Name: "支出" },
        { ID: 2, Name: "收入" },
        { ID: 3, Name: "转出" },
        { ID: 4, Name: "转入" },
      ],
      totalAmount: 0,
    };
  },
  mounted() {},
  methods: {
    async onSubmit() {
      let startTime = this.dateArray[0];
      let endTime = this.dateArray[1];
      let start = (startTime && DateUtil(startTime).format("YYYY-MM-dd")) || "";
      let end = (endTime && DateUtil(endTime).format("YYYY-MM-dd")) || "";

      const res = await this.$http.get("/api/record/list", {
        params: {
          startDate: start,
          endDate: end,
          typeId: this.accountType,
          accountId: this.$refs.accountCom.account,
        },
      });

      this.tableData = [];
      this.totalAmount = 0;
      if (res.data && res.data.Data && res.data.Data.length > 0) {
        this.tableData = res.data.Data;
        this.totalAmount = 0;
        for (let v of this.tableData) {
          this.totalAmount += v.Amount;
          v.Amount = (v.Amount / 100.0).toFixed(2);
        }
        this.totalAmount = (this.totalAmount / 100.0).toFixed(2);
      }
    },
  },
};
</script>
