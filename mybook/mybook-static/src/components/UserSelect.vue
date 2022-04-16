<template>
  <el-select
      v-model="tempVal"
      size="mini"
      clearable
      placeholder="请选择"
      @change="onChanges"
  >
    <el-option
        v-for="item in users"
        :key="item.ID"
        :label="item.Name"
        :value="item.ID"
    >
    </el-option>
  </el-select>
</template>
<script>
export default {
  name: "UserSelect",
  props: {
    user: {
      type: Number,
    },
  },
  data: function () {
    return {tempVal: this.user, users: []};
  },
  mounted() {
    this.fillUserList();
  },
  methods: {
    async fillUserList() {
      const res = await this.$http.get(window.api.user.listAll);
      this.users = res.data.data;
    },
    onChanges(val) {
      this.$emit("hasChange", val);
    },
  },
};
</script>