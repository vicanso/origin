<template>
  <BaseEditor
    v-if="fields"
    title="添加/更新产品信息"
    icon="el-icon-files"
    :fields="fields"
    :id="productID"
    :findByID="getProductByID"
    :updateByID="updateProductByID"
  />
</template>
<script>
import { mapActions } from "vuex";
import BaseEditor from "@/components/base/Editor.vue";

const productStatuses = [];
const productCategories = [];
const fields = [
  {
    label: "名称：",
    key: "name",
    clearable: true,
    placeholder: "请输入产品名称"
  },
  {
    label: "单价：",
    key: "price",
    clearable: true,
    placeholder: "请输入产品单价"
  },
  {
    label: "单位：",
    key: "unit",
    clearable: true,
    placeholder: "请输入产品单位"
  },
  {
    label: "状态：",
    key: "status",
    placeholder: "请选择产品状态",
    type: "select",
    options: productStatuses
  },
  {
    label: "分类：",
    key: "categories",
    placeholder: "请选择产品分类",
    type: "select",
    multiple: true,
    options: productCategories
  },
  {
    label: "品牌：",
    key: "brand",
    placeholder: "请选择产品品牌",
    type: "brand"
  },
  {
    label: "产地：",
    key: "origin",
    placeholder: "请选择产品产地",
    type: "region",
    maxLevel: 2,
    showAllLevels: true
  },
  {
    label: "SN：",
    key: "sn",
    placeholder: "请输入产品SN码"
  },
  {
    label: "主图：",
    key: "mainPic",
    placeholder: "请输入主图位置"
  },
  {
    label: "图片：",
    key: "files",
    span: 24,
    type: "upload",
    bucket: "origin-pics"
  },
  {
    label: "开始时间：",
    key: "startedAt",
    type: "datePicker",
    pickerType: "datetime",
    placeholder: "请选择产品开始销售时间",
    labelWidth: "100px"
  },
  {
    label: "结束时间：",
    key: "endedAt",
    type: "datePicker",
    pickerType: "datetime",
    placeholder: "请选择产品结束销售时间",
    labelWidth: "100px"
  },
  {
    label: "关键字：",
    key: "keywords",
    placeholder: "请输入产品关键字，多个关键字以空格分开"
  },
  {
    label: "简介：",
    key: "catalog",
    type: "textarea",
    autosize: { minRows: 5, maxRows: 10 },
    span: 24,
    placeholder: "请输入产品简介"
  }
];

export default {
  name: "Product",
  components: {
    BaseEditor
  },
  data() {
    return {
      fields: null,
      productID: 0,
      processing: false
    };
  },
  methods: {
    ...mapActions([
      "listProductStatus",
      "listBrand",
      "getProductByID",
      "updateProductByID"
    ])
  },
  async beforeMount() {
    try {
      const { statuses } = await this.listProductStatus();
      productStatuses.length = 0;
      productStatuses.push(...statuses);
      this.fields = fields;
    } catch (err) {
      this.$message.error(err.message);
    }
  }
};
</script>
<style lang="sass" scoped>
$uploadWidth: 350px
.product
  .upload
    float: right
    width: $uploadWidth
  h5
    margin: 0 0 10px 0
  .form
    margin-right: $uploadWidth + 10
  .selector, .btn, .dateRange
    width: 100%
</style>
