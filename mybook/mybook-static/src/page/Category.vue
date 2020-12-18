<template>
  <div>
    <el-form :inline="true" class="demo-form-inline form-nav">
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

      <el-form-item>
        <el-button type="primary" @click="onSubmit" size="mini">查询</el-button>
      </el-form-item>
      <el-form-item>
        <div>总金额:{{ totalAmount }}</div>
      </el-form-item>
    </el-form>

    <el-dialog title="明细" :visible.sync="detailTableDialogVisible" width="760px">
      <el-table :data="detailData" stripe style="width: 100%" height="660">
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
            label="金额"
            align="right"
            width="120"
        >
          <template slot-scope="scope">
            <span>{{ (scope.row.Amount / 100.0).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column sortable prop="RecordTime" label="时间" width="100">
        </el-table-column>
        <el-table-column prop="Comment" label="备注" width="180">
        </el-table-column>
      </el-table>
    </el-dialog>
    <el-dialog title="报表" :visible.sync="reportDialogVisible" width="1820px">
      <Echart ref="echart" v-if="reportDialogVisible" class="categoryMonth"/>
    </el-dialog>

    <el-table :data="tableData" stripe height="860" class="main-box">
      <el-table-column
          sortable
          prop="CategoryId"
          label="ID"
          width="80"
          align="right"
      >
      </el-table-column>
      <el-table-column
          prop="RecordTypeName"
          label="类型"
          width="80"
      ></el-table-column>
      <el-table-column
          prop="Name"
          label="明细类型"
          width="100"
          align="right"
      ></el-table-column>

      <!-- <el-table-column prop="RecordType" label="类型" width="180"> -->

      <el-table-column
          sortable
          prop="Amount"
          label="金额"
          align="right"
          width="140"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.Amount.toFixed(2) }}</span>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="180">
        <template slot-scope="scope">
          <el-button
              @click.native.prevent="detail(scope.row.CategoryId)"
              type="text"
              size="mini"
          >
            详情
          </el-button>
          <el-button
              @click.native.prevent="weekDetail(scope.row.CategoryId)"
              type="text"
              size="mini"
          >
            周统计
          </el-button>
          <el-button
              @click.native.prevent="monthDetail(scope.row.CategoryId)"
              type="text"
              size="mini"
          >
            月统计
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<style scoped>
/deep/ .el-dialog__body {
  padding: 0 0 0 10px;
}

/deep/ .el-dialog__header {
  padding-bottom: 0;
}

.form-nav {
  margin-bottom: -16px;
}

.main-box {
  width: 600px;
  margin-left: 2vw;
  border-radius: 4px;
}

.categoryMonth {
  width: 1800px;
  height: 650px;
}
</style>
<script>
import {dateShortCut, formatter, getFormatByPeriod, monthPeriod} from "@/util/DateUtil";
import Echart from "../components/Echart";

export default {
  components: {
    Echart,
  },
  data: function () {
    return {
      detailTableDialogVisible: false,
      reportDialogVisible: false,
      accountId: null,
      pickerOptions: {
        shortcuts: dateShortCut
      },
      dateArray: [],
      visible: false,
      tableData: [],
      detailData: [],
      accountType: "",
      accountTypes: [
        {ID: 1, Name: "支出"},
        {ID: 2, Name: "收入"},
        {ID: 3, Name: "转出"},
        {ID: 4, Name: "转入"},
      ],
      totalAmount: 0,
      echartOption: {
        title: {
          text: '',
        },
        tooltip: {
          trigger: 'axis',
        },
        color: ['#409EFF'],
        legend: {
          data: [],
        },
        toolbox: {
          show: true,
          feature: {
            // mark: { show: true },
            // dataView: { show: true, readOnly: false },
            // magicType: { show: true, type: ['line', 'bar'] },
            // restore: { show: true },
            saveAsImage: {show: true},
          },
        },
        calculable: true,
        xAxis: [
          {
            type: 'category',
            boundaryGap: true,
            data: [],
          },
        ],
        yAxis: [
          {
            type: 'value',
            axisLabel: {
              formatter: '{value}',
            },
          },
        ],
        series: [],
      },
    };
  },
  mounted() {
  },
  methods: {
    getFormatDate() {
      let startTime = this.dateArray[0];
      let endTime = this.dateArray[1];
      let start = (startTime && formatter(startTime).formatDate()) || "";
      let end = (endTime && formatter(endTime).formatDate()) || "";
      return {
        start,
        end,
      };
    },
    async onSubmit() {
      const {start, end} = this.getFormatDate();

      const res = await this.$http.get(window.api.record.byCategory, {
        params: {
          startDate: start,
          endDate: end,
          typeId: this.accountType,
        },
      });

      this.tableData = [];
      this.totalAmount = 0;
      if (res.data && res.data.Data && res.data.Data.length > 0) {
        this.tableData = res.data.Data;
        this.totalAmount = 0;
        for (let v of this.tableData) {
          this.totalAmount += v.Amount;
          v.Amount = v.Amount / 100.0;
        }
        this.totalAmount = this.totalAmount / 100.0;
      }
    },
    listenAccount(val) {
      this.accountId = val;
    },
    async detail(categoryId) {
      this.detailTableDialogVisible = true;
      const {start, end} = this.getFormatDate();
      const res = await this.$http.get(window.api.record.categoryDetail, {
        params: {
          startDate: start,
          endDate: end,
          categoryId: categoryId,
          typeId: this.accountType,
        },
      });
      this.detailData = [];
      this.detailData = res.data.Data;
    },
    async report(categoryId, period) {
      this.reportDialogVisible = true;
      let startTime = this.dateArray[0];
      let endTime = this.dateArray[1];
      let start = (startTime && formatter(startTime).format(getFormatByPeriod(period))) || "";
      let end = (endTime && formatter(endTime).format(getFormatByPeriod(period))) || "";

      let resp = await this.$http.get(window.api.report.categoryPeriod, {
        params: {
          startDate: start,
          endDate: end,
          chartType: 'line',
          typeId: this.accountType,
          categoryId: categoryId,
          showLabel: false,
          period: period,
        },
      });

      let respData = resp.data;
      if (respData.Code !== 0 || respData.Data == null || !respData.Data.lines) {
        this.$message({
          message: "分类统计数据为空 " + respData.Msg,
          type: "warning",
        });
        return;
      }

      this.$nextTick(() => {
        let finalLines = respData.Data.lines
        finalLines = this.$refs.echart.appendSumLine(finalLines)

        this.echartOption.xAxis[0].data = respData.Data.xAxis
        this.echartOption.legend.data = respData.Data.legends
        this.echartOption.series = finalLines

        this.$refs.echart.setOption(this.echartOption);
      })
    },
    weekDetail(categoryId) {
      console.log(categoryId);
    },
    async monthDetail(categoryId) {
      this.report(categoryId, monthPeriod)
    },
  },
};
</script>