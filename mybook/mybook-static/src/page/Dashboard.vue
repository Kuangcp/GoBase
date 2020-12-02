<template>
  <div id="categoryMonthDiv" ref="chart" style="width: 1200px; height: 400px"></div>
</template>
<script>
var echarts = require("echarts");

function appendSumLine(lines) {
  let sumData = []
  let first = lines[0];
  for (let i = 0; i < first.data.length; i++) {
    let temp = 0
    for (let j = 0; j < lines.length; j++) {
      temp += lines[j].data[i]
    }
    sumData.push(temp.toFixed(2))
  }

  lines.push({ //新的一个柱子 注意不设stack
    name: '累计',
    type: 'bar',
    barGap: '-100%', // 左移100%，stack不再与上面两个同列
    label: {
      normal: {
        show: true, //显示数值
        position: 'top', //  位置设为top
        formatter: '{c}',
        textStyle: {color: '#213e53'} //设置数值颜色
      }
    },
    itemStyle: {
      normal: {
        color: 'rgba(128, 128, 128, 0)' // 设置背景颜色为透明
      }
    },
    data: sumData,
  })
  return lines
}

export default {
  data() {
    return {};
  },
  mounted() {
    this.drawLine();
  },
  methods: {
    async drawLine() {
      let categoryMonthDiv = this.$refs.chart;
      if (categoryMonthDiv) {
        let myChart = echarts.init(categoryMonthDiv);

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
        let resp = await this.$http.get(window.api.report.categoryMonth, {
          params: {
            startDate: "2019-01",
            endDate: "2020-12",
            typeId: 1,
            chartType: "bar",
            showLabel: false,
            period: 'month',
          },
        });
        myChart.setOption({
          xAxis: {
            data: resp.data.Data.xAxis,
          },
          legend: {
            data: resp.data.Data.legends,
          },
          series: appendSumLine(resp.data.Data.lines),
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