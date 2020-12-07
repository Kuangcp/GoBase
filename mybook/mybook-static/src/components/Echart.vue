<template>
  <div
      ref="echart"
      class="u-echart"
      :style="echartStyle"/>
</template>

<script>
import echarts from 'echarts';
import 'echarts/theme/macarons';

export default {
  name: 'Echart',
  components: {},
  props: {
    echartStyle: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      chart: null,
    };
  },
  computed: {},
  watch: {},
  created() {
  },
  mounted() {
    this.initChart();
  },
  beforeDestroy() {
    if (!this.chart) {
      return;
    }
    this.chart.dispose();
    this.chart = null;
  },
  methods: {
    initChart() {
      // 基于准备好的dom，初始化echarts实例
      this.chart = echarts.init(this.$refs.echart, 'macarons');
    },

    setOption(options) {
      this.chart.setOption(options);
    },
    appendSumLine(lines) {
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
  },
};
</script>

<style lang="scss">
</style>