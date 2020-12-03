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

      <el-form-item label="类型">
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

      <el-form-item label="明细标记">
        <el-switch
            v-model="detailLabel"
            active-color="#13ce66"
            inactive-color="#ff4949">
        </el-switch>
      </el-form-item>

      <el-form-item label="求和">
        <el-switch
            v-model="showSumLabel"
            active-color="#13ce66"
            inactive-color="#ff4949">
        </el-switch>
      </el-form-item>

      <el-form-item label="柱/线">
        <el-switch
            v-model="lineChartType"
            active-color="#13ce66"
            inactive-color="#ff4949">
        </el-switch>
      </el-form-item>

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
            style="width:240px"
        >
        </el-date-picker>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="drawLine" size="mini">查询</el-button>
      </el-form-item>
    </el-form>

    <div
        id="categoryMonthDiv"
        ref="chart"
        class="categoryMonth"
    ></div>
  </div>
</template>

<style scoped>
.categoryMonth {
  width: 1860px;
  height: 760px;
}
</style>

<script>
import DateUtil from "../util/DateUtil.js";

var echarts = require("echarts");

function appendSumLine(lines) {
  let sumData = [];
  let first = lines[0];
  for (let i = 0; i < first.data.length; i++) {
    let temp = 0;
    for (let j = 0; j < lines.length; j++) {
      temp += lines[j].data[i];
    }
    if (temp === 0) {
      sumData.push(0);
    } else {
      sumData.push(temp.toFixed(2));
    }
  }

  lines.push({
    //新的一个柱子 注意不设stack
    name: "累计",
    type: "bar",
    barGap: "-100%", // 左移100%，stack不再与上面两个同列
    label: {
      normal: {
        show: true, //显示数值
        position: "top", //  位置设为top
        formatter: "{c}",
        textStyle: {color: "#213e53"}, //设置数值颜色
      },
    },
    itemStyle: {
      normal: {
        color: "rgba(128, 128, 128, 0)", // 设置背景颜色为透明
      },
    },
    data: sumData,
  });
  return lines;
}

function fillDate(picker, offset) {
  const end = new Date();
  const start = new Date();
  start.setTime(start.getTime() - offset);
  picker.$emit("pick", [start, end]);
}

export default {
  components: {},
  data: function () {
    return {
      accountType: 1,
      accountTypes: [
        {ID: 1, Name: "支出"},
        {ID: 2, Name: "收入"},
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
      let categoryMonthDiv = this.$refs.chart;
      if (categoryMonthDiv) {
        let startTime = this.dateArray[0];
        let endTime = this.dateArray[1];
        let start = (startTime && DateUtil(startTime).format(this.getFormat())) || "";
        let end = (endTime && DateUtil(endTime).format(this.getFormat())) || "";

        let resp = await this.$http.get(window.api.report.categoryMonth, {
          params: {
            startDate: start,
            endDate: end,
            typeId: this.accountType,
            chartType: this.lineChartType ? 'line' : 'bar',
            showLabel: this.detailLabel,
            period: this.timePeriod,
          },
        });

        if (resp.data.Data == null || resp.data.Data.length === 0) {
          this.$message({
            message: "分类统计数据为空",
            type: "warning",
          });
          return;
        }

        if (this.monthChart !== "") {
          this.monthChart.dispose()
        }
        let myChart = echarts.init(categoryMonthDiv);
        this.monthChart = myChart
        let option = {
          title: {
            text: "",
          },
          tooltip: {
            trigger: "axis",
            axisPointer: {
              type: "cross",
              label: {
                backgroundColor: "#6a7985",
              },
            },
          },
          legend: {
            data: [],
          },
          toolbox: {
            feature: {
              saveAsImage: {},
            },
          },
          grid: {
            left: "3%",
            right: "4%",
            bottom: "3%",
            containLabel: true,
          },
          xAxis: [
            {
              type: "category",
              boundaryGap: false,
              data: [],
            },
          ],
          yAxis: [
            {
              type: "value",
            },
          ],
          series: [],
        };
        myChart.setOption(option);

        let finalLines = resp.data.Data.lines
        if (this.showSumLabel) {
          finalLines = appendSumLine(finalLines)
        }
        myChart.setOption({
          xAxis: {
            data: resp.data.Data.xAxis,
          },
          legend: {
            data: resp.data.Data.legends,
          },
          series: finalLines,
        });
      } else {
        console.error("div 不存在");
      }
    },
  },
};
</script>