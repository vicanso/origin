<template>
  <div class="brand">
    <BaseEditor
      v-if="!processing && fields"
      title="更新地区信息"
      icon="el-icon-position"
      :id="id"
      :findByID="getRegionByID"
      :updateByID="updateRegionByID"
      :fields="fields"
    />
  </div>
</template>
<script>
import BaseEditor from "@/components/base/Editor.vue";
import { mapActions } from "vuex";
const brandStatuses = [];
const fields = [
  {
    label: "名称：",
    key: "name",
    clearable: true,
    span: 6,
    rules: [
      {
        required: true,
        message: "地区名称不能为空"
      }
    ]
  },
  {
    label: "代码：",
    key: "code",
    span: 6,
    disabled: true
  },
  {
    label: "优先级",
    key: "priority",
    dataType: "number",
    span: 6,
    clearable: true,
    placeholder: "请输入优先级(1-100)"
  },
  {
    label: "状态：",
    key: "status",
    type: "select",
    span: 6,
    options: brandStatuses
  }
];

export default {
  name: "Brand",
  components: {
    BaseEditor
  },
  data() {
    return {
      fields: null,
      processing: false,
      id: 0
    };
  },
  methods: mapActions(["listStatus", "updateRegionByID", "getRegionByID"]),
  async beforeMount() {
    this.processing = true;
    const { id } = this.$route.query;
    if (id) {
      this.id = Number(id);
    }
    try {
      const { statuses } = await this.listStatus();
      brandStatuses.length = 0;
      brandStatuses.push(...statuses);
      this.fields = fields;
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
