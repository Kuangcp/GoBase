<template>
  <el-select v-model="account" size="mini" clearable placeholder="请选择">
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
  data: function () {
    return { account: "", accounts: [] };
  },
  mounted() {
    this.fillAccount();
  },
  methods: {
    async queryAllAccount() {
      const res = await this.$http.get("/api/account/list");
      // console.log("ren", res.data);
      return res.data.Data;
    },
    async fillAccount() {
      this.accounts = [];
      let result = await this.queryAllAccount();
      this.accounts = result;
    },
  },
};
</script>