<template>
  <el-select
      v-model="tempVal"
      size="mini"
      clearable
      placeholder="请选择"
      @change="onChanges"
  >
    <el-option
        v-for="item in accounts"
        :key="item.ID"
        :label="item.Name"
        :value="item.ID"
    >
    </el-option>
  </el-select>
</template>
<script>
export default {
  props: {
    account: {
      type: Number,
      // default: 1,
    },
  },
  data: function () {
    return {tempVal: this.account, accounts: []};
  },
  mounted() {
    this.fillAccount();
  },
  methods: {
    async queryAllAccount() {
      const res = await this.$http.get(window.api.account.listAll);
      // console.log("ren", res.data);
      return res.data.data;
    },
    async fillAccount() {
      this.accounts = [];
      let result = await this.queryAllAccount();
      this.accounts = result;
    },
    onChanges(val) {
      console.log("child", val);
      this.$emit("hasChange", val);
    },
  },
};
</script>