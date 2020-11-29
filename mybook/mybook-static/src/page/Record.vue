<template>
  <div>
    <el-form :inline="true" class="demo-form-inline">
      <el-form-item>
        <AccountSelect />
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
  data() {
    return {
      activeName: "detail",
    };
  },
  methods: {},
};
</script>