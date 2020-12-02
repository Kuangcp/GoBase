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
            v-model.number="amount"
            size="mini"
            clearable
            min="0"
            style="width: 120px"
        />
      </el-form-item>

      <el-form-item label="时间" required>
        <el-date-picker
            v-model="recordDate"
            type="dates"
            size="mini"
            clearable
            style="width: 200px"
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
      accountId: 2,
      targetAccountId: 1,
      amount: 0,
      recordDate: [],
      comment: "",
    };
  },
  methods: {
    async onSubmit() {
      let ids = this.$refs.categoryCom.categoryId;
      if (this.recordDate.length === 0) {
        this.$message({
          message: "时间为空",
          type: "warning",
        });
        return;
      }

      let resultDate = this.recordDate.map((v) => DateUtil(v).formatDate());

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
      if (resp.data.Code !== 0) {
        this.$message({
          message: resp.data.Msg,
          type: "warning",
        });
      } else {
        this.$message({
          message: "新增 " + resp.data.Data.length + " 条",
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