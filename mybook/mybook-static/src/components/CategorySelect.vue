<template>
  <div>
    <el-cascader
        v-model="categoryId"
        :options="categoryMap"
        :props="{ expandTrigger: 'hover' }"
        @change="handleChange"
        size="mini"
        style="width: 356px"
        clearable
    ></el-cascader>

    <!-- <span :key="index" v-for="(item, index) in categoryMap">
      <el-radio v-model="radio" :label="item.value">{{ item.label }}</el-radio>
    </span>

    <el-divider></el-divider>
    <div :key="index" v-for="(item, index) in categoryMap">
      <div v-if="item.value == radio" style="height: 60px">
        <span
          :key="childIndex"
          v-for="(childItem, childIndex) in item.children"
        >
          <el-radio v-model="leaf" :label="childItem.label"></el-radio>
        </span>
      </div>
    </div> -->
  </div>
</template>
<script>
export default {
  data: function () {
    return {leaf: 0, radio: 1, categoryId: [], categoryMap: []};
  },
  mounted() {
    this.fillCategoryList();
  },
  methods: {
    async fillCategoryList() {
      const res = await this.$http.get(window.api.category.tree);
      this.categoryMap = res.data.data;
    },
    handleChange(value) {
      console.log(value);
    },
  },
};
</script>
