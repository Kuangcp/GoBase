<template>
  <el-cascader
    v-model="categoryId"
    :options="categoryMap"
    :props="{ expandTrigger: 'hover' }"
    @change="handleChange"
  ></el-cascader>
</template>
<script>
export default {
  data: function () {
    return { categoryId: [], categoryMap: [] };
  },
  mounted() {
    this.fillAccount();
  },
  methods: {
    async queryAllAccount() {
      const res = await this.$http.get("/api/category/listTree");
      console.log("ren", res.data);
      return res.data.Data;
    },
    async fillAccount() {
      this.categoryMap = [];
      let result = await this.queryAllAccount();
      this.categoryMap = result;
    },
    handleChange(value) {
      console.log(value);
    },
  },
};
</script>