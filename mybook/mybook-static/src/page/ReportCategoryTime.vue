<template>
  <div>
    <el-form
        :inline="true"
        ref="ruleForm"
        class="demo-form-inline"
    >
      <el-form-item label="类型">
        <el-select
            v-model="accountType"
            @change="changeAccountType"
            size="mini"
            placeholder="请选择"
            style="width:90px"
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

      <el-form-item label="周期">
        <el-select
            @change="drawLine"
            v-model="timePeriod"
            size="mini"
            placeholder="请选择"
            style="width:60px"
        >
          <el-option
              v-for="item in timePeriods"
              :key="item.ID"
              :label="item.Name"
              :value="item.ID"
          >
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item>
        <el-dropdown size="mini">
          <el-button type="primary" size="mini">
            标记<em class="el-icon-arrow-down el-icon--right"></em>
          </el-button>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item>
              <div class="k-flex">
                <span>明细标记</span>
                <el-switch
                    @change="drawLine"
                    v-model="detailLabel"
                    active-color="#13ce66"
                    inactive-color="#ff4949">
                </el-switch>
              </div>
            </el-dropdown-item>
            <el-dropdown-item>
              <div class="k-flex">
                <span>柱/线</span>
                <el-switch
                    @change="drawLine"
                    v-model="lineChartType"
                    active-color="#13ce66"
                    inactive-color="#ff4949">
                </el-switch>
              </div>
            </el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </el-form-item>

      <el-form-item>
        <el-date-picker
            @change="drawLine"
            v-model="dateArray"
            type="daterange"
            align="right"
            unlink-panels
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            size="mini"
            :picker-options="pickerOptions"
            style="width:240px"
        >
        </el-date-picker>
      </el-form-item>
    </el-form>

    <Echart
        ref="echart"
        v-if="showChart"
        class="categoryMonth"/>
  </div>
</template>

<style scoped>
.categoryMonth {
  width: 1960px;
  height: 760px;
}

.k-flex {
  width: 100px;
  height: 30px;
  align-items: center;
  display: flex;
  justify-content: space-between;
}
</style>

<script>
import {
  dateShortCut,
  dayPeriod,
  formatter,
  getFormatByPeriod,
  monthPeriod,
  weekPeriod,
  yearPeriod,
  calDate,
  passMonth,
} from "@/util/DateUtil";
import Echart from "../components/Echart";

export default {
  components: {
    Echart,
  },
  data: function () {
    return {
      showChart: false,
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
      accountType: 21,
      accountTypes: [
        {ID: 20, Name: "余额"},
        {ID: 21, Name: "收支图"},
        {ID: 1, Name: "支出"},
        {ID: 22, Name: "支出聚合"},
        {ID: 2, Name: "收入"},
        {ID: 23, Name: "收入聚合"},
        {ID: 3, Name: "转帐"},
      ],
      monthChart: "",
      lineChartType: false,
      detailLabel: false,
      showSumLabel: true,
      timePeriod: dayPeriod,
      timePeriods: [
        {ID: yearPeriod, Name: "年"},
        {ID: monthPeriod, Name: "月"},
        {ID: weekPeriod, Name: "周"},
        {ID: dayPeriod, Name: "日"},
      ],
      pickerOptions: {
        shortcuts: dateShortCut
      },
      dateArray: [],
    };
  },
  mounted() {
    this.dateArray = calDate(passMonth)
    this.drawLine()
  },
  methods: {
    async changeAccountType() {
      this.lineChartType = this.accountType === 20 || this.accountType === 21;
      this.drawLine()
    },
    async drawLine() {
      console.log(this.dateArray)
      let startTime = this.dateArray[0];
      let endTime = this.dateArray[1];

      this.showChart = false
      let resp
      if (this.accountType === 20) {
        let start = (startTime && formatter(startTime).format(getFormatByPeriod(this.dayPeriod))) || "";
        let end = (endTime && formatter(endTime).format(getFormatByPeriod(this.dayPeriod))) || "";

        resp = await this.$http.get(window.api.report.balanceReport, {
          params: {
            startDate: start,
            endDate: end,
            chartType: this.lineChartType ? 'line' : 'bar',
          },
        });
      } else {
        let start = (startTime && formatter(startTime).format(getFormatByPeriod(this.timePeriod))) || "";
        let end = (endTime && formatter(endTime).format(getFormatByPeriod(this.timePeriod))) || "";

        resp = await this.$http.get(window.api.report.categoryPeriod, {
          params: {
            startDate: start,
            endDate: end,
            typeId: this.accountType,
            chartType: this.lineChartType ? 'line' : 'bar',
            showLabel: this.detailLabel,
            period: this.timePeriod,
          },
        });
      }
      this.showChart = true

      let respData = resp.data;
      if (respData.code !== 0 || respData.data == null || !respData.data.lines) {
        this.$message({
          message: "分类统计数据为空 " + respData.msg,
          type: "warning",
        });
        return;
      }

      this.$nextTick(() => {
        let finalLines = respData.data.lines
        if (this.accountType !== 21 && this.accountType !== 20 && this.showSumLabel) {
          finalLines = this.$refs.echart.appendSumLine(finalLines)
        }

        this.echartOption.xAxis[0].data = respData.data.xAxis
        this.echartOption.legend.data = respData.data.legends
        this.echartOption.series = finalLines

        this.$refs.echart.setOption(this.echartOption);
      })
    },
  },
};
</script>