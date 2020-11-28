<template>
  <div>
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
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
      <el-form-item label="账户">
        <el-select v-model="account" size="mini" placeholder="请选择">
          <el-option
            v-for="item in accounts"
            :key="item.ID"
            :label="item.Name"
            :value="item.ID"
          >
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="onSubmit" size="mini">查询</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="tableData" stripe style="width: 100%">
      <el-table-column prop="BillDay" label="日期" width="180">
      </el-table-column>
      <el-table-column prop="CurrentAmount" label="姓名" width="180">
      </el-table-column>
      <el-table-column prop="ID" label="地址"> </el-table-column>
    </el-table>
  </div>
</template>

<script>
import DateUtil from "../util/DateUtil.js";

function fillDate(picker, offset) {
  const end = new Date();
  const start = new Date();
  start.setTime(start.getTime() - offset);
  picker.$emit("pick", [start, end]);
}

export default {
  data: function () {
    return {
      formInline: {
        user: "",
        region: "",
      },
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
      account: "",
      accounts: [],
    };
  },
  mounted() {
    this.fillAccount();
  },
  methods: {
    async onSubmit() {
      console.log(this.account);
      let startTime = this.dateArray[0];
      let endTime = this.dateArray[1];

      const res = await this.$http.get("/api/account/list");
      console.log("ren", res.data);

      let startStr =
        (startTime && DateUtil(startTime).format("YYYY-MM-dd")) || "";
      console.log(startTime, startStr, endTime);
      this.tableData = res.data.Data;
      // 添加属性
      //   this.$set(this.obj, "inputVal", 0);
    },
    async queryAllAccount() {
      const res = await this.$http.get("/api/account/list");
      console.log("ren", res.data);
      return res.data.Data;
    },
    async fillAccount() {
      this.accounts = [];
      let result = await this.queryAllAccount();
      console.log("result:::::::::::", result);
      this.accounts = result;
    },
  },
};
</script>
