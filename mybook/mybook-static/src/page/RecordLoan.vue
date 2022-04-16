<template>
  <el-form :inline="false" label-width="80px" ref="ruleForm" class="demo-form-inline">
    <el-form-item label="借贷方" required>
      <UserSelect
          ref="userCom"
          :user="userId"
          @hasChange="userChange"
          style="width: 120px"
      />
    </el-form-item>

    <el-form-item label="交易账户" required>
      <AccountSelect
          ref="accountCom"
          :account="accountId"
          @hasChange="accountChange"
          style="width: 120px"
      />
    </el-form-item>

    <el-form-item label="金额" required>
      <el-input
          v-model="amount"
          size="mini"
          clearable
          min="0"
          style="width: 120px"
      />
    </el-form-item>

    <el-form-item label="时间" required>
      <el-date-picker
          v-model="recordDate"
          type="date"
          size="mini"
          clearable
          style="width: 140px"
          placeholder="选择日期"
      >
      </el-date-picker>
    </el-form-item>
    <el-form-item label="备注">
      <el-input v-model="comment" size="mini" clearable style="width: 200px"/>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="onSubmit" size="mini">新增</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import AccountSelect from "../components/AccountSelect";
import UserSelect from "../components/UserSelect";

export default {
  name: "LoanRecord",
  components: {
    UserSelect,
    AccountSelect,
  },
  data: function () {
    return {
      userId: 1,
      accountId: 7,
      amount: "",
      recordDate: "",
      comment: "",
    };
  },
  methods: {
    accountChange(val) {
      this.accountId = val;
    },
    userChange(val) {
      this.userId = val;
    },
    async onSubmit() {
      console.log('submit')
    }
  }
}
</script>

<style scoped>

</style>