<template>
  <div>
    <el-form :inline="true" label-width="80px" size="mini">
      <el-form-item label="" required>
        <el-input placeholder="用户名" v-model="userName" clearable/>
      </el-form-item>

      <el-form-item label="" required>
        <el-button @click="addUser">提交</el-button>
      </el-form-item>
    </el-form>

    <el-form :inline="true" label-width="80px" size="mini">
      <el-form-item label="" required>
        <el-input placeholder="账户名" v-model="accountName" clearable/>
      </el-form-item>

      <el-form-item label="" required>
        <el-button @click="addAccount">提交</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: "Setting",
  data: function () {
    return {
      userName: '',
      accountName: ''
    }
  },
  methods: {
    async addUser() {
      let resp = await this.$http.get(window.api.user.addUser + "?name=" + this.userName);
      console.log(resp)

      if (resp.data.code !== 0) {
        this.$message({
          message: resp.data.msg,
          type: "warning",
        });
      } else {
        this.$message({
          message: "新增 " + this.userName,
          type: "success",
        });
      }
    },
    async addAccount() {

    }
  }
}
</script>

<style scoped>

</style>