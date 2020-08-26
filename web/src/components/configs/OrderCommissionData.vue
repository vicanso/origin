<template>
  <div class="orderCommissionData">
    <el-col :span="12">
      <el-form-item label="佣金分组：">
        <el-input
          type="text"
          placeholder="请输入佣金分组，通用的则填 *"
          v-model="group"
        />
      </el-form-item>
    </el-col>
    <el-col :span="12">
      <el-form-item label="佣金比例：">
        <el-input
          type="number"
          placeholder="请输入佣金比例（小数）"
          v-model="ratio"
        />
      </el-form-item>
    </el-col>
  </div>
</template>
<script>
export default {
  name: "OrderCommissionData",
  props: {
    data: String
  },
  data() {
    const data = {
      group: "",
      ratio: 0
    };
    if (this.$props.data) {
      Object.assign(data, JSON.parse(this.$props.data));
    }
    return data;
  },
  watch: {
    group() {
      this.handleChange();
    },
    ratio() {
      this.handleChange();
    }
  },
  methods: {
    handleChange() {
      const { group, ratio } = this;
      let value = "";
      if (ratio > 0.1) {
        this.$message.error("佣金比例不能大于0.1");
        return;
      }
      if (group) {
        value = JSON.stringify({
          group,
          ratio
        });
      }
      this.$emit("change", value);
    }
  }
};
</script>
