<template>
  <div>
    <el-form :inline="true" class="demo-form-inline">
      <el-form-item label="账户">
        <AccountSelect ref="accountCom" />
      </el-form-item>
      <el-form-item label="转账目标账户">
        <AccountSelect ref="targetAccountCom" />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="onSubmit" size="mini">查询</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import AccountSelect from "../components/AccountSelect";

export default {
  components: {
    AccountSelect,
  },
  data: function () {
    return {};
  },
  methods: {
    async onSubmit() {
      console.log(this.$refs.accountCom.account);
      console.log(this.$refs.targetAccountCom.account);
    },
  },
};
</script>