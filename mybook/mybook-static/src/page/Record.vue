<template>
  <div>
    <el-form :inline="true" ref="ruleForm" class="demo-form-inline">
      <el-form-item>
        <CategorySelect ref="categoryCom" />
      </el-form-item>
      <!-- <br> -->
      <el-form-item label="操作账户">
        <AccountSelect ref="accountCom" style="width: 120px" />
      </el-form-item>
      <el-form-item label="=> 转账目标账户">
        <AccountSelect ref="targetAccountCom" style="width: 120px" />
      </el-form-item>
      <el-form-item label="金额">
        <el-input
          v-model.number="amount"
          size="mini"
          clearable
          min="0"
          style="width: 100px"
        />
      </el-form-item>
      <el-form-item label="时间">
        <el-date-picker
          v-model="recordDate"
          type="date"
          size="mini"
          style="width: 132px"
          placeholder="选择日期"
        >
        </el-date-picker>
      </el-form-item>
      <el-form-item label="备注">
        <el-input v-model="comment" size="mini" style="width: 100px" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit" size="mini">新增</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import DateUtil from "../util/DateUtil.js";
import AccountSelect from "../components/AccountSelect";
import CategorySelect from "../components/CategorySelect";

export default {
  components: {
    AccountSelect,
    CategorySelect,
  },
  data: function () {
    return {
      amount: 0,
      recordDate: "",
      comment: "",
    };
  },
  methods: {
    async onSubmit() {
      // console.log(param);
      let ids = this.$refs.categoryCom.categoryId;
      console.log(ids);

      let param = {
        typeId: ids[0] + "",
        accountId: this.$refs.accountCom.account + "",
        targetAccountId: this.$refs.targetAccountCom.account + "",
        categoryId: ids[ids.length - 1] + "",
        amount: this.amount + "",
        date: DateUtil(this.recordDate).formatDate(),
        comment: this.comment,
      };
      this.$http.post("/api/record/createRecord", param);
    },
  },
};
</script>