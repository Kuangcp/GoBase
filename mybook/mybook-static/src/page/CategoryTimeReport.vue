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
            标记<i class="el-icon-arrow-down el-icon--right"></i>
          </el-button>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item>
              <div class="k-flex">
                <span>明细标记</span>
                <el-switch
                    v-model="detailLabel"
                    active-color="#13ce66"
                    inactive-color="#ff4949">
                </el-switch>
              </div>
            </el-dropdown-item>
            <el-dropdown-item>
              <div class="k-flex">
                <span>求和</span>
                <el-switch
                    v-model="showSumLabel"
                    active-color="#13ce66"
                    inactive-color="#ff4949">
                </el-switch>
              </div>
            </el-dropdown-item>
            <el-dropdown-item>
              <div class="k-flex">
                <span>柱/线</span>
                <el-switch
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

      <el-form-item>
        <el-button type="primary" @click="drawLine" size="mini">查询</el-button>
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
import DateUtil from "../util/DateUtil.js";
import Echart from "../components/Echart";

function fillDate(picker, offset) {
  const end = new Date();
  const start = new Date();
  start.setTime(start.getTime() - offset);
  picker.$emit("pick", [start, end]);
}

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
      accountType: 1,
      accountTypes: [
        {ID: 1, Name: "支出"},
        {ID: 2, Name: "收入"},
        {ID: 0, Name: "收支图"},
        {ID: 3, Name: "转帐"},
      ],
      monthChart: "",
      lineChartType: false,
      detailLabel: false,
      showSumLabel: false,
      timePeriod: "month",
      timePeriods: [
        {ID: "year", Name: "年"},
        {ID: "month", Name: "月"},
        {ID: "week", Name: "周"},
        {ID: "day", Name: "日"},
      ],
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
          {
            text: "最近一年",
            onClick(picker) {
              fillDate(picker, 3600 * 1000 * 24 * 360);
            },
          },
        ],
      },
      dateArray: [],
    };
  },
  mounted() {
  },
  methods: {
    getFormat() {
      switch (this.timePeriod) {
        case "year":
          return "YYYY"
        case "month":
          return "YYYY-MM"
        case "day":
          return "YYYY-MM-dd"
      }
    },

    async drawLine() {
      let startTime = this.dateArray[0];
      let endTime = this.dateArray[1];
      let start = (startTime && DateUtil(startTime).format(this.getFormat())) || "";
      let end = (endTime && DateUtil(endTime).format(this.getFormat())) || "";

      this.showChart = false
      let resp = await this.$http.get(window.api.report.categoryPeriod, {
        params: {
          startDate: start,
          endDate: end,
          typeId: this.accountType,
          chartType: this.lineChartType ? 'line' : 'bar',
          showLabel: this.detailLabel,
          period: this.timePeriod,
        },
      });
      this.showChart = true

      let respData = resp.data;
      if (respData.Code !== 0) {
        this.$message({
          message: respData.Msg,
          type: "warning",
        });
        return;
      }
      if (respData.Data == null || !respData.Data.lines) {
        this.$message({
          message: "分类统计数据为空",
          type: "warning",
        });
        return;
      }

      // console.log(this.echartOption)
      this.$nextTick(() => {
        let finalLines = respData.Data.lines
        if (this.showSumLabel) {
          finalLines = this.$refs.echart.appendSumLine(finalLines)
        }

        this.echartOption.xAxis[0].data = respData.Data.xAxis
        this.echartOption.legend.data = respData.Data.legends
        this.echartOption.series = finalLines

        this.$refs.echart.setOption(this.echartOption);
      })
    },
  },
};
</script>