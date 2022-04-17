<template>
  <el-form :inline="false" label-width="80px" ref="ruleForm" class="demo-form-inline">
    <el-form-item label="借贷方" required>
      <UserSelect
          ref="userCom"
          :user="userId"
          @hasChange="userChange"
          style="width: 140px"
      />
    </el-form-item>

    <el-form-item label="借贷" required>
      <el-select v-model="loanType" placeholder="请选择" size="mini" style="width: 140px;">
        <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value">
        </el-option>
      </el-select>
    </el-form-item>

    <el-form-item label="账户" required>
      <AccountSelect
          ref="accountCom"
          :account="accountId"
          @hasChange="accountChange"
          style="width: 140px"
      />
    </el-form-item>

    <el-form-item label="金额" required>
      <el-input
          v-model="amount"
          size="mini"
          clearable
          min="0"
          style="width: 140px"
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
    <el-form-item label="偿还">
      <el-date-picker
          v-model="exceptedDate"
          type="date"
          size="mini"
          clearable
          style="width: 140px"
          placeholder="选择日期"
      >
      </el-date-picker>
    </el-form-item>
    <el-form-item label="备注">
      <el-input v-model="comment" size="mini" clearable style="width: 140px"/>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="onSubmit" size="mini">新增</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import {formatter} from "@/util/DateUtil";
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
      loanType: 1,
      amount: "",
      recordDate: "",
      exceptedDate: "",
      comment: "",
      options: [{
        value: 1,
        label: '借入'
      }, {
        value: 2,
        label: '贷出'
      }, {
        value: 3,
        label: '贷出-负债'
      }, {
        value: 4,
        label: '借入-负债'
      }],
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
      // TODO 提交， 借贷 ，操作账户： 应付款，应收款

      let recordDateFmt = formatter(this.recordDate).formatDate()
      let exceptedDateFmt = formatter(this.exceptedDate).formatDate()

      let param = {
        userId: this.userId,
        accountId: this.$refs.accountCom.account || 0,
        amount: this.amount,
        loanType: this.loanType,
        date: recordDateFmt,
        exceptedDate: exceptedDateFmt,
        comment: this.comment,
      };

      console.log(param);
      let resp = await this.$http.post(window.api.loan.create, param);
      // console.log(resp);

      if (resp.data.code !== 0) {
        this.$message({
          message: resp.data.msg,
          type: "warning",
        });
      } else {
        this.$message({
          message: "新增成功",
          type: "success",
        });
      }
    }
  }
}
</script>

<style scoped>

</style>