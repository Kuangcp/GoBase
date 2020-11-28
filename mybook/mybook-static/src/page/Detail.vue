<template>
  <div>
    <el-button @click="showDiv()" ref="viewBtn">Button</el-button>

    <el-dialog :visible.sync="visible">
      <div
        style="display: none"
        :style="{ display: visible ? 'block' : 'none' }"
      >
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
          <el-form-item label="分类">
            <el-select v-model="categoryType" size="mini" placeholder="请选择">
              <el-option
                v-for="item in categoryTypes"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="onSubmit" size="mini"
              >查询</el-button
            >
          </el-form-item>
        </el-form>

        <el-table :data="tableData" stripe style="width: 100%">
          <el-table-column prop="date" label="日期" width="180">
          </el-table-column>
          <el-table-column prop="name" label="姓名" width="180">
          </el-table-column>
          <el-table-column prop="address" label="地址"> </el-table-column>
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import DateUtil from '../util/DateUtil.js'

export default {
  data: function() {
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
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近一个月",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近三个月",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
              picker.$emit("pick", [start, end]);
            },
          },
        ],
      },
      dateArray: [],
      visible: false,
      tableData: [],
      obj: {
        // inputVal: 1
      },
      categoryTypes: [
        {
          value: "选项1",
          label: "黄金糕",
        },
        {
          value: "选项2",
          label: "双皮奶",
        },
        {
          value: "选项3",
          label: "蚵仔煎",
        },
        {
          value: "选项4",
          label: "龙须面",
        },
        {
          value: "选项5",
          label: "北京烤鸭",
        },
      ],
      categoryType: "",
    };
  },
  methods: {
    watchInput(val) {
      console.log(val);
    },
    showDiv() {
      this.visible = true;
      console.log(this.$refs.viewBtn);
    },
    onSubmit() {
      let startTime = this.dateArray[0];
      let endTime = this.dateArray[1];

      let startStr =
        (startTime && DateUtil(startTime).format("YYYY-MM-dd")) || "";
      console.log(startTime, startStr, endTime);
      this.tableData = [
        {
          date: "2016-05-02",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1518 弄",
        },
        {
          date: "2016-05-02",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1518 弄",
        },
      ];
      this.$set(this.obj, "inputVal", 0);
    },
  },
};
</script>
