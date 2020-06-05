<template>
  <BaseEditor
    v-if="!processing && fields"
    title="添加/更新品牌信息"
    icon="el-icon-good"
    :id="id"
    :findByID="getBrandByID"
    :updateByID="updateBrandByID"
    :fields="fields"
    :add="addBrand"
  />
</template>
<script>
import { mapActions } from "vuex";
import BaseEditor from "@/components/base/Editor.vue";

const brandStatuses = [];
const fields = [
  {
    label: "名称：",
    key: "name",
    clearable: true,
    placeholder: "请输入品牌名称",
    span: 12,
    rules: [
      {
        required: true,
        message: "品牌名称不能为空"
      }
    ]
  },
  {
    label: "状态：",
    key: "status",
    span: 12,
    type: "select",
    placeholder: "请选择产品状态",
    options: brandStatuses,
    rules: [
      {
        required: true,
        message: "品牌状态不能为空"
      }
    ]
  },
  {
    label: "LOGO：",
    key: "files",
    span: 24,
    type: "upload",
    bucket: "origin-pics",
    rules: [
      {
        required: true,
        message: "品牌图标不能为空"
      }
    ]
  },
  {
    label: "简介：",
    placeholder: "请输入品牌简介",
    key: "catalog",
    type: "textarea",
    span: 24,
    autosize: { minRows: 5, maxRows: 10 },
    rules: [
      {
        required: true,
        message: "品牌简介不能为空"
      }
    ]
  }
];

export default {
  name: "ProductBrand",
  components: {
    BaseEditor
  },
  data() {
    return {
      fields: null,
      id: 0,
      processing: false
    };
  },
  methods: {
    ...mapActions(["listStatus", "addBrand", "getBrandByID", "updateBrandByID"])
  },
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
