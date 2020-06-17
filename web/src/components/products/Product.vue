<template>
  <BaseEditor
    v-if="fields"
    title="添加/更新产品信息"
    icon="el-icon-files"
    :fields="fields"
    :id="id"
    :findByID="getProductByID"
    :updateByID="updateProductByID"
    :add="addProduct"
  />
</template>
<script>
import { mapActions } from "vuex";
import BaseEditor from "@/components/base/Editor.vue";

const productStatuses = [];
const fields = [
  {
    label: "名称：",
    key: "name",
    clearable: true,
    placeholder: "请输入产品名称",
    rules: [
      {
        required: true,
        message: "产品名称不能为空"
      }
    ]
  },
  {
    label: "单价：",
    key: "price",
    clearable: true,
    dataType: "number",
    placeholder: "请输入产品单价",
    rules: [
      {
        required: true,
        message: "产品单价不能为空"
      }
    ]
  },
  {
    label: "单位：",
    type: "specsUnit",
    key: "specs",
    dataType: "number",
    placeholder: "请输入产品规格",
    selectKey: "unit",
    clearable: true,
    selectPlaceholder: "请输入产品单位",
    options: [
      {
        name: "盒",
        value: "盒"
      },
      {
        name: "份",
        value: "份 "
      },
      {
        name: "克",
        value: "克"
      },
      {
        name: "个",
        value: "个"
      }
    ],
    rules: [
      {
        required: true,
        message: "产品规格不能为空"
      }
    ]
  },
  {
    label: "状态：",
    key: "status",
    placeholder: "请选择产品状态",
    type: "select",
    options: productStatuses,
    rules: [
      {
        required: true,
        message: "产品状态不能为空"
      }
    ]
  },
  {
    label: "分类：",
    key: "categories",
    placeholder: "请选择产品分类",
    type: "productCategory",
    rules: [
      {
        required: true,
        message: "产品分类不能为空"
      }
    ]
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
    label: "供应商",
    key: "supplier",
    placeholder: "请选择供应商",
    type: "supplier",
    rules: [
      {
        required: true,
        message: "商品供应商不能为空"
      }
    ]
  },
  {
    label: "SN：",
    key: "sn",
    placeholder: "请输入产品SN码"
  },
  {
    label: "排序：",
    key: "rank",
    dataType: "number",
    placeholder: "请输入产品排序(1-1000)"
  },
  {
    label: "主图：",
    key: "mainPic",
    placeholder: "请输入主图位置",
    dataType: "number"
  },
  {
    label: "图片：",
    key: "files",
    span: 24,
    type: "upload",
    limit: 10
  },
  {
    label: "开始时间：",
    key: "startedAt",
    type: "datePicker",
    pickerType: "datetime",
    placeholder: "请选择产品开始销售时间",
    labelWidth: "100px",
    rules: [
      {
        required: true,
        message: "产品开始销售时间不能为空"
      }
    ]
  },
  {
    label: "结束时间：",
    key: "endedAt",
    type: "datePicker",
    pickerType: "datetime",
    placeholder: "请选择产品结束销售时间",
    labelWidth: "100px",
    rules: [
      {
        required: true,
        message: "产品结束销售时间不能为空"
      }
    ]
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
    placeholder: "请输入产品简介",
    rules: [
      {
        required: true,
        message: "产品简介不能为空"
      }
    ]
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
      id: 0,
      processing: false
    };
  },
  methods: {
    ...mapActions([
      "listStatus",
      "listBrand",
      "getProductByID",
      "updateProductByID",
      "addProduct"
    ])
  },
  async beforeMount() {
    const { id } = this.$route.query;
    if (id) {
      this.id = Number(id);
    }
    try {
      const { statuses } = await this.listStatus();
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
