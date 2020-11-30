<template>
  <div id="bar_dv" ref="chart" style="width: 600px; height: 400px"></div>
</template>
 
<script>
var echarts = require("echarts");

export default {
  name: "EHistogram",
  data() {
    return {
      msg: "Welcome to Your Vue.js App",
      label: "",
      itemColor: "red",
      backgroundColor: "white",
      itemDataType: "",
      xAxisName: "x",
      yAxisName: "y",
      eventType: "line",
    };
  },
  mounted() {
    this.drawLine();
  },
  methods: {
    async drawLine() {
      var bar_dv = this.$refs.chart;
      if (bar_dv) {
        let myChart = echarts.init(bar_dv);

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
        let resp = await this.$http.get("/api/report/categoryMonth", {
          params: {
            startDate: "2018-11",
            endDate: "2020-11",
            typeId: 1,
            chartType: "line",
            showLabel: false,
          },
        });
        myChart.setOption({
          xAxis: {
            data: resp.data.Data.xAxis,
          },
          legend: {
            data: resp.data.Data.legends,
          },
          series: resp.data.Data.lines,
        });
      } else {
        console.error("div 不存在");
      }
    },
  },
};
</script>
 
<style scoped>
</style> 