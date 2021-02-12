<template>
  <div>
    <el-form :inline="false" label-width="80px" ref="ruleForm" class="demo-form-inline">
      <el-form-item label="分类" required>
        <CategorySelect ref="categoryCom"/>
      </el-form-item>

      <el-form-item label="操作账户" required>
        <AccountSelect
            ref="accountCom"
            :account="accountId"
            @hasChange="hasChange"
            style="width: 120px"
        />
      </el-form-item>

      <el-form-item label="目标账户">
        <el-tooltip class="item" effect="dark" content="仅在转账时有效" placement="top-start">
          <AccountSelect
              ref="targetAccountCom"
              :account="targetAccountId"
              @hasChange="listenTargetAccount"
              style="width: 120px"
          />
        </el-tooltip>
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
        OR
        <el-date-picker
            v-model="recordDates"
            type="dates"
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
  </div>
</template>

<script>
import {formatter} from "@/util/DateUtil";
import AccountSelect from "../components/AccountSelect";
import CategorySelect from "../components/CategorySelect";

export default {
  components: {
    AccountSelect,
    CategorySelect,
  },
  data: function () {
    return {
      accountId: 2,
      targetAccountId: 1,
      amount: "",
      recordDate: "",
      recordDates: [],
      comment: "",
    };
  },
  methods: {
    async onSubmit() {
      let ids = this.$refs.categoryCom.categoryId;
      if ((!this.recordDates || this.recordDates.length === 0)
          && (!this.recordDate || this.recordDate.length === 0)) {
        this.$message({
          message: "时间为空",
          type: "warning",
        });
        return;
      }

      let resultDate = [formatter(this.recordDate).formatDate()]
      if (!this.recordDate || this.recordDate.length === 0) {
        resultDate = this.recordDates.map((v) => formatter(v).formatDate());
      }

      let param = {
        typeId: ids[0],
        accountId: this.$refs.accountCom.account || 0,
        targetAccountId: this.$refs.targetAccountCom.account || 0,
        categoryId: ids[ids.length - 1],
        amount: this.amount,
        date: resultDate,
        comment: this.comment,
      };

      console.log(param);
      let resp = await this.$http.post(window.api.record.create, param);
      console.log(resp);
      if (resp.data.code !== 0) {
        this.$message({
          message: resp.data.msg,
          type: "warning",
        });
      } else {
        this.$message({
          message: "新增 " + resp.data.data.length + " 条",
          type: "success",
        });
      }
    },
    hasChange(val) {
      this.accountId = val;
    },
    listenTargetAccount(val) {
      this.targetAccountId = val;
    },
  },
};
</script>