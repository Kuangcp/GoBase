<template>
  <div>
    <el-form :inline="true" class="demo-form-inline">
      <el-form-item>
        <el-radio v-model="recordType" label="1">支出</el-radio>
        <el-radio v-model="recordType" label="2">收入</el-radio>
        <el-radio v-model="recordType" label="3">转账</el-radio>
        <!-- 由转出(3)触发生成转入(4) -->
      </el-form-item>
      <br />
      <el-form-item label="操作账户">
        <AccountSelect ref="accountCom" style="width:120px;" />
      </el-form-item>
      <el-form-item label="=> 转账目标账户">
        <AccountSelect ref="targetAccountCom" style="width:120px;"/>
      </el-form-item>
      <el-form-item label="金额">
        <el-input v-model="amount" size="mini" clearable min="0" style="width:100px;" />
      </el-form-item>
      <el-form-item label="时间">
        <el-date-picker v-model="recordDate" type="date" size="mini" style="width:132px;" placeholder="选择日期">
        </el-date-picker>
      </el-form-item>
      <el-form-item label="备注">
        <el-input v-model="comment" size="mini" style="width:100px;"/>
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

export default {
  components: {
    AccountSelect,
  },
  data: function () {
    return {
      recordType: "1",
      amount: 0,
      recordDate: "",
      comment: "",
    };
  },
  methods: {
    async onSubmit() {
      let param = {
        typeId: this.recordType,
        accountId: this.$refs.accountCom.account,
        targetAccountId: this.$refs.targetAccountCom.account,
        categoryId: 0,
        amount: this.amount,
        date: DateUtil(this.recordDate).formatDate(),
        comment: this.comment,
      };
      console.log(param);
    },
  },
};
</script>