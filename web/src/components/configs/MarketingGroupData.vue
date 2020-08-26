<template>
  <div class="marketingGroupData">
    <el-col :span="12">
      <el-form-item label="分组名称：">
        <el-input
          type="text"
          placeholder="请输入分组名称"
          v-model="groupName"
        />
      </el-form-item>
    </el-col>
    <el-col :span="12">
      <el-form-item label="群组分红归属：" label-width="150px">
        <el-input
          type="number"
          placeholder="请输入群组归属人员账户ID"
          v-model="userID"
        />
      </el-form-item>
    </el-col>
  </div>
</template>
<script>
export default {
  name: "MarketingGroupData",
  props: {
    data: String
  },
  data() {
    const data = {
      groupName: "",
      userID: null
    };
    if (this.$props.data) {
      Object.assign(data, JSON.parse(this.$props.data));
    }
    return data;
  },
  watch: {
    groupName() {
      this.handleChange();
    },
    userID() {
      this.handleChange();
    }
  },
  methods: {
    handleChange() {
      const { groupName, userID } = this;
      let value = "";
      if (groupName) {
        value = JSON.stringify({
          groupName,
          userID
        });
      }
      this.$emit("change", value);
    }
  }
};
</script>
