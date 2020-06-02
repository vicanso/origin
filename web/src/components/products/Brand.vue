<template>
  <BaseEditor
    v-if="!processing && fields"
    title="添加/更新品牌信息"
    icon="el-icon-good"
    :id="brandID"
    :findByID="getBrandByID"
    :updateByID="updateBrandByID"
    :fields="fields"
    :add="addBrand"
    :rules="rules"
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
    span: 12
  },
  {
    label: "状态：",
    key: "status",
    span: 12,
    type: "select",
    placeholder: "请选择产品状态",
    options: brandStatuses
  },
  {
    label: "LOGO：",
    key: "files",
    span: 24,
    type: "upload",
    bucket: "origin-pics"
  },
  {
    label: "简介：",
    placeholder: "请输入品牌简介",
    key: "catalog",
    type: "textarea",
    span: 24,
    autosize: { minRows: 5, maxRows: 10 }
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
      brandID: 0,
      processing: false,

      rules: {
        name: [
          {
            required: true,
            message: "品牌名称不能为空"
          }
        ],
        status: [
          {
            required: true,
            message: "品牌状态不能为空"
          }
        ],
        files: [
          {
            required: true,
            message: "品牌图标不能为空"
          }
        ],
        catalog: [
          {
            required: true,
            message: "品牌简介不能为空"
          }
        ]
      }
    };
  },
  methods: {
    ...mapActions([
      "listStatus",
      "addBrand",
      "getBrandByID",
      "updateBrandByID",
      "addBrand"
    ])
  },
  async beforeMount() {
    this.processing = true;
    const { id } = this.$route.query;
    if (id) {
      this.brandID = Number(id);
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
